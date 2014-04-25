package rats

import (
	"github.com/wmbest2/android/adb"
	"time"
    "sync"
)

var devices map[string]*Device
var lock sync.Mutex

type Device struct {
	adb.Device
	InUse bool
}

func Poll(in chan map[string]*Device, out chan map[string]*Device) {
    devices := <- in
	new_devices := adb.ListDevices(nil)
	new_map := make(map[string]*Device)
	for _, d := range new_devices {
		if devices[d.String()] != nil {
			new_map[d.String()] = devices[d.String()]
		} else {
			new_map[d.String()] = &Device{Device: *d, InUse: false}
		}
	}
    out <- new_map
}

func PollDevices() {
    in := make(chan map[string]*Device);
    out := make(chan map[string]*Device);
    go Poll(in, out)
    lock.Lock()
    in <- devices
    devices = <- out
    lock.Unlock()
}

func UpdateAdb(seconds time.Duration) {
	PollDevices()

	c := time.Tick(seconds * time.Second)
	for _ = range c {
		PollDevices()
	}
}

func GetDevices() chan []*Device {
	out := make(chan []*Device)

    lock.Lock()
    v := make([]*Device, 0, len(devices))

    for _, value := range devices {
        v = append(v, value)
    }
    lock.Unlock()

	go func() {
		out <- v
	}()
	return out
}
