package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/docker/libchan"
	"github.com/docker/libchan/spdy"
	"github.com/wmbest2/android/adb"
	"github.com/wmbest2/rats/agent/android"
	"github.com/wmbest2/rats/agent/proto"
	"github.com/wmbest2/rats/core"
	"github.com/wmbest2/rats/test"
)

var (
	adb_address = flag.String("adb_address", "localhost", "Address of ADB server")
	adb_port    = flag.Int("adb_port", 5037, "Port of ADB server")
)

func tryStartAdb() {
	path := os.ExpandEnv("$ANDROID_HOME")
	if path != "" {
		path = filepath.Join(path, "platform-tools", "adb")
		b, err := exec.Command(path, "start-server").CombinedOutput()
		if err != nil {
			log.Println(err)
		} else {
			log.Println(string(b))
		}
	}
}

func refreshDevices(a *adb.Adb, inRecover bool) {
	defer func() {
		if e := recover(); e != nil {
			if !inRecover && a == adb.Default {
				log.Println("Couldn't connect to adb, attempting to recover")
				tryStartAdb()
				refreshDevices(a, true)
			} else if inRecover {
				log.Fatalf("Still couldn't connect.  Make sure adb exists in $ANDROID_HOME\n\tor manually start it with 'adb start-server'")
			} else {
				log.Fatal(e)
			}
		}
	}()
	core.UpdateAdb(a)
}

func run(p *proto.Run) *test.TestRun {
	log.Println("Starting new test run")
	start := time.Now()

	var serials []string
	if p.Metadata["serialList"] != "" {
		serials = strings.Split(p.Metadata["serialList"], ",")
	}

	size := len(p.Binary["test"])

	count, _ := strconv.Atoi(p.Metadata["count"])
	filter := &core.DeviceFilter{
		Count:  count,
		Strict: p.Metadata["strict"] == "true",
	}
	filter.Serials = serials

	manifest := core.GetManifest(bytes.NewReader(p.Binary["test"]), int64(size))
	filter.MinSdk = manifest.Sdk.Min
	filter.MaxSdk = manifest.Sdk.Max

	devices := <-core.GetDevices(filter)
	core.Reserve(devices...)

	// Remove old if left over
	core.Uninstall(manifest.Package, devices...)
	core.Uninstall(manifest.Instrument.Target, devices...)

	// Install New
	if main != nil {
		core.Install("main.apk", bytes.NewBuffer(p.Binary["main"]), devices...)
	}
	core.Install("test.apk", bytes.NewBuffer(p.Binary["test"]), devices...)

	core.Unlock(devices)

	artifacts := []string{"coverage.ec"}

	finished, out := android.RunTests(manifest, devices, artifacts)

	var s *test.TestRun
SuitesLoop:
	for {
		select {
		case s = <-out:
			break SuitesLoop
		case device := <-finished:
			go func() {
				core.Uninstall(manifest.Package, device)
				core.Uninstall(manifest.Instrument.Target, device)
				core.Release(device)
			}()
		}
	}

	s.Name = p.Metadata["uuid"]
	s.Timestamp = time.Now()
	if p.Metadata["msg"] != "" {
		s.Message = test.NewNullString(p.Metadata["msg"])
	}

	log.Printf("Test run completed in %s across %d device(s)\n", time.Since(start), len(devices))

	return s
}

func init() {
	flag.Parse()

	conn := adb.Default
	if *adb_address != "localhost" || *adb_port != 5037 {
		conn = adb.Connect(*adb_address, *adb_port)
	}
	go refreshDevices(conn, false)
}

func main() {
	client, err := net.Dial("tcp", "localhost:4000")
	if err != nil {
		log.Fatal(err)
	}

	p, err := spdy.NewSpdyStreamProvider(client, true)
	if err != nil {
		log.Print(err)
	}
	tran := spdy.NewTransport(p)

	sender, err := tran.NewSendChannel()
	if err != nil {
		log.Fatal(err)
	}

	receiver, remoteSender := libchan.Pipe()

	msg := proto.Message{
		Command:   proto.Register,
		Responder: remoteSender,
	}

	err = sender.Send(msg)
	if err != nil {
		log.Fatal(err)
	}

	// Wait for the response to unblock
	resp := &proto.Message{}
	err = receiver.Receive(resp)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %s\n", resp.Command)

	for {
		log.Printf("HERHEHRHER START")
		msg := &proto.Message{}
		err := receiver.Receive(msg)
		if err != nil {
			log.Fatalf("Receive error: %s\n", err)
		}
		log.Printf("Response: %s\n", msg.Command)

		switch msg.Command {
		case proto.Init:
			result := run(msg.Run)

			b, err := json.Marshal(result)
			if err != nil {
				log.Fatal(err)
			}

			err = msg.Responder.Send(proto.Message{
				Command: proto.Complete,
				Result:  b,
			})
			if err != nil {
				log.Fatal(err)
			}
		case proto.Devices:
			b, err := json.Marshal(<-core.GetAllDevices())

			if err != nil {
				log.Fatal(err)
			}

			err = msg.Responder.Send(proto.Message{
				Command: proto.Complete,
				Result:  b,
			})
			if err != nil {
				log.Fatal(err)
			}
		}

		log.Printf("HERHEHRHER STOP")
	}
}
