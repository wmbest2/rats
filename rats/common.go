package rats

import (
	"archive/zip"
	"fmt"
	"github.com/wmbest2/android/adb"
	"github.com/wmbest2/android/apk"
	"io/ioutil"
	"log"
	"sync"
)


func RunOnDevice(wg *sync.WaitGroup, d adb.AdbRunner, params []string) {
	defer wg.Done()
	d.ExecSync(params...)
}

func RunOnAll(params ...string) {
	var wg sync.WaitGroup
	for _, d := range <-GetDevices() {
		wg.Add(1)
		go RunOnDevice(&wg, d, params)
	}
	wg.Wait()
}

func Install(file string) {
	RunOnAll("install", "-r", file)
}

func Uninstall(pack string) {
	RunOnAll("uninstall", pack)
}

func GetFileFromZip(file string, subFile string) []byte {
	r, err := zip.OpenReader(file)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
		if f.Name == subFile {
			var body []byte
			rc, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}
			body, err = ioutil.ReadAll(rc)
			if err != nil {
				log.Fatal(err)
			}
			rc.Close()

			return body
		}
	}
	return []byte{}
}

func GetManifest(file string) *apk.Manifest {
	var manifest apk.Manifest

	body := GetFileFromZip(file, "AndroidManifest.xml")
	err := apk.Unmarshal([]byte(body), &manifest)

	if err != nil {
		fmt.Printf("error: %v", err)
		return nil
	}

	return &manifest
}
