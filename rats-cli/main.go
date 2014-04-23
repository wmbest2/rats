package main

import (
	"encoding/xml"
	"fmt"
	"github.com/wmbest2/android/apk"
	"github.com/wmbest2/rats_server/rats"
	"github.com/wmbest2/rats_server/test"
	"os"
)


func runTests(manifest *apk.Manifest) *test.TestSuites {
	out := make(chan *test.TestSuite)
	suites := &test.TestSuites{}

	for _, d := range rats.Devices {
		go test.RunTest(d, manifest, out)
	}

	for _ = range rats.Devices {
		suite := <-out
		suites.TestSuites = append(suites.TestSuites, suite)
		suites.Time += suite.Time
	}

    return suites
}

func main() {
	rats.UpdateDevices()

	argCount := len(os.Args)
	if argCount != 2 && argCount != 3 {
		fmt.Println("Usage: cli-client <main apk [optional]> <test apk>")
		fmt.Println("   * main apk not required for library tests")
		return
	}

	for _, arg := range os.Args[1:] {
		rats.Install(arg)
	}

	testFile := os.Args[len(os.Args)-1]
	manifest := rats.GetManifest(testFile)

    s := runTests(manifest)
    str, err := xml.Marshal(s)
    if err == nil {
        fmt.Println(string(str))
    }

	rats.Uninstall(manifest.Package)
	rats.Uninstall(manifest.Instrument.Target)
}
