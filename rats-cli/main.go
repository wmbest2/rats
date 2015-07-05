package main

import (
	"encoding/xml"
	"fmt"
	"github.com/wmbest2/android/apk"
	"github.com/wmbest2/rats/agent/android"
	"github.com/wmbest2/rats/rats"
	"os"
)

func main() {
	argCount := len(os.Args)
	if argCount != 2 && argCount != 3 {
		fmt.Println("Usage: cli-client <main apk [optional]> <test apk>")
		fmt.Println("   * main apk not required for library tests")
		return
	}

	devices := <-rats.GetAllDevices()

	var manifest *apk.Manifest

	for _, arg := range os.Args[1:] {
		file, err := os.Open(arg)
		if err != nil {
			panic(err)
		}
		rats.Install("apk.apk", file, devices...)
		if manifest == nil {
			fi, _ := file.Stat()
			manifest = rats.GetManifest(file, fi.Size())
		}

		file.Close()
	}

	for _, device := range devices {
		device.SetScreenOn(true)
		device.Unlock()
	}

	//testFile := os.Args[len(os.Args)-1]

	_, runs := android.RunTests(manifest, devices)

	str, err := xml.Marshal(*<-runs)
	if err == nil {
		fmt.Println(string(str))
	}

	rats.Uninstall(manifest.Package, devices...)
	rats.Uninstall(manifest.Instrument.Target, devices...)
}
