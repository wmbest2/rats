package android

import (
	"bytes"
	"fmt"
	"github.com/wmbest2/android/adb"
	"github.com/wmbest2/android/apk"
	"github.com/wmbest2/rats/core"
	"github.com/wmbest2/rats/test"
	"log"
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
	LONG_MSG
)

type instToken struct {
	Timestamp time.Time
	Type      int
	Value     []byte
}

type RunPair struct {
	Tests  *test.TestSuite
	Device *core.Device
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
		case "longMsg":
			token.Type = LONG_MSG
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

func processLastToken(test *test.TestCase, token *instToken) bool {
	if token == nil {
		return false
	}
	switch token.Type {
	case TEST:
		test.Name = strings.TrimSpace(string(token.Value))
	case CLASS:
		test.Classname = strings.TrimSpace(string(token.Value))
	case LONG_MSG:
		fallthrough
	case STACK:
		test.Stack = string(token.Value) + "\n"
	default:
		return false
	}
	return true
}

func parseInstrumentation(suite *test.TestSuite, in chan []byte) {
	instrumentCheck := regexp.MustCompile("INSTRUMENTATION_(?:STATUS|RESULT)(?:(?:_(CODE): (.*))|(?:: ([^=\n]*)=(.*)))")
	var currentTest *test.TestCase
	var lastToken *instToken
	var startTime, endTime time.Time
	for v := range in {
		if instrumentCheck.Match(v) {

			if currentTest == nil {
				currentTest = &test.TestCase{}
			}

			vals := instrumentCheck.FindSubmatch(v)
			lastToken = tokenForLine(vals)

			if suite.Tests == 0 && lastToken.Type == NUM_TESTS {
				suite.Tests, _ = strconv.Atoi(string(lastToken.Value))
			}

			processLastToken(currentTest, lastToken)
			if lastToken.Type == CODE && string(lastToken.Value) == "1" {
				startTime = lastToken.Timestamp
			} else if lastToken.Type == CODE || lastToken.Type == LONG_MSG {
				endTime = lastToken.Timestamp
				v := string(lastToken.Value)
				switch {
				case lastToken.Type == LONG_MSG:
					currentTest.Stack = fmt.Sprintf("Test failed to run to completion. Reason: '%s'. Check device logcat for details", currentTest.Stack)
					fallthrough
				case v == "-3":
					currentTest.Skipped = true
					suite.Skipped++
				case v == "-2":
					currentTest.Failures = append(currentTest.Failures, currentTest.Stack)
					suite.Failures++
				case v == "-1":
					currentTest.Errors = append(currentTest.Errors, currentTest.Stack)
					suite.Errors++
				}

				currentTest.Time = endTime.Sub(startTime).Seconds()
				suite.Time += currentTest.Time
				suite.TestCases = append(suite.TestCases, *currentTest)
				currentTest = nil
				lastToken = nil

				if suite.Tests != 0 && suite.Tests == len(suite.TestCases) {
					// return early
					return
				}
			}
		} else {
			if lastToken != nil && (lastToken.Type == STACK || lastToken.Type == LONG_MSG) {
				currentTest.Stack += fmt.Sprintf("%s\n", v)
			}
		}
	}
}

func RunTests(manifest *apk.Manifest, devices []*core.Device, artifacts []string) (chan *core.Device, chan *test.TestRun) {
	finished := make(chan *core.Device)
	out := make(chan *test.TestRun)
	go func(out chan *test.TestRun, finished chan *core.Device, artifacts []string) {
		in := make(chan *RunPair)
		count := 0
		suites := &test.TestRun{Success: test.NewNullBool(true)}

		for _, d := range devices {
			go RunTest(d, manifest, in)
			count++
		}

		log.Printf("\t -> Waiting for %d Tests\n", count)

		for {
			select {
			case run := <-in:
				log.Printf("\t -> Sent result\n")

				for _, artifact := range artifacts {
					buffer := bytes.NewBuffer(nil)
					file := fmt.Sprintf("%s", artifact)

					log.Printf("\t -> Fetching file %s\n", file)
					adb.ShellSync(run.Device, fmt.Sprintf("run-as %s chmod 666 %s", manifest.Instrument.Target, file))
					err := adb.Pull(run.Device, buffer, file)

					if err == nil {
						result := test.Artifact{Name: test.NewNullString(artifact), Data: buffer.Bytes()}
						run.Tests.Artifacts = append(run.Tests.Artifacts, result)
					} else {
						log.Printf("\t -> %s\n", err)
					}
				}

				finished <- run.Device
				suites.TestSuites = append(suites.TestSuites, *run.Tests)
				suites.Time.Float64 += run.Tests.Time
				suites.Success.Bool = suites.Success.Bool && run.Tests.Failures == 0 && run.Tests.Errors == 0
				count--
			}

			if count == 0 {
				break
			}
		}
		out <- suites
		close(finished)
		close(out)
	}(out, finished, artifacts)

	return finished, out
}

func LogTestSuite(device *core.Device, manifest *apk.Manifest, out chan *RunPair) {
	testRunner := fmt.Sprintf("%s/%s", manifest.Package, manifest.Instrument.Name)
	in := adb.Shell(device, "am", "instrument", "-r", "true", "-w", "-e", "log", testRunner)
	suite := test.TestSuite{Hostname: test.NewNullString(device.Serial), Name: test.NewNullString(device.String())}
	parseInstrumentation(&suite, in)
	out <- &RunPair{Tests: &suite, Device: device}
}

func RunTest(device *core.Device, manifest *apk.Manifest, out chan *RunPair) {
	testRunner := fmt.Sprintf("%s/%s", manifest.Package, manifest.Instrument.Name)

	in := adb.Shell(device, "am", "instrument", "-r", "-w", "-e", "coverage", "true", testRunner)
	suite := test.TestSuite{Hostname: test.NewNullString(device.Serial), Name: test.NewNullString(device.String())}
	parseInstrumentation(&suite, in)

	log.Printf("\t -> Finished test on %s\n", device.Serial)

	out <- &RunPair{Tests: &suite, Device: device}
}
