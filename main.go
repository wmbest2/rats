package main

import (
    "github.com/go-martini/martini"
    "github.com/wmbest2/adb"
    "net/http"
    "encoding/json"
)

func GetDevices(parms martini.Params) (int, string) {
    deviceLock.Lock()
    b, _ := json.Marshal(devices)
    deviceLock.Unlock()
    return http.StatusOK, string(b)
}


func main() {
    go updateAdb(5)

    m := martini.Classic()
    r := martini.NewRouter()
    r.Get(`/api/devices`, GetDevices)
    m.Action(r.Handle)
    m.Run()
}
