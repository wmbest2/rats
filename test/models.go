package test

import (
	"encoding/xml"
)

type Error struct {
	Contents string `xml:",chardata"`
}

type Failure struct {
	Contents string `xml:",chardata"`
}

type TestCase struct {
	XMLName   xml.Name `xml:"testcase"`
	Classname string   `xml:"classname,attr"`
	Name      string   `xml:"name,attr"`
	Time      float64  `xml:"time,attr"`
	Failure   *Failure `xml:"failure,omitempty"`
	Error     *Error   `xml:"error,omitempty"`
    Stack     string   `xml:"-"`
}

type TestSuite struct {
	XMLName   xml.Name `xml:"testsuite"`
	Tests     int      `xml:"tests,attr"`
	Failures  int      `xml:"failures,attr"`
	Errors    int      `xml:"errors,attr"`
	Hostname  string   `xml:"hostname,attr"`
	Time      float64  `xml:"time,attr"`
	Name      string   `xml:"name,attr"`
	TestCases []*TestCase
}

type TestSuites struct {
	XMLName   xml.Name `xml:"testsuites"`
	TestSuites []*TestSuite 
	Time       float64   `xml:"time,attr"`
}
