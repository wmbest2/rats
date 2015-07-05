package device

import (
	"github.com/wmbest2/rats/test"
	"io"
)

type Metadata map[string]string

type Device interface {
	Serial() string
	Metadata() *Metadata

	Reserve()
	Release()
	InUse() bool

	RunTest(app io.Reader, test io.Reader) (chan test.TestSuite, chan bool)
}
