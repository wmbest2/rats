package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"github.com/wmbest2/rats_server/rats"
	"net/http"
    "runtime"
    "path/filepath"
)

func RunTests(w http.ResponseWriter, r *http.Request) {
    //apk, header,_ := r.FormFile("apk")
    //test_apk, test_header, err := r.FormFile("test-apk")

    //if err != nil {
        //panic(err)
    //}
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
