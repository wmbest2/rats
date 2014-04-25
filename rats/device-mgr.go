package rats

import (
	"github.com/wmbest2/android/adb"
	"time"
)

var devices map[string]*adb.Device

func PollDevices() {
    out := make(chan []*adb.Device)

    go func() {
        out <- adb.ListDevices(nil)
    }()

    new_devices := <-out
    new_map := make(map[string]*adb.Device)
    for _, d := range new_devices {
        if devices[d.String()] != nil {
             new_map[d.String()] = devices[d.String()]
        } else {
             new_map[d.String()] = d
        }
    }
    devices = new_map
}

func UpdateAdb(seconds time.Duration) {
	PollDevices()

	c := time.Tick(seconds * time.Second)
	for _ = range c {
		PollDevices()
	}
}

func GetDevices() chan []*adb.Device {
    out := make(chan []*adb.Device)
    
    go func() {
        v := make([]*adb.Device, 0, len(devices))

        for  _, value := range devices {
           v = append(v, value)
        }
        out <- v 
    }()
    return out
}
