package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/wmbest2/android/adb"
	"github.com/wmbest2/android/apk"
	"github.com/wmbest2/rats/agent/android"
	"github.com/wmbest2/rats/core"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	jsonEnabled = flag.Bool("json", false, "Log debug information")
)

func init() {
	flag.Parse()
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
	core.UpdateAdb(a)
}

func main() {
	argCount := len(flag.Args())
	if argCount != 2 && argCount != 3 {
		fmt.Println("Usage: cli-client <main apk [optional]> <test apk>")
		fmt.Println("   * main apk not required for library tests")
		return
	}

	var devices []*core.Device
	for len(devices) == 0 {
		devices = <-core.GetAllDevices()
	}

	var manifest *apk.Manifest

	for _, arg := range flag.Args()[1:] {
		file, err := os.Open(arg)
		if err != nil {
			panic(err)
		}

		fi, _ := file.Stat()
		man := core.GetManifest(file, fi.Size())
		log.Printf("Installing package %s\n", man.Package)

		manifest = man

		//rats.Install(file.Name(), file, devices...)

		log.Printf("\t -> Install Complete\n")

		file.Close()
	}

	for _, device := range devices {
		device.SetScreenOn(true)
		device.Unlock()
	}

	log.Printf("Running Tests\n")

	coverageFile := fmt.Sprintf("/data/data/%s/files/coverage.ec", manifest.Instrument.Target)

	finished, runs := android.RunTests(manifest, devices, []string{coverageFile})

	for range finished {
		run := <-runs
		var str []byte
		var err error
		if *jsonEnabled {
			str, err = json.Marshal(run)
		} else {
			str, err = xml.Marshal(run)
		}

		log.Printf("\t -> Received results...\n")

		if err == nil {
			fmt.Println(string(str))
		}
	}

	//rats.Uninstall(manifest.Package, devices...)
	//rats.Uninstall(manifest.Instrument.Target, devices...)
}
