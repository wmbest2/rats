package test

import (
	"encoding/json"
	"encoding/xml"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"testing"

	"github.com/wmbest2/rats/db"
)

var (
	testcase = `<?xml version="1.0" encoding="UTF-8" ?>
<testsuites name="abcd" time="0.069">
	<testsuite errors="1" failures="1" hostname="test-macbook-pro.local"
		name="tests.ATest" tests="3" time="0.069" timestamp="2009-12-19T17:58:59">
	  <properties>
		<property name="java.vendor" value="Sun Microsystems Inc." />
		<property name="jtreg.home" value="/Users/mahmood/research/jtreg" />
		<property name="sun.java.launcher" value="SUN_STANDARD" />
		<property name="sun.management.compiler" value="HotSpot 64-Bit Server Compiler" />
		<property name="os.name" value="Darwin" />
		<property name="java.specification.version" value="1.6" />
	  </properties>
	  <testcase classname="tests.ATest" name="error" time="0.0060">
		<error type="java.lang.RuntimeException">java.lang.RuntimeException
		at tests.ATest.error(ATest.java:11)
	</error>
	  </testcase>
	  <testcase classname="tests.ATest" name="fail" time="0.0020">
		<failure type="junit.framework.AssertionFailedError">junit.framework.AssertionFailedError: 
		at tests.ATest.fail(ATest.java:9)
	</failure>
	  </testcase>
	  <testcase classname="tests.ATest" name="sucess" time="0.0" />
	  <system-out><![CDATA[here
	]]></system-out>
	  <system-err><![CDATA[]]></system-err>
	</testsuite>
</testsuites>	`
	withNoTestSuites    = `<testsuites name="SomeTest" time="0"></testsuites>`
	withSingleTestSuite = `<testsuites name="SomeTest" time="0"><testsuite tests="0" time="0" name="AllAndroidTests"><properties></properties></testsuite></testsuites>`
)

func TestXmlOutput(t *testing.T) {

	Convey("Given a TestSuites object", t, func() {
		suites := &TestRun{Name: "SomeTest"}
		So(suites, ShouldNotBeNil)
		log.Println(suites.Name)

		out, err := xml.MarshalIndent(suites, "", "     ")
		So(err, ShouldBeNil)
		log.Println(string(out))
		So(string(out), ShouldEqual, withNoTestSuites)

		Convey("When a TestSuite is added", func() {
			suite := TestSuite{Name: db.NewNullString("AllAndroidTests"), Tests: 0}
			suites.TestSuites = append(suites.TestSuites, suite)

			out, err := xml.Marshal(&suites)
			So(err, ShouldBeNil)

			Convey("Then the xml should contain the correct Test Suite information", func() {
				So(string(out), ShouldEqual, withSingleTestSuite)
			})
		})
	})

	Convey("Given a Impoted XML object", t, func() {
		suites := TestRun{}
		err := xml.Unmarshal([]byte(testcase), &suites)

		Convey("When there are no errors", func() {
			So(err, ShouldBeNil)
			Convey("Then suites should contain the appropriate info", func() {
				So(suites, ShouldNotBeNil)
				So(suites.Time.Float64, ShouldEqual, 0.069)

				out, _ := json.Marshal(&suites)
				log.Println(string(out))
				out, _ = xml.MarshalIndent(&suites, "", "     ")
				log.Println(string(out))
			})

			Convey("And should contain a test suite", func() {
				So(suites.TestSuites, ShouldNotBeNil)
				So(1, ShouldEqual, len(suites.TestSuites))
			})
		})
	})
}
