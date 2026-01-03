package output

import (
	"bytes"
	"encoding/xml"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/lint"
)

func TestWriteJUnit_EmptyMatches(t *testing.T) {
	var buf bytes.Buffer
	err := WriteJUnit(&buf, []lint.Match{})
	if err != nil {
		t.Fatalf("WriteJUnit failed: %v", err)
	}

	var suites JUnitTestSuites
	if err := xml.Unmarshal(buf.Bytes(), &suites); err != nil {
		t.Fatalf("Failed to parse JUnit output: %v", err)
	}

	if suites.Name != "cfn-lint" {
		t.Errorf("Expected name 'cfn-lint', got %s", suites.Name)
	}

	// Empty matches should create one empty test suite
	if len(suites.TestSuites) != 1 {
		t.Errorf("Expected 1 test suite for empty matches, got %d", len(suites.TestSuites))
	}

	if suites.Tests != 0 {
		t.Errorf("Expected 0 tests, got %d", suites.Tests)
	}

	if suites.Failures != 0 {
		t.Errorf("Expected 0 failures, got %d", suites.Failures)
	}
}

func TestWriteJUnit_WithErrors(t *testing.T) {
	matches := []lint.Match{
		{
			Rule: lint.MatchRule{ID: "E1001"},
			Location: lint.MatchLocation{
				Filename: "template.yaml",
				Start:    lint.MatchPosition{LineNumber: 10, ColumnNumber: 5},
			},
			Level:   "Error",
			Message: "Missing Type property",
		},
		{
			Rule: lint.MatchRule{ID: "E1002"},
			Location: lint.MatchLocation{
				Filename: "template.yaml",
				Start:    lint.MatchPosition{LineNumber: 15, ColumnNumber: 3},
			},
			Level:   "Error",
			Message: "Invalid property",
		},
	}

	var buf bytes.Buffer
	err := WriteJUnit(&buf, matches)
	if err != nil {
		t.Fatalf("WriteJUnit failed: %v", err)
	}

	var suites JUnitTestSuites
	if err := xml.Unmarshal(buf.Bytes(), &suites); err != nil {
		t.Fatalf("Failed to parse JUnit output: %v", err)
	}

	if suites.Tests != 2 {
		t.Errorf("Expected 2 tests, got %d", suites.Tests)
	}

	if suites.Failures != 2 {
		t.Errorf("Expected 2 failures, got %d", suites.Failures)
	}
}

func TestWriteJUnit_WithWarnings(t *testing.T) {
	matches := []lint.Match{
		{
			Rule: lint.MatchRule{ID: "W2001"},
			Location: lint.MatchLocation{
				Filename: "template.yaml",
				Start:    lint.MatchPosition{LineNumber: 5, ColumnNumber: 1},
			},
			Level:   "Warning",
			Message: "Parameter not used",
		},
	}

	var buf bytes.Buffer
	err := WriteJUnit(&buf, matches)
	if err != nil {
		t.Fatalf("WriteJUnit failed: %v", err)
	}

	var suites JUnitTestSuites
	if err := xml.Unmarshal(buf.Bytes(), &suites); err != nil {
		t.Fatalf("Failed to parse JUnit output: %v", err)
	}

	// Warnings count as failures in JUnit
	if suites.Failures != 1 {
		t.Errorf("Expected 1 failure for warning, got %d", suites.Failures)
	}
}

func TestWriteJUnit_InformationalNotFailure(t *testing.T) {
	matches := []lint.Match{
		{
			Rule: lint.MatchRule{ID: "I1001"},
			Location: lint.MatchLocation{
				Filename: "template.yaml",
				Start:    lint.MatchPosition{LineNumber: 5, ColumnNumber: 1},
			},
			Level:   "Informational",
			Message: "Consider adding tags",
		},
	}

	var buf bytes.Buffer
	err := WriteJUnit(&buf, matches)
	if err != nil {
		t.Fatalf("WriteJUnit failed: %v", err)
	}

	var suites JUnitTestSuites
	if err := xml.Unmarshal(buf.Bytes(), &suites); err != nil {
		t.Fatalf("Failed to parse JUnit output: %v", err)
	}

	// Informational messages should not be counted as failures
	if suites.Tests != 1 {
		t.Errorf("Expected 1 test, got %d", suites.Tests)
	}
	if suites.Failures != 0 {
		t.Errorf("Expected 0 failures for informational, got %d", suites.Failures)
	}
}

func TestWriteJUnit_MultipleFiles(t *testing.T) {
	matches := []lint.Match{
		{
			Rule: lint.MatchRule{ID: "E1001"},
			Location: lint.MatchLocation{
				Filename: "template1.yaml",
				Start:    lint.MatchPosition{LineNumber: 10, ColumnNumber: 1},
			},
			Level:   "Error",
			Message: "Error in template1",
		},
		{
			Rule: lint.MatchRule{ID: "E1002"},
			Location: lint.MatchLocation{
				Filename: "template2.yaml",
				Start:    lint.MatchPosition{LineNumber: 5, ColumnNumber: 1},
			},
			Level:   "Error",
			Message: "Error in template2",
		},
		{
			Rule: lint.MatchRule{ID: "E1003"},
			Location: lint.MatchLocation{
				Filename: "template1.yaml",
				Start:    lint.MatchPosition{LineNumber: 20, ColumnNumber: 1},
			},
			Level:   "Error",
			Message: "Another error in template1",
		},
	}

	var buf bytes.Buffer
	err := WriteJUnit(&buf, matches)
	if err != nil {
		t.Fatalf("WriteJUnit failed: %v", err)
	}

	var suites JUnitTestSuites
	if err := xml.Unmarshal(buf.Bytes(), &suites); err != nil {
		t.Fatalf("Failed to parse JUnit output: %v", err)
	}

	// Should group by file into separate test suites
	if len(suites.TestSuites) != 2 {
		t.Errorf("Expected 2 test suites (one per file), got %d", len(suites.TestSuites))
	}

	if suites.Tests != 3 {
		t.Errorf("Expected 3 total tests, got %d", suites.Tests)
	}
}

func TestWriteJUnit_EmptyFilename(t *testing.T) {
	matches := []lint.Match{
		{
			Rule: lint.MatchRule{ID: "E1001"},
			Location: lint.MatchLocation{
				Filename: "",
				Start:    lint.MatchPosition{LineNumber: 10, ColumnNumber: 1},
			},
			Level:   "Error",
			Message: "Error with no filename",
		},
	}

	var buf bytes.Buffer
	err := WriteJUnit(&buf, matches)
	if err != nil {
		t.Fatalf("WriteJUnit failed: %v", err)
	}

	var suites JUnitTestSuites
	if err := xml.Unmarshal(buf.Bytes(), &suites); err != nil {
		t.Fatalf("Failed to parse JUnit output: %v", err)
	}

	// Empty filename should default to "unknown"
	if len(suites.TestSuites) != 1 {
		t.Fatalf("Expected 1 test suite, got %d", len(suites.TestSuites))
	}

	if suites.TestSuites[0].Name != "unknown" {
		t.Errorf("Expected test suite name 'unknown', got %s", suites.TestSuites[0].Name)
	}
}

func TestWriteJUnit_XMLHeader(t *testing.T) {
	var buf bytes.Buffer
	err := WriteJUnit(&buf, []lint.Match{})
	if err != nil {
		t.Fatalf("WriteJUnit failed: %v", err)
	}

	output := buf.String()
	if len(output) < 5 || output[:5] != "<?xml" {
		t.Error("Expected output to start with XML declaration")
	}
}
