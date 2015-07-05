package test

import (
	"encoding/xml"
	"github.com/wmbest2/rats/db"
	"time"
)

type StackType int64

const (
	Error StackType = iota
	Failure
)

type Property struct {
	XMLName xml.Name      `xml:"property" json:"-" bson:"-"`
	Id      int64         `xml:"-" json:"id"`
	Parent  int64         `xml:"-" json:"parent_id"`
	Name    string        `xml:"name,attr" json:"name"`
	Value   db.NullString `xml:"value,attr" json:"value"`
}

type TestCase struct {
	XMLName    xml.Name      `xml:"testcase" json:"-" bson:"-"`
	Id         int64         `xml:"-" json:"id"`
	Parent     int64         `xml:"-" json:"id"`
	Classname  string        `xml:"classname,attr" json:"classname"`
	Name       string        `xml:"name,attr" json:"name"`
	Status     db.NullString `xml:"status,attr,omitempty" json:"status"`
	Assertions db.NullString `xml:"assertions,attr,omitempty" json:"assertions"`
	Time       float64       `xml:"time,attr" json:"time"`
	Failures   []string      `xml:"failure" json:"failures,omitempty" bson:"failures,omitempty"`
	Errors     []string      `xml:"error" json:"errors,omitempty" bson:"errors,omitempty"`
	Skipped    bool          `xml:"skipped,omitempty" json:"skipped,omitempty" bson:"skipped,omitempty"`
	Stack      string        `xml:"-" json:"-" bson:"-"`
}

type TestSuite struct {
	XMLName    xml.Name      `xml:"testsuite" json:"-" bson:"-"`
	Id         int64         `xml:"-" json:"id"`
	Parent     int64         `xml:"-" json:"parent_id"`
	Properties []Property    `xml:"properties>property,omitempty" json:"properties, omitempty"`
	Tests      int           `xml:"tests,attr" json:"tests"`
	Failures   int           `xml:"failures,attr,omitempty" json:"failures"`
	Errors     int           `xml:"errors,attr,omitempty" json:"errors"`
	Skipped    int           `xml:"skipped,attr,omitempty" json:"skipped"`
	Hostname   db.NullString `xml:"hostname,attr,omitempty" json:"host"`
	Time       float64       `xml:"time,attr" json:"time"`
	Name       db.NullString `xml:"name,attr" json:"name"`
	SystemOut  db.NullString `xml:"system-out,omitempty" json:"system-out,omitempty"`
	SystemErr  db.NullString `xml:"system-err,omitempty" json:"system-out,omitempty"`
	TestCases  []TestCase    `xml:"testcase" json:"cases"`
}

type Artifact struct {
	Id    int64  `xml:"-" json:"id"`
	RunId int64  `xml:"-" json:"run_id"`
	Data  []byte `xml:"-" json:"data"`
}

type TestRun struct {
	XMLName    xml.Name       `xml:"testsuites" json:"-" bson:"-"`
	TestSuites []TestSuite    `xml:"testsuite" json:"suites,omitempty"`
	Name       string         `xml:"name,attr" json:"name"`
	Time       db.NullFloat64 `xml:"time,attr" json:"time"`

	Id        int64         `xml:"-" json:"id"`
	TokenId   int64         `xml:"-" json:"token_id"`
	ProjectId int64         `xml:"-" json:"project_id"`
	Timestamp time.Time     `xml:"-" json:"timestamp"`
	CommitId  db.NullString `xml:"-" json:"commit"`
	Message   db.NullString `xml:"-" json:"description"`
	Artifacts []Artifact    `xml:"-" json:"artifacts"`
	Success   db.NullBool   `xml:"-" json:"success"`
}
