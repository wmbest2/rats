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
	"github.com/wmbest2/rats/rats"
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
	rats.UpdateAdb(a)
}

func run(p *proto.Run) *test.TestSuites {
	log.Println("Starting new test run")
	start := time.Now()

	var serials []string
	if p.Metadata["serialList"] != "" {
		serials = strings.Split(p.Metadata["serialList"], ",")
	}

	size := len(p.Binary["test"])

	count, _ := strconv.Atoi(p.Metadata["count"])
	filter := &rats.DeviceFilter{
		Count:  count,
		Strict: p.Metadata["strict"] == "true",
	}
	filter.Serials = serials

	manifest := rats.GetManifest(bytes.NewReader(p.Binary["test"]), int64(size))
	filter.MinSdk = manifest.Sdk.Min
	filter.MaxSdk = manifest.Sdk.Max

	devices := <-rats.GetDevices(filter)
	rats.Reserve(devices...)

	// Remove old if left over
	rats.Uninstall(manifest.Package, devices...)
	rats.Uninstall(manifest.Instrument.Target, devices...)

	// Install New
	if main != nil {
		rats.Install("main.apk", bytes.NewBuffer(p.Binary["main"]), devices...)
	}
	rats.Install("test.apk", bytes.NewBuffer(p.Binary["test"]), devices...)

	rats.Unlock(devices)

	finished, out := android.RunTests(manifest, devices)

	var s *test.TestSuites
SuitesLoop:
	for {
		select {
		case s = <-out:
			break SuitesLoop
		case device := <-finished:
			go func() {
				rats.Uninstall(manifest.Package, device)
				rats.Uninstall(manifest.Instrument.Target, device)
				rats.Release(device)
			}()
		}
	}

	s.Name = p.Metadata["uuid"]
	s.Timestamp = time.Now()
	if p.Metadata["msg"] != "" {
		s.Message = p.Metadata["msg"]
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

	tran, err := spdy.NewClientTransport(client)
	if err != nil {
		log.Fatal(err)
	}

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
		msg := &proto.Message{}
		err := receiver.Receive(msg)
		if err != nil {
			log.Fatalf("Receive error: %s\n", err)
		}

		result := run(msg.Run)

		b, err := json.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}

		err = msg.Responder.Send(proto.Message{
			Command:   proto.Complete,
			Result:    b,
			Responder: sender,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
