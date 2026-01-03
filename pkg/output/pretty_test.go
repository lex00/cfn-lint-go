package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/lint"
)

func TestWritePretty_EmptyMatches(t *testing.T) {
	var buf bytes.Buffer
	err := WritePretty(&buf, []lint.Match{}, false)
	if err != nil {
		t.Fatalf("WritePretty failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "No issues found") {
		t.Error("Expected 'No issues found' message for empty matches")
	}
}

func TestWritePretty_EmptyMatchesNoColor(t *testing.T) {
	var buf bytes.Buffer
	err := WritePretty(&buf, []lint.Match{}, true)
	if err != nil {
		t.Fatalf("WritePretty failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "No issues found") {
		t.Error("Expected 'No issues found' message for empty matches")
	}
	// Should not contain color codes
	if strings.Contains(output, "\033[") {
		t.Error("Expected no color codes when noColor=true")
	}
}

func TestWritePretty_WithErrors(t *testing.T) {
	matches := []lint.Match{
		{
			Rule: lint.MatchRule{
				ID:               "E1001",
				ShortDescription: "Missing Type",
			},
			Location: lint.MatchLocation{
				Filename: "template.yaml",
				Start:    lint.MatchPosition{LineNumber: 10, ColumnNumber: 5},
			},
			Level:   "Error",
			Message: "Resource is missing Type property",
		},
	}

	var buf bytes.Buffer
	err := WritePretty(&buf, matches, true) // noColor=true for easier testing
	if err != nil {
		t.Fatalf("WritePretty failed: %v", err)
	}

	output := buf.String()

	// Check file header
	if !strings.Contains(output, "template.yaml") {
		t.Error("Expected filename in output")
	}

	// Check error symbol
	if !strings.Contains(output, "\u2716") { // ✖
		t.Error("Expected error symbol in output")
	}

	// Check line number
	if !strings.Contains(output, "Line 10:5") {
		t.Error("Expected line number in output")
	}

	// Check rule ID
	if !strings.Contains(output, "E1001") {
		t.Error("Expected rule ID in output")
	}

	// Check message
	if !strings.Contains(output, "Resource is missing Type property") {
		t.Error("Expected error message in output")
	}

	// Check summary
	if !strings.Contains(output, "1 errors") {
		t.Error("Expected error count in summary")
	}
}

func TestWritePretty_WithWarnings(t *testing.T) {
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
	err := WritePretty(&buf, matches, true)
	if err != nil {
		t.Fatalf("WritePretty failed: %v", err)
	}

	output := buf.String()

	// Check warning symbol
	if !strings.Contains(output, "\u26A0") { // ⚠
		t.Error("Expected warning symbol in output")
	}

	// Check summary
	if !strings.Contains(output, "1 warnings") {
		t.Error("Expected warning count in summary")
	}
}

func TestWritePretty_WithInformational(t *testing.T) {
	matches := []lint.Match{
		{
			Rule: lint.MatchRule{ID: "I1001"},
			Location: lint.MatchLocation{
				Filename: "template.yaml",
				Start:    lint.MatchPosition{LineNumber: 3, ColumnNumber: 1},
			},
			Level:   "Informational",
			Message: "Consider adding tags",
		},
	}

	var buf bytes.Buffer
	err := WritePretty(&buf, matches, true)
	if err != nil {
		t.Fatalf("WritePretty failed: %v", err)
	}

	output := buf.String()

	// Check info symbol
	if !strings.Contains(output, "\u2139") { // ℹ
		t.Error("Expected info symbol in output")
	}

	// Check summary
	if !strings.Contains(output, "1 info") {
		t.Error("Expected info count in summary")
	}
}

func TestWritePretty_MultipleFiles(t *testing.T) {
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
	}

	var buf bytes.Buffer
	err := WritePretty(&buf, matches, true)
	if err != nil {
		t.Fatalf("WritePretty failed: %v", err)
	}

	output := buf.String()

	// Both filenames should appear
	if !strings.Contains(output, "template1.yaml") {
		t.Error("Expected template1.yaml in output")
	}
	if !strings.Contains(output, "template2.yaml") {
		t.Error("Expected template2.yaml in output")
	}

	// Summary should show 2 errors
	if !strings.Contains(output, "2 errors") {
		t.Error("Expected 2 errors in summary")
	}
}

func TestWritePretty_EmptyFilename(t *testing.T) {
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
	err := WritePretty(&buf, matches, true)
	if err != nil {
		t.Fatalf("WritePretty failed: %v", err)
	}

	output := buf.String()

	// Empty filename should default to "unknown"
	if !strings.Contains(output, "unknown") {
		t.Error("Expected 'unknown' for empty filename")
	}
}

func TestWritePretty_ColorOutput(t *testing.T) {
	matches := []lint.Match{
		{
			Rule: lint.MatchRule{ID: "E1001"},
			Location: lint.MatchLocation{
				Filename: "template.yaml",
				Start:    lint.MatchPosition{LineNumber: 10, ColumnNumber: 1},
			},
			Level:   "Error",
			Message: "Test error",
		},
	}

	var buf bytes.Buffer
	err := WritePretty(&buf, matches, false) // noColor=false
	if err != nil {
		t.Fatalf("WritePretty failed: %v", err)
	}

	output := buf.String()

	// Should contain ANSI color codes
	if !strings.Contains(output, "\033[") {
		t.Error("Expected ANSI color codes when noColor=false")
	}
}

func TestWritePretty_MixedLevels(t *testing.T) {
	matches := []lint.Match{
		{
			Rule: lint.MatchRule{ID: "E1001"},
			Location: lint.MatchLocation{
				Filename: "template.yaml",
				Start:    lint.MatchPosition{LineNumber: 10, ColumnNumber: 1},
			},
			Level:   "Error",
			Message: "Error message",
		},
		{
			Rule: lint.MatchRule{ID: "W2001"},
			Location: lint.MatchLocation{
				Filename: "template.yaml",
				Start:    lint.MatchPosition{LineNumber: 5, ColumnNumber: 1},
			},
			Level:   "Warning",
			Message: "Warning message",
		},
		{
			Rule: lint.MatchRule{ID: "I3001"},
			Location: lint.MatchLocation{
				Filename: "template.yaml",
				Start:    lint.MatchPosition{LineNumber: 3, ColumnNumber: 1},
			},
			Level:   "Informational",
			Message: "Info message",
		},
	}

	var buf bytes.Buffer
	err := WritePretty(&buf, matches, true)
	if err != nil {
		t.Fatalf("WritePretty failed: %v", err)
	}

	output := buf.String()

	// Check all counts in summary
	if !strings.Contains(output, "1 errors") {
		t.Error("Expected 1 error in summary")
	}
	if !strings.Contains(output, "1 warnings") {
		t.Error("Expected 1 warning in summary")
	}
	if !strings.Contains(output, "1 info") {
		t.Error("Expected 1 info in summary")
	}
}

func TestReadFileLines(t *testing.T) {
	// Test with non-existent file
	lines := readFileLines("/non/existent/file.txt")
	if lines != nil {
		t.Error("Expected nil for non-existent file")
	}
}
