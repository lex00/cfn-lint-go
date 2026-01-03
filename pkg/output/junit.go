package output

import (
	"encoding/xml"
	"io"
	"time"

	"github.com/lex00/cfn-lint-go/pkg/lint"
)

// JUnitTestSuites represents the root element of JUnit XML output.
type JUnitTestSuites struct {
	XMLName    xml.Name         `xml:"testsuites"`
	Name       string           `xml:"name,attr"`
	Tests      int              `xml:"tests,attr"`
	Failures   int              `xml:"failures,attr"`
	Time       string           `xml:"time,attr"`
	TestSuites []JUnitTestSuite `xml:"testsuite"`
}

// JUnitTestSuite represents a test suite (one per file).
type JUnitTestSuite struct {
	Name      string          `xml:"name,attr"`
	Tests     int             `xml:"tests,attr"`
	Failures  int             `xml:"failures,attr"`
	Time      string          `xml:"time,attr"`
	TestCases []JUnitTestCase `xml:"testcase"`
}

// JUnitTestCase represents a single test case (one per rule violation).
type JUnitTestCase struct {
	Name      string          `xml:"name,attr"`
	Classname string          `xml:"classname,attr"`
	Time      string          `xml:"time,attr"`
	Failure   *JUnitFailure   `xml:"failure,omitempty"`
}

// JUnitFailure represents a test failure.
type JUnitFailure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Content string `xml:",chardata"`
}

// WriteJUnit writes matches in JUnit XML format.
func WriteJUnit(w io.Writer, matches []lint.Match) error {
	// Group matches by file
	fileMatches := make(map[string][]lint.Match)
	for _, m := range matches {
		filename := m.Location.Filename
		if filename == "" {
			filename = "unknown"
		}
		fileMatches[filename] = append(fileMatches[filename], m)
	}

	var testSuites []JUnitTestSuite
	totalTests := 0
	totalFailures := 0

	for filename, matches := range fileMatches {
		var testCases []JUnitTestCase
		failures := 0

		for _, m := range matches {
			testCase := JUnitTestCase{
				Name:      m.Rule.ID,
				Classname: filename,
				Time:      "0",
			}

			// Add failure for errors and warnings
			if m.Level == "Error" || m.Level == "Warning" {
				testCase.Failure = &JUnitFailure{
					Message: m.Message,
					Type:    m.Level,
					Content: m.Message,
				}
				failures++
			}

			testCases = append(testCases, testCase)
		}

		testSuite := JUnitTestSuite{
			Name:      filename,
			Tests:     len(testCases),
			Failures:  failures,
			Time:      "0",
			TestCases: testCases,
		}

		testSuites = append(testSuites, testSuite)
		totalTests += len(testCases)
		totalFailures += failures
	}

	// If no matches, create an empty test suite
	if len(testSuites) == 0 {
		testSuites = []JUnitTestSuite{
			{
				Name:      "cfn-lint",
				Tests:     0,
				Failures:  0,
				Time:      "0",
				TestCases: []JUnitTestCase{},
			},
		}
	}

	suites := JUnitTestSuites{
		Name:       "cfn-lint",
		Tests:      totalTests,
		Failures:   totalFailures,
		Time:       time.Now().Format("2006-01-02T15:04:05"),
		TestSuites: testSuites,
	}

	output, err := xml.MarshalIndent(suites, "", "  ")
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}

	_, err = w.Write(output)
	return err
}
