package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
    "fmt"
	"github.com/go-martini/martini"
	"github.com/wmbest2/rats_server/rats"
	"github.com/wmbest2/rats_server/test"
    "io"
	"net/http"
	"path/filepath"
	"runtime"
    "os"
)

func uuid() (string, error) {
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// TODO: verify the two lines implement RFC 4122 correctly
	uuid[8] = 0x80 // variant bits see page 5
	uuid[4] = 0x40 // version 4 Pseudo Random, see page 7

	return hex.EncodeToString(uuid), nil
}

func RunTests(w http.ResponseWriter, r *http.Request) (int, string) {
	uuid, err := uuid()
	if err != nil {
		panic(err)
	}

    dir := fmt.Sprintf("test_runs/%s", uuid);
	os.MkdirAll(dir, os.ModeDir | os.ModePerm | os.ModeTemporary)

	apk, _, _ := r.FormFile("apk")
    if apk != nil {
        f := fmt.Sprintf("%s/main.apk", dir)
        apk_file, err := os.Create(f);
        defer apk_file.Close()

        if err != nil {
            panic(err)
        }

        _, err = io.Copy(apk_file, apk)
        if err != nil {
            panic(err)
        }

		rats.Install(f)
    }

	test_apk, _, err := r.FormFile("test-apk")

    if err != nil {
        panic("A Test Apk must be supplied")
    }

    f := fmt.Sprintf("%s/test.apk", dir)
    test_apk_file, err := os.Create(f);
    defer test_apk_file.Close()

    if err != nil {
        panic(err)
    }

    _, err = io.Copy(test_apk_file, test_apk)
    if err != nil {
        panic(err)
    }
	rats.Install(f)
    manifest := rats.GetManifest(f)

	s := test.RunTests(manifest)
	rats.Uninstall(manifest.Package)
	rats.Uninstall(manifest.Instrument.Target)

	str, err := json.Marshal(s)
	if err != nil {
        panic(err)
	}
    return http.StatusOK, string(str)
}

func GetDevices(parms martini.Params) (int, string) {
	b, _ := json.Marshal(<-rats.GetDevices())
	return http.StatusOK, string(b)
}

func serveStatic(m *martini.Martini) {
	_, file, _, _ := runtime.Caller(0)
	here := filepath.Dir(file)
	static := filepath.Join(here, "/public")
	m.Use(martini.Static(string(static)))
}

func main() {
	go rats.UpdateAdb(5)

	m := martini.New()
	m.Use(martini.Recovery())
	m.Use(martini.Logger())
	serveStatic(m)
	r := martini.NewRouter()
	r.Get(`/api/devices`, GetDevices)
	r.Post("/api/run", RunTests)
	m.Action(r.Handle)
	m.Run()
}
