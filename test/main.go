package test

import (
	"encoding/xml"
	"time"
)

type Property struct {
	XMLName xml.Name `xml:"property" json:"-" bson:"-"`
	Name    string   `xml:"name,attr" json:"name"`
	Value   string   `xml:"value,attr" json:"value"`
}

type TestCase struct {
	XMLName    xml.Name `xml:"testcase" json:"-" bson:"-"`
	Classname  string   `xml:"classname,attr" json:"classname"`
	Name       string   `xml:"name,attr" json:"name"`
	Status     string   `xml:"status,attr,omitempty" json:"status"`
	Assertions string   `xml:"assertions,attr,omitempty" json:"assertions"`
	Time       float64  `xml:"time,attr" json:"time"`
	Failures   []string `xml:"failure" json:"failures,omitempty" bson:"failures,omitempty"`
	Errors     []string `xml:"error" json:"errors,omitempty" bson:"errors,omitempty"`
	Skipped    bool     `xml:"skipped,omitempty" json:"skipped,omitempty" bson:"skipped,omitempty"`
	Stack      string   `xml:"-" json:"-" bson:"-"`
}

type TestSuite struct {
	XMLName    xml.Name   `xml:"testsuite" json:"-" bson:"-"`
	Properties []Property `xml:"properties>property,omitempty" json:"properties, omitempty"`
	Tests      int        `xml:"tests,attr" json:"tests"`
	Failures   int        `xml:"failures,attr,omitempty" json:"failures"`
	Errors     int        `xml:"errors,attr,omitempty" json:"errors"`
	Skipped    int        `xml:"skipped,attr,omitempty" json:"skipped"`
	Hostname   string     `xml:"hostname,attr,omitempty" json:"host"`
	Time       float64    `xml:"time,attr" json:"time"`
	Name       string     `xml:"name,attr" json:"name"`
	SystemOut  string     `xml:"system-out,omitempty" json:"system-out,omitempty"`
	SystemErr  string     `xml:"system-err,omitempty" json:"system-out,omitempty"`
	TestCases  []TestCase `xml:"testcase" json:"cases"`
}

type TestSuites struct {
	XMLName    xml.Name    `xml:"testsuites" json:"-" bson:"-"`
	TestSuites []TestSuite `xml:"testsuite" json:"suites,omitempty"`
	Name       string      `xml:"name,attr" json:"name"`
	Time       float64     `xml:"time,attr" json:"time"`

	// Rats specifics? Should this be a separate TestRun object?
	Project   int64     `xml:"-" json:"project"`
	Timestamp time.Time `xml:"-" json:"timestamp"`
	Message   string    `xml:"-" json:"description"`
	Success   bool      `xml:"-" json:"success"`
}
