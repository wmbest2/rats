package test

import (
	"encoding/xml"
)

type Error struct {
	Contents string `xml:",chardata" json:"data"`
}

type Failure struct {
    Contents string `xml:",chardata" json:"data"`
}

type TestCase struct {
    XMLName   xml.Name `xml:"testcase" json:"-"`
    Classname string   `xml:"classname,attr" json:"classname"`
    Name      string   `xml:"name,attr" json:"name"`
    Time      float64  `xml:"time,attr" json:"time"`
    Failure   *Failure `xml:"failure,omitempty" json:"failure"`
    Error     *Error   `xml:"error,omitempty" json:"error"`
    Stack     string   `xml:"-" json:"-"`
}

type TestSuite struct {
    XMLName   xml.Name `xml:"testsuite" json:"-"`
    Tests     int      `xml:"tests,attr" json:"tests"`
    Failures  int      `xml:"failures,attr" json:"failures"`
    Errors    int      `xml:"errors,attr" json:"errors"`
    Hostname  string   `xml:"hostname,attr" json:"host"`
	Time      float64  `xml:"time,attr" json"time"`
	Name      string   `xml:"name,attr" json"name"`
	TestCases []*TestCase `json:"cases"`
}

type TestSuites struct {
    XMLName    xml.Name `xml:"testsuites" json:"-"`
	TestSuites []*TestSuite `json:"suites"`
	Time       float64 `xml:"time,attr"`
}
