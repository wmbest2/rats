package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/docker/libchan"
	"github.com/gorilla/mux"
	"github.com/wmbest2/rats/agent/proto"
)

var (
	port    = flag.Int("port", 3000, "Port to serve HTTP connections")
	rpcport = flag.Int("rpcport", 4000, "Port to serve RPC connections")
	debug   = flag.Bool("debug", false, "Log debug information")
)

type RatsHandler func(http.ResponseWriter, *http.Request) error

func (rh RatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := rh(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type PageMeta struct {
	Page  int `json:"page"`
	Count int `json:"count"`
	Total int `json:"total"`
}

func RunTests(w http.ResponseWriter, r *http.Request) error {
	receiver, remoteSender := libchan.Pipe()
	msg := proto.NewRun(remoteSender)

	uuid, err := uuid()
	if err != nil {
		return err
	}

	msg.Run.Metadata["uuid"] = uuid

	main, _, err := r.FormFile("apk")
	testApk, _, err := r.FormFile("test-apk")
	if err != nil {
		return errors.New("A Test Apk must be supplied")
	}

	// TODO: this can just copy r.FormValues -> map[string][]string
	mainBuf := &bytes.Buffer{}
	testBuf := &bytes.Buffer{}

	_, err = mainBuf.ReadFrom(main)
	if err != nil {
		return err
	}

	_, err = testBuf.ReadFrom(testApk)
	if err != nil {
		return err
	}

	msg.Run.Binary["main"] = mainBuf.Bytes()
	msg.Run.Binary["test"] = testBuf.Bytes()

	msg.Run.Metadata["count"] = r.FormValue("count")

	msg.Run.Metadata["serialList"] = r.FormValue("serials")
	msg.Run.Metadata["strict"] = r.FormValue("strict")
	msg.Run.Metadata["msg"] = r.FormValue("message")

	//if dbErr := db.C("runs").Insert(&s); dbErr != nil {
	//w.WriteHeader(http.StatusConflict)
	//json.NewEncoder(w).Encode(dbErr.Error())
	//}

	err = daemon.Send(msg)
	if err != nil {
		return err
	}

	msg = &proto.Message{}
	err = receiver.Receive(msg)
	if err != nil {
		return err
	}

	if v, ok := msg.Result.([]byte); ok {
		fmt.Fprint(w, string(v))
		return nil
	}

	w.WriteHeader(http.StatusInternalServerError)
	return nil
}

func GetDevices(w http.ResponseWriter, r *http.Request) error {
	receiver, remoteSender := libchan.Pipe()

	err := daemon.Send(proto.Message{
		Command:   proto.Devices,
		Responder: remoteSender,
	})

	if err != nil {
		return err
	}

	msg := &proto.Message{}
	err = receiver.Receive(msg)
	if err != nil {
		return err
	}

	if v, ok := msg.Result.([]byte); ok {
		fmt.Fprint(w, string(v))
		return nil
	}

	w.WriteHeader(http.StatusInternalServerError)
	return nil
}

func PingHandler(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "pong")

	return nil
}

func init() {
	flag.Parse()

	r := mux.NewRouter()

	r.Handle("/api/ping", RatsHandler(PingHandler))
	r.Handle("/api/devices", RatsHandler(GetDevices))
	r.Handle("/api/run", RatsHandler(RunTests))
	//r.Handle("/api/runs", RatsHandler(GetRuns))
	//r.Handle("/api/runs/{id}", RatsHandler(GetRun))
	//r.Handle("/api/runs/{id}/{device}", RatsHandler(GetRunDevice))
	r.PathPrefix("/").Handler(http.FileServer(rice.MustFindBox(`public`).HTTPBox()))

	http.Handle("/", r)
}

func main() {
	go listenRpc()

	log.Printf("Listening on port %d\n", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatalf("Error starting server: %s\n", err.Error())
	}
}
