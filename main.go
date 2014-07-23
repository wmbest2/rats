package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"github.com/wmbest2/android/adb"
	"github.com/wmbest2/rats-server/rats"
	"github.com/wmbest2/rats-server/test"
	"io"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	mgoSession *mgo.Session

	mongodb = flag.String("db", "mongodb://localhost/rats", "Mongo db url")
	port    = flag.Int("port", 3000, "Port to serve")
	debug   = flag.Bool("debug", false, "Log debug information")
)

type RatsHandler func(http.ResponseWriter, *http.Request, *mgo.Database) error

func (rh RatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s := mgoSession.Clone()
	defer s.Close()

	err := rh(w, r, s.DB("rats"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func uuid() (string, error) {
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// TODO: verify the two lines implement RFC 4122 correctly
	uuid[8] = 0x80 // variant bits see page 5
	uuid[4] = 0x40 // version 4 Pseudo Random, see page 7

	return hex.EncodeToString(uuid), nil
}

type byteCounter struct {
	Total int64
}

func (b *byteCounter) Write(p []byte) (int, error) {
	b.Total += int64(len(p))
	return len(p), nil
}

func getLength(reader io.ReadSeeker) int64 {
	reader.Seek(0, 0) // reset for new reader
	var counter byteCounter

	t := io.TeeReader(reader, &counter)

	buf := make([]byte, 4096)
	for {
		_, err := t.Read(buf)
		if err != nil {
			break
		}
	}

	reader.Seek(0, 0) // reset
	return counter.Total
}

func RunTests(w http.ResponseWriter, r *http.Request, db *mgo.Database) error {
	uuid, err := uuid()
	if err != nil {
		return err
	}
	main, _, err := r.FormFile("apk")
	testApk, _, err := r.FormFile("test-apk")
	if err != nil {
		return errors.New("A Test Apk must be supplied")
	}

	size := getLength(testApk)

	count, _ := strconv.Atoi(r.FormValue("count"))
	serialList := r.FormValue("serials")
	strict := r.FormValue("strict")

	var serials []string
	if serialList != "" {
		serials = strings.Split(serialList, ",")
	}

	filter := &rats.DeviceFilter{
		Count:  count,
		Strict: strict == "true",
	}
	filter.Serials = serials

	manifest := rats.GetManifest(testApk, size)
	filter.MinSdk = manifest.Sdk.Min
	filter.MaxSdk = manifest.Sdk.Max

	testApk.Seek(0, 0) // reset for new reader

	devices := <-rats.GetDevices(filter)
	rats.Reserve(devices...)

	if main != nil {
		rats.Install("main.apk", main, devices...)
	}

	rats.Install("test.apk", testApk, devices...)

	rats.Unlock(devices)

	finished, out := test.RunTests(manifest, devices)

	var s *test.TestSuites
SuitesLoop:
	for {
		select {
		case s = <-out:
			break SuitesLoop
		case device := <-finished:
			go func() {
				rats.Uninstall(manifest.Package, device)
				rats.Uninstall(manifest.Instrument.Target, device)

				rats.Release(device)
			}()
		}
	}

	s.Name = uuid
	s.Timestamp = time.Now()
	s.Project = manifest.Instrument.Target

	if dbErr := db.C("runs").Insert(&s); dbErr != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(dbErr.Error())
	}

	if !s.Success {
		w.WriteHeader(http.StatusInternalServerError)
	}

	return json.NewEncoder(w).Encode(s)
}

func GetRunDevice(w http.ResponseWriter, r *http.Request, db *mgo.Database) error {
	vars := mux.Vars(r)

	var runs test.TestSuites
	q := bson.M{"name": vars["id"]}
	s := bson.M{"testsuites": bson.M{"$elemMatch": bson.M{"hostname": vars["device"]}}}
	query := db.C("runs").Find(q).Select(s).Limit(1)
	query.One(&runs)

	return json.NewEncoder(w).Encode(runs.TestSuites[0])
}

func GetRun(w http.ResponseWriter, r *http.Request, db *mgo.Database) error {
	vars := mux.Vars(r)

	var runs test.TestSuites
	query := db.C("runs").Find(bson.M{"name": vars["id"]}).Limit(1)
	query.One(&runs)

	return json.NewEncoder(w).Encode(runs)
}

func GetRuns(w http.ResponseWriter, r *http.Request, db *mgo.Database) error {
	page := 0
	p := r.URL.Query().Get("page")
	if p != "" {
		page, _ = strconv.Atoi(p)
	}

	size := 10
	s := r.URL.Query().Get("count")
	if s != "" {
		size, _ = strconv.Atoi(s)
	}

	var runs []*test.TestSuites
	query := db.C("runs").Find(bson.M{}).Limit(size).Skip(page * size)
	query.Select(bson.M{"testsuites.testcases": 0, "testsuites.device.inuse": 0})
	query.Sort("-timestamp")
	query.All(&runs)

	return json.NewEncoder(w).Encode(runs)
}

func GetDevices(w http.ResponseWriter, r *http.Request, db *mgo.Database) error {
	return json.NewEncoder(w).Encode(<-rats.GetAllDevices())
}

func PingHandler(w http.ResponseWriter, r *http.Request, db *mgo.Database) error {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "pong")

	return nil
}

func init() {
	flag.Parse()

	var err error
	mgoSession, err = mgo.Dial(*mongodb)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	go rats.UpdateAdb(adb.Default)

	r := mux.NewRouter()

	r.Handle("/api/ping", RatsHandler(PingHandler))
	r.Handle("/api/devices", RatsHandler(GetDevices))
	r.Handle("/api/run", RatsHandler(RunTests))
	r.Handle("/api/runs", RatsHandler(GetRuns))
	r.Handle("/api/runs/{id}", RatsHandler(GetRun))
	r.Handle("/api/runs/{id}/{device}", RatsHandler(GetRunDevice))
	r.PathPrefix("/").Handler(http.FileServer(rice.MustFindBox(`public`).HTTPBox()))

	http.Handle("/", r)

	log.Printf("Listening on port %d\n", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatalf("Error starting server: %s\n", err.Error())
	}
}
