package main

import (
    "github.com/wmbest2/adb"
    "github.com/wmbest2/adb/apk"
    "fmt"
    "sync"
    "log"
    "os"
    "archive/zip"
    "io/ioutil"
    "encoding/xml"
)

func runOnDevice(wg *sync.WaitGroup, d *adb.Device, params []string) {
    defer wg.Done()
    fmt.Printf("%s\n", d)
    v,_ := d.AdbExec(params...)
    fmt.Printf("%s\n", string(v))
}

func runOnAll(params ...string) {
    var wg sync.WaitGroup
    deviceLock.Lock()
    for _,d := range devices {
        wg.Add(1)
        go runOnDevice(&wg, d, params)
    }
    wg.Wait()
    deviceLock.Unlock()
}

func install(file string) {
    runOnAll("install", file)
}

func getReaderFromZip(file string, subFile string) []byte {
    r, err := zip.OpenReader(file)
    if err != nil {
        log.Fatal(err)
    }
    defer r.Close()

    // Iterate through the files in the archive,
    // printing some of their contents.
    for _, f := range r.File {
        if (f.Name == subFile) {
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

func getTestInfo(file string) *apk.Manifest {
    body := getReaderFromZip(file, "AndroidManifest.xml") 
    manifest_body := apk.DecompressXML(body)

    var manifest apk.Manifest

    err := xml.Unmarshal([]byte(manifest_body), &manifest)
    if err != nil {
            fmt.Printf("error: %v", err)
            return nil
    }

    return &manifest 
}

func main() {
    updateDevices()

    argCount := len(os.Args)
    if (argCount != 2 && argCount != 3) {
        fmt.Println("Usage: cli-client <main apk [optional]> <test apk>")
        fmt.Println("   * main apk not required for library tests")
        return
    }

    for _, arg := range os.Args[1:] {
        install(arg)
    }

    testFile := os.Args[len(os.Args) -1]
    manifest := getTestInfo(testFile)

    fmt.Printf("Test File: %s\n", testFile)
    fmt.Printf("Test Package: %s\n", manifest.Package)
    fmt.Printf("Instrumentation: %s\n", manifest.Instrument.Name)
}
