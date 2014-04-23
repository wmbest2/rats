package rats

import (
	"github.com/wmbest2/android/adb"
	"sync"
	"time"
)

var Devices []*adb.Device
var DeviceLock sync.Mutex

func UpdateDevices() {
	DeviceLock.Lock()
	Devices = adb.ListDevices(nil)
	DeviceLock.Unlock()
}

func UpdateAdb(seconds time.Duration) {
	UpdateDevices()

	c := time.Tick(seconds * time.Second)
	for _ = range c {
		UpdateDevices()
	}
}
