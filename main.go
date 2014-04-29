package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/wmbest2/rats_server/rats"
	"github.com/wmbest2/rats_server/test"
	"io"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
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

func RunTests(w http.ResponseWriter, r *http.Request, db *mgo.Database) (int, string) {
	uuid, err := uuid()
	if err != nil {
		panic(err)
	}

	dir := fmt.Sprintf("test_runs/%s", uuid)
	os.MkdirAll(dir, os.ModeDir|os.ModePerm|os.ModeTemporary)

	apk, _, _ := r.FormFile("apk")
	if apk != nil {
		defer apk.Close()
		f := fmt.Sprintf("%s/main.apk", dir)
		apk_file, err := os.Create(f)
		defer apk_file.Close()

		if err != nil {
			panic(err)
		}

		_, err = io.Copy(apk_file, apk)
		if err != nil {
			panic(err)
		}

		rats.Install(f)
	}

	test_apk, _, err := r.FormFile("test-apk")

	if err != nil {
		panic("A Test Apk must be supplied")
	}

	defer test_apk.Close()

	f := fmt.Sprintf("%s/test.apk", dir)
	test_apk_file, err := os.Create(f)
	defer test_apk_file.Close()

	if err != nil {
		panic(err)
	}

	_, err = io.Copy(test_apk_file, test_apk)
	if err != nil {
		panic(err)
	}
	rats.Install(f)
	manifest := rats.GetManifest(f)

	s := test.RunTests(manifest)

	s.Name = uuid
	s.Timestamp = time.Now()
	s.Project = manifest.Instrument.Target

	rats.Uninstall(manifest.Package)
	rats.Uninstall(manifest.Instrument.Target)
	os.RemoveAll(dir)

	if dbErr := db.C("runs").Insert(&s); dbErr != nil {
		return http.StatusConflict, string(dbErr.Error())
	}

	str, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return http.StatusOK, string(str)
}

func GetRunDevice(r *http.Request, parms martini.Params, db *mgo.Database) (int, string) {
	var runs test.TestSuites
    q := bson.M{"name": parms["id"]}
    fmt.Printf("%#v\n", q)
    query := db.C("runs").Find(q).Limit(1)
	query.One(&runs)
    for _, run := range runs.TestSuites {
        if run.Hostname == parms["device"] {
            b, _ := json.Marshal(run)
            return http.StatusOK, string(b)
        }
    }
    return http.StatusNotFound, fmt.Sprintf("Run on Device %s Not Found", parms["device"])
}

func GetRun(r *http.Request, parms martini.Params, db *mgo.Database) (int, string) {
	var runs test.TestSuites
    query := db.C("runs").Find(bson.M{"name": parms["id"]}).Limit(1)
	query.One(&runs)
	b, _ := json.Marshal(runs)
	return http.StatusOK, string(b)
}

func GetRuns(r *http.Request, parms martini.Params, db *mgo.Database) (int, string) {
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
	//query.Select(bson.M{"name": 1, "project": 1, "timestamp": 1, "time": 1, "success": 1})
	query.Sort("-timestamp")
	query.All(&runs)
	b, _ := json.Marshal(runs)
	return http.StatusOK, string(b)
}

func GetDevices(parms martini.Params) (int, string) {
	b, _ := json.Marshal(<-rats.GetDevices())
	return http.StatusOK, string(b)
}

func serveStatic(m *martini.Martini) {
	_, file, _, _ := runtime.Caller(0)
	here := filepath.Dir(file)
	static := filepath.Join(here, "/public")
	m.Use(martini.Static(string(static)))
}

func main() {
	go rats.UpdateAdb(5)

	m := martini.New()
	m.Use(martini.Recovery())
	m.Use(martini.Logger())
	m.Use(Mongo("localhost/rats"))
	serveStatic(m)
	r := martini.NewRouter()
	r.Get(`/api/devices`, GetDevices)
	r.Post("/api/run", RunTests)
	r.Get("/api/runs", GetRuns)
    r.Get("/api/runs/:id", GetRun)
    r.Get("/api/runs/:id/:device", GetRunDevice)

	r.Get("/debug/pprof", pprof.Index)
	r.Get("/debug/pprof/cmdline", pprof.Cmdline)
	r.Get("/debug/pprof/profile", pprof.Profile)
	r.Get("/debug/pprof/symbol", pprof.Symbol)
	r.Post("/debug/pprof/symbol", pprof.Symbol)
	r.Get("/debug/pprof/block", pprof.Handler("block").ServeHTTP)
	r.Get("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
	r.Get("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
	r.Get("/debug/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)

	m.Action(r.Handle)
	m.Run()
}
