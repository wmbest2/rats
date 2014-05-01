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

    count,_ := strconv.Atoi(r.FormValue("count"))
    
    filter := &rats.DeviceFilter{Count: count}

	f := fmt.Sprintf("%s/main.apk", dir)
	install, err := save("apk", f, r)

	if err != nil {
		panic(err)
	}

    devices := <-rats.GetDevices(filter)
    rats.Reserve(devices)

	if install {
		rats.Install(f, devices)
	}

	f = fmt.Sprintf("%s/test.apk", dir)
	_, err = save("test-apk", f, r)

	if err != nil {
		panic("A Test Apk must be supplied")
	}

	rats.Install(f, devices)
	manifest := rats.GetManifest(f)

	rats.Unlock(devices)

	s := test.RunTests(manifest, devices)
	s.Name = uuid
	s.Timestamp = time.Now()
	s.Project = manifest.Instrument.Target

	if dbErr := db.C("runs").Insert(&s); dbErr != nil {
		return http.StatusConflict, []byte(dbErr.Error())
	}

	rats.Uninstall(manifest.Package, devices)
	rats.Uninstall(manifest.Instrument.Target, devices)
    rats.Release(devices)

	os.RemoveAll(dir)

	str, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	return http.StatusOK, str
}

func GetRunDevice(r *http.Request, parms martini.Params, db *mgo.Database) (int, []byte) {
	var runs test.TestSuites
    q := bson.M{"name": parms["id"], "testsuites.hostname": parms["device"]}
	fmt.Printf("%#v\n", q)
	query := db.C("runs").Find(q).Limit(1)
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

	m.Action(r.Handle)
	m.Run()
}
