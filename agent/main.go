package agent

import (
	"github.com/wmbest2/rats-server/rats/device"
)

type Agent interface {
	List() []Device
	Reserve(d ...Device) bool
	Release(d ...Device) bool
	RunTest(devices []Device, app io.Reader, test io.Reader)
}
