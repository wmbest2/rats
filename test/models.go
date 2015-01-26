package test

import (
	"encoding/xml"
	"github.com/wmbest2/rats-server/rats"
	"time"
)

type Issue struct {
	XMLName   xml.Name `json:"-" bson:"-"`
	Exception *string  `xml:",chardata" json:",string"`
}

type TestCase struct {
	XMLName    xml.Name `xml:"testcase" json:"-" bson:"-"`
	Classname  string   `xml:"classname,attr" json:"classname"`
	Name       string   `xml:"name,attr" json:"name"`
	Status     string   `xml:"status,attr" json:"status"`
	Assertions string   `xml:"assertions,attr" json:"assertions"`
	Time       float64  `xml:"time,attr" json:"time"`
	Failures   []Issue  `json:"failure,omitempty" bson:"failures,omitempty"`
	Errors     []Issue  `json:"error,omitempty" bson:"errors,omitempty"`
	Skipped    bool     `xml:"skipped,omitempty" json:"skipped,omitempty" bson:"skipped,omitempty"`
	Stack      string   `xml:"-" json:"-" bson:"-"`
}

type TestSuite struct {
	XMLName   xml.Name     `xml:"testsuite" json:"-" bson:"-"`
	Tests     int          `xml:"tests,attr" json:"tests"`
	Failures  int          `xml:"failures,attr" json:"failures"`
	Errors    int          `xml:"errors,attr" json:"errors"`
	Skipped   int          `xml:"skipped,attr" json:"skipped"`
	Hostname  string       `xml:"hostname,attr" json:"host"`
	Time      float64      `xml:"time,attr" json:"time"`
	Name      string       `xml:"name,attr" json:"name"`
	Device    *rats.Device `xml:"-" json:"device,omitempty" "device,omitempty"`
	TestCases []*TestCase  `json:"cases"`
}

type TestSuites struct {
	XMLName    xml.Name     `xml:"testsuites" json:"-" bson:"-"`
	TestSuites []*TestSuite `json:"suites,omitempty"`
	Name       string       `xml:"name,attr" json:"name"`
	Project    string       `json:"project"`
	Timestamp  time.Time    `json:"timestamp"`
	Time       float64      `xml:"time,attr" json:"time"`
	Message    string       `json:"description"`
	Success    bool         `json:"success"`
}
