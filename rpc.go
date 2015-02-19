package main

import (
	"fmt"
	"log"
	"net"

	"github.com/docker/libchan"
	"github.com/docker/libchan/spdy"
	"github.com/wmbest2/rats/agent/proto"
)

const (
	Android = "android"
)

var (
	daemon libchan.Sender
)

func readChan(receiver libchan.Receiver) {
	for {
		msg := &proto.Message{}
		err := receiver.Receive(msg)
		if err != nil {
			log.Print(err)
			break
		}

		switch msg.Command {
		case proto.Register:
			log.Println("Registering new daemon worker")
			daemon = msg.Responder
			err = daemon.Send(proto.Message{
				Command: "done",
			})
			if err != nil {
				log.Print(err)
				break
			}
		}
	}
}

func listenRpc() {
	log.Printf("RPC listening on port %d\n", *rpcport)
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *rpcport))
	if err != nil {
		log.Fatal(err)
	}

	tran, err := spdy.NewTransportListener(l, spdy.NoAuthenticator)
	if err != nil {
		log.Fatal(err)
	}

	for {
		t, err := tran.AcceptTransport()
		if err != nil {
			log.Print(err)
			break
		}

		go func() {
			for {
				receiver, err := t.WaitReceiveChannel()
				if err != nil {
					log.Print(err)
					break
				}

				go readChan(receiver)
			}
		}()
	}
}
