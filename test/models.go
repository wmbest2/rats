package test

import (
	"encoding/xml"
    "time"
    "github.com/wmbest2/rats_server/rats"
)

type TestCase struct {
	XMLName   xml.Name `xml:"testcase" json:"-"`
	Classname string   `xml:"classname,attr" json:"classname"`
	Name      string   `xml:"name,attr" json:"name"`
	Time      float64  `xml:"time,attr" json:"time"`
	Failure   *string   `xml:"failure,omitempty" json:"failure,omitempty"`
	Error     *string   `xml:"error,omitempty" json:"error,omitempty"`
	Stack     string   `xml:"-" json:"-"`
}

type TestSuite struct {
	XMLName   xml.Name    `xml:"testsuite" json:"-"`
	Tests     int         `xml:"tests,attr" json:"tests"`
	Failures  int         `xml:"failures,attr" json:"failures"`
	Errors    int         `xml:"errors,attr" json:"errors"`
	Hostname  string      `xml:"hostname,attr" json:"host"`
    Time      float64     `xml:"time,attr" json:"time"`
	Name      string      `xml:"name,attr" json:"name"`
    Device    *rats.Device `xml:"-" json:"device,omitempty"`
	TestCases []*TestCase `json:"cases"`
}

type TestSuites struct {
	XMLName    xml.Name     `xml:"testsuites" json:"-"`
	TestSuites []*TestSuite `json:"suites,omitempty"`
    Name       string       `xml:"name,attr" json:"name"`
    Project    string       `json:"project"`
    Timestamp  time.Time    `json:"timestamp"`
    Time       float64      `xml:"time,attr" json:"time"`
    Success    bool         `json:"success"`
}
