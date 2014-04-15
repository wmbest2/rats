package main

import (
    "github.com/wmbest2/android/adb"
    "github.com/wmbest2/android/apk"
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

func getFileFromZip(file string, subFile string) []byte {
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
    var manifest apk.Manifest

    body := getFileFromZip(file, "res/xml/volley_ball_routes.xml") 
    err := apk.Unmarshal([]byte(body), &manifest)

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

    /*for _, arg := range os.Args[1:] {*/
        /*install(arg)*/
    /*}*/

    testFile := os.Args[len(os.Args) -1]
    manifest := getTestInfo(testFile)

    o, _ := xml.Marshal(manifest)
    fmt.Printf("Manifest: %s",o)

    /*fmt.Printf("Test File: %s\n", testFile)*/
    /*fmt.Printf("Test Package: %s\n", manifest.Package)*/
    /*fmt.Printf("Instrumentation: %s\n", manifest.Instrument.Name)*/
}
