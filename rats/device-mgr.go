package rats

import (
	"github.com/wmbest2/android/adb"
	"sync"
	"time"
)

var devices map[string]*Device
var lock sync.Mutex

type Device struct {
	adb.Device
	InUse bool
}

type DeviceFilter struct {
	adb.DeviceFilter
	Count  int
	Strict bool
}

func Poll(in chan map[string]*Device, out chan map[string]*Device) {
	devices := <-in
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
	in := make(chan map[string]*Device)
	out := make(chan map[string]*Device)
	go Poll(in, out)
	lock.Lock()
	in <- devices
	devices = <-out
	lock.Unlock()
}

func UpdateAdb(seconds time.Duration) {
	PollDevices()

	c := time.Tick(seconds * time.Second)
	for _ = range c {
		PollDevices()
	}
}

func GetAllDevices() chan []*Device {
    return GetDevices(nil)
}

func GetDevices(filter *DeviceFilter) chan []*Device {
	out := make(chan []*Device)

	go func() {
		lock.Lock()
		v := make([]*Device, 0, len(devices))
		lock.Unlock()

        count := -1
		if filter != nil && filter.Count > 0 {
			count = filter.Count
		}
		for {
			lock.Lock()
			for _, value := range devices {
				if (filter == nil || (value.MatchFilter(&filter.DeviceFilter)) && !value.InUse) {
					v = append(v, value)
					if count > 1 {
						count--
					} else if count != -1 {
						break
					}
				}
			}
			lock.Unlock()

			if filter == nil || !filter.Strict || count == 0  {
				break
			}

            <-time.After(5 * time.Second)
		}

		out <- v
	}()
	return out
}

func Reserve(devices []*Device) {
	lock.Lock()
	for _, value := range devices {
		value.InUse = true
	}
	lock.Unlock()
}

func Release(devices []*Device) {
	lock.Lock()
	for _, value := range devices {
		value.InUse = false
	}
	lock.Unlock()
}
