package main

import (
    "github.com/wmbest2/android/adb"
    "time"
    "sync"
)

var devices []*adb.Device
var deviceLock sync.Mutex

func updateDevices() {
    deviceLock.Lock()
    devices = adb.ListDevices(nil)
    deviceLock.Unlock()
}

func updateAdb(seconds time.Duration) {
    updateDevices()

    c := time.Tick(seconds * time.Second)
    for _ = range c {
        updateDevices()
    }
}
