package agent

import (
	"github.com/wmbest2/rats/device"
	"io"
)

type Agent interface {
	List() []device.Device
	Reserve(d ...device.Device) bool
	Release(d ...device.Device) bool
	RunTest(devices []device.Device, app io.Reader, test io.Reader)
}
