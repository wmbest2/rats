package main

import (
	"encoding/xml"
	"fmt"
	"github.com/wmbest2/android/adb"
	"github.com/wmbest2/android/apk"
	"github.com/wmbest2/rats/agent/android"
	"github.com/wmbest2/rats/rats"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func init() {
	conn := adb.Default
	go refreshDevices(conn, false)
}

func tryStartAdb() {
	path := os.ExpandEnv("$ANDROID_HOME")
	if path != "" {
		path = filepath.Join(path, "platform-tools", "adb")
		exec.Command(path, "start-server").CombinedOutput()
	}
}

func refreshDevices(a *adb.Adb, inRecover bool) {
	defer func() {
		if e := recover(); e != nil {
			if !inRecover && a == adb.Default {
				tryStartAdb()
				refreshDevices(a, true)
			}
		}
	}()
	rats.UpdateAdb(a)
}

func main() {
	argCount := len(os.Args)
	if argCount != 2 && argCount != 3 {
		log.Println("Usage: cli-client <main apk [optional]> <test apk>")
		log.Println("   * main apk not required for library tests")
		return
	}

	var devices []*rats.Device
	for len(devices) == 0 {
		devices = <-rats.GetAllDevices()
	}

	var manifest *apk.Manifest

	for _, arg := range os.Args[1:] {
		file, err := os.Open(arg)
		if err != nil {
			panic(err)
		}

		fi, _ := file.Stat()
		man := rats.GetManifest(file, fi.Size())
		log.Printf("Installing package %s\n", man.Package)

		manifest = man

		rats.Install(file.Name(), file, devices...)

		log.Printf("\t -> Install Complete\n")

		file.Close()
	}

	for _, device := range devices {
		device.SetScreenOn(true)
		device.Unlock()
	}

	log.Printf("Running Tests\n")

	finished, runs := android.RunTests(manifest, devices)

	for range finished {
		run := <-runs
		str, err := xml.Marshal(run)

		log.Printf("\t -> Received results...\n")

		if err == nil {
			fmt.Println(string(str))
		}
	}

	rats.Uninstall(manifest.Package, devices...)
	rats.Uninstall(manifest.Instrument.Target, devices...)
}
