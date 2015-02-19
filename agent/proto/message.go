package proto

import "github.com/docker/libchan"

const (
	Register = "rats:proto:register"
	Init     = "rats:proto:init"
	Complete = "rats:proto:complete"
)

type Run struct {
	Binary   map[string][]byte
	Metadata map[string]string
}

type Message struct {
	Command   string
	Run       *Run
	Result    []byte
	Responder libchan.Sender
}

func NewRun(s libchan.Sender) *Message {
	return &Message{
		Command: Init,
		Run: &Run{
			Binary:   make(map[string][]byte),
			Metadata: make(map[string]string),
		},
		Responder: s,
	}
}
