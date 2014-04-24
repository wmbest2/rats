package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"github.com/wmbest2/rats_server/rats"
	"net/http"
    "runtime"
    "path/filepath"
)

func GetDevices(parms martini.Params) (int, string) {
	rats.DeviceLock.Lock()
	b, _ := json.Marshal(rats.Devices)
	rats.DeviceLock.Unlock()
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

	m := martini.Classic()
	r := martini.NewRouter()
	r.Get(`/api/devices`, GetDevices)
	m.Action(r.Handle)
	m.Run()
}
