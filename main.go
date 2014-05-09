package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/wmbest2/rats_server/rats"
	"github.com/wmbest2/rats_server/test"
	"io"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

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

func Mongo(db string) martini.Handler {
	session, err := mgo.Dial(db)
	if err != nil {
		panic(err)
	}

	return func(c martini.Context) {
		reqSession := session.Clone()
		c.Map(reqSession.DB("rats"))
		defer reqSession.Close()

		c.Next()
	}
}

func makeTestFolder() (string, string) {
	uuid, err := uuid()
	if err != nil {
		panic(err)
	}

	dir := fmt.Sprintf("test_runs/%s", uuid)
	os.MkdirAll(dir, os.ModeDir|os.ModePerm|os.ModeTemporary)
	return uuid, dir
}

func save(key string, filename string, r *http.Request) (bool, error) {
	apk, _, _ := r.FormFile(key)
	if apk != nil {
		apk_file, err := os.Create(filename)

		if err != nil {
			return false, err
		}

		_, err = io.Copy(apk_file, apk)
		apk.Close()
		apk_file.Close()

		if err != nil {
			return false, err
		}

		return true, nil
	}
	return false, nil
}

func RunTests(w http.ResponseWriter, r *http.Request, db *mgo.Database) (int, []byte) {
	uuid, dir := makeTestFolder()

	main := fmt.Sprintf("%s/main.apk", dir)
	install, err := save("apk", main, r)

	if err != nil {
		panic(err)
	}

    f := fmt.Sprintf("%s/test.apk", dir)
	_, err = save("test-apk", f, r)

	if err != nil {
		panic("A Test Apk must be supplied")
	}

	count, _ := strconv.Atoi(r.FormValue("count"))
	serialList := r.FormValue("serials")
    strict := r.FormValue("strict")

    var serials []string
    if serialList != "" {
        serials = strings.Split(serialList, ",")
    }

    filter := &rats.DeviceFilter{
        Count: count, 
        Strict: strict == "true",
    }
    filter.Serials = serials

	manifest := rats.GetManifest(f)
    filter.MinSdk = manifest.Sdk.Min
    filter.MaxSdk = manifest.Sdk.Max

	devices := <-rats.GetDevices(filter)
	rats.Reserve(devices)

	if install {
		rats.Install(main, devices)
	}

	rats.Install(f, devices)

	rats.Unlock(devices)

	s := test.RunTests(manifest, devices)
	s.Name = uuid
	s.Timestamp = time.Now()
	s.Project = manifest.Instrument.Target

	if dbErr := db.C("runs").Insert(&s); dbErr != nil {
		return http.StatusConflict, []byte(dbErr.Error())
	}

    go func() {
        rats.Uninstall(manifest.Package, devices)
        rats.Uninstall(manifest.Instrument.Target, devices)

        rats.Release(devices)
    }()

	os.RemoveAll(dir)

	str, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

    code := http.StatusOK
    if !s.Success {
        code = http.StatusInternalServerError
    }

	return code, str
}

func GetRunDevice(r *http.Request, parms martini.Params, db *mgo.Database) (int, []byte) {
	var runs test.TestSuites
    q := bson.M{"name": parms["id"]}
    s := bson.M{ "testsuites": bson.M{"$elemMatch": bson.M{"hostname": parms["device"]}}}
    query := db.C("runs").Find(q).Select(s).Limit(1)
	query.One(&runs)
	b, _ := json.Marshal(runs.TestSuites[0])
	return http.StatusOK, b
}

func GetRun(r *http.Request, parms martini.Params, db *mgo.Database) (int, []byte) {
	var runs test.TestSuites
	query := db.C("runs").Find(bson.M{"name": parms["id"]}).Limit(1)
	query.One(&runs)
	b, _ := json.Marshal(runs)
	return http.StatusOK, b
}

func GetRuns(r *http.Request, parms martini.Params, db *mgo.Database) (int, []byte) {
	page := 0
	p := r.URL.Query().Get("page")
	if p != "" {
		page, _ = strconv.Atoi(p)
	}

	size := 25
	s := r.URL.Query().Get("count")
	if s != "" {
		size, _ = strconv.Atoi(s)
	}

	var runs []*test.TestSuites
	query := db.C("runs").Find(bson.M{}).Limit(size).Skip(page * size)
	query.Select(bson.M{"testsuites.testcases": 0, "testsuites.device.inuse": 0})
	query.Sort("-timestamp")
	query.All(&runs)
	b, _ := json.Marshal(runs)
	return http.StatusOK, b
}

func GetDevices(parms martini.Params) (int, []byte) {
	b, _ := json.Marshal(<-rats.GetAllDevices())
	return http.StatusOK, b
}

func serveStatic(m *martini.Martini) {
	_, file, _, _ := runtime.Caller(0)
	here := filepath.Dir(file)
	static := filepath.Join(here, "/public")
	m.Use(martini.Static(string(static)))
	m.Use(martini.Static("public"))
}

func main() {
	go rats.UpdateAdb(5)

    var mongodb = flag.String("db", "mongodb://localhost/rats", "Mongo db url")
    var port = flag.Int("port", 3000, "Port to serve")
    var debug = flag.Bool("debug", false, "Log debug information")

    flag.Parse()

	m := martini.New()
	m.Use(martini.Recovery())
	m.Use(Mongo(*mongodb))
	serveStatic(m)

    if *debug {
        m.Use(martini.Logger())
    }

	r := martini.NewRouter()
	r.Get(`/api/devices`, GetDevices)
	r.Post("/api/run", RunTests)
	r.Get("/api/runs", GetRuns)
	r.Get("/api/runs/:id", GetRun)
	r.Get("/api/runs/:id/:device", GetRunDevice)

	m.Action(r.Handle)
    fmt.Printf("Listening on port %d\n", *port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), m))
}
