package test

import (
	"fmt"
	"github.com/wmbest2/android/adb"
	"github.com/wmbest2/android/apk"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	NUM_TESTS = iota
	STREAM
	ID
	TEST
	CLASS
	CURRENT
	STACK
	CODE
)

type instToken struct {
	Timestamp time.Time
	Type      int
	Value     []byte
}

func tokenForLine(line [][]byte) *instToken {
	token := &instToken{Timestamp: time.Now()}
	if string(line[1]) == "CODE" {
		token.Type = CODE
		token.Value = line[2]
	} else {
		switch string(line[3]) {
		case "numTests":
			token.Type = NUM_TESTS
		case "stream":
			token.Type = STREAM
		case "stack":
			token.Type = STACK
		case "id":
			token.Type = ID
		case "test":
			token.Type = TEST
		case "current":
			token.Type = CURRENT
		case "class":
			token.Type = CLASS
		}

		token.Value = line[4]
	}
	return token
}

func processLastToken(test *TestCase, token *instToken) bool {
	if token == nil {
		return false
	}
	switch token.Type {
	case TEST:
		test.Name = strings.TrimSpace(string(token.Value))
	case CLASS:
		test.Classname = strings.TrimSpace(string(token.Value))
	case STACK:
		test.Stack = string(token.Value) + "\n"
	default:
		return false
	}
	return true
}

func parseInstrumentation(suite *TestSuite, in chan interface{}, out chan *TestSuite) {
	instrumentCheck := regexp.MustCompile("INSTRUMENTATION_STATUS(?:(?:_(CODE): (.*))|(?:: ([^=\n]*)=(.*)))")
	var currentTest *TestCase
	var lastToken *instToken
	var startTime, endTime time.Time
	var v interface{}

	ok := true

	for {
		if !ok {
			break
		}

		switch v, ok = <-in; v.(type) {
		case []byte:
			if instrumentCheck.Match(v.([]byte)) {

				if currentTest == nil {
					currentTest = &TestCase{}
					suite.Tests++
				}

				vals := instrumentCheck.FindSubmatch(v.([]byte))
				lastToken = tokenForLine(vals)

				if suite.Tests == 0 && lastToken.Type == NUM_TESTS {
					suite.Tests, _ = strconv.Atoi(string(lastToken.Value))
				}

				processLastToken(currentTest, lastToken)
				if lastToken.Type == CODE && string(lastToken.Value) == "1" {
					startTime = lastToken.Timestamp
				} else if lastToken.Type == CODE {
					endTime = lastToken.Timestamp
					switch string(lastToken.Value) {
					case "-2":
						currentTest.Failure = &Failure{Contents: currentTest.Stack}
						suite.Failures++
					case "-1":
						currentTest.Error = &Error{Contents: currentTest.Stack}
						suite.Errors++
					}

					currentTest.Time = endTime.Sub(startTime).Seconds()
					suite.Time += currentTest.Time
					suite.TestCases = append(suite.TestCases, currentTest)
					currentTest = nil
					lastToken = nil
				}
			} else {
				if lastToken != nil && lastToken.Type == STACK {
					currentTest.Stack += string(v.([]byte))
				}
			}
		case error:
			fmt.Println(v.(error))
		}
	}
	out <- suite
}

func LogTestSuite(device *adb.Device, manifest *apk.Manifest, out chan *TestSuite) {
	testRunner := fmt.Sprintf("%s/%s", manifest.Package, manifest.Instrument.Name)
	in := device.Exec("shell", "am", "instrument", "-r", "-e", "log", "true","-w", testRunner)
	suite := TestSuite{Hostname: device.Serial, Name: device.String()}
	parseInstrumentation(&suite, in, out)
}

func RunTest(device *adb.Device, manifest *apk.Manifest, out chan *TestSuite) {
	testRunner := fmt.Sprintf("%s/%s", manifest.Package, manifest.Instrument.Name)
	in := device.Exec("shell", "am", "instrument", "-r", "-w", testRunner)
	suite := TestSuite{Hostname: device.Serial, Name: device.String()}
	parseInstrumentation(&suite, in, out)
}
