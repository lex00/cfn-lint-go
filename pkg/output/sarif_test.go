package output

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/lint"
)

func TestWriteSARIF_EmptyMatches(t *testing.T) {
	var buf bytes.Buffer
	err := WriteSARIF(&buf, []lint.Match{}, "1.0.0")
	if err != nil {
		t.Fatalf("WriteSARIF failed: %v", err)
	}

	var sarif SARIF
	if err := json.Unmarshal(buf.Bytes(), &sarif); err != nil {
		t.Fatalf("Failed to parse SARIF output: %v", err)
	}

	if sarif.Version != "2.1.0" {
		t.Errorf("Expected SARIF version 2.1.0, got %s", sarif.Version)
	}

	if len(sarif.Runs) != 1 {
		t.Errorf("Expected 1 run, got %d", len(sarif.Runs))
	}

	if len(sarif.Runs[0].Results) != 0 {
		t.Errorf("Expected empty results, got %d", len(sarif.Runs[0].Results))
	}
}

func TestWriteSARIF_WithMatches(t *testing.T) {
	matches := []lint.Match{
		{
			Rule: lint.MatchRule{
				ID:               "E1001",
				ShortDescription: "Missing Type",
				Description:      "Resource is missing Type property",
			},
			Location: lint.MatchLocation{
				Filename: "template.yaml",
				Start:    lint.MatchPosition{LineNumber: 10, ColumnNumber: 5},
				End:      lint.MatchPosition{LineNumber: 10, ColumnNumber: 20},
			},
			Level:   "Error",
			Message: "Resource 'MyBucket' is missing Type property",
		},
		{
			Rule: lint.MatchRule{
				ID:               "W2001",
				ShortDescription: "Unused parameter",
				Description:      "Parameter is not used",
			},
			Location: lint.MatchLocation{
				Filename: "template.yaml",
				Start:    lint.MatchPosition{LineNumber: 5, ColumnNumber: 3},
				End:      lint.MatchPosition{LineNumber: 5, ColumnNumber: 15},
			},
			Level:   "Warning",
			Message: "Parameter 'Env' is not used",
		},
		{
			Rule: lint.MatchRule{
				ID:               "I1001",
				ShortDescription: "Suggestion",
				Description:      "Consider using tags",
			},
			Location: lint.MatchLocation{
				Filename: "template.yaml",
				Start:    lint.MatchPosition{LineNumber: 15, ColumnNumber: 1},
				End:      lint.MatchPosition{LineNumber: 15, ColumnNumber: 10},
			},
			Level:   "Informational",
			Message: "Consider adding tags to the bucket",
		},
	}

	var buf bytes.Buffer
	err := WriteSARIF(&buf, matches, "0.15.0")
	if err != nil {
		t.Fatalf("WriteSARIF failed: %v", err)
	}

	var sarif SARIF
	if err := json.Unmarshal(buf.Bytes(), &sarif); err != nil {
		t.Fatalf("Failed to parse SARIF output: %v", err)
	}

	// Check version in driver
	if sarif.Runs[0].Tool.Driver.Version != "0.15.0" {
		t.Errorf("Expected driver version 0.15.0, got %s", sarif.Runs[0].Tool.Driver.Version)
	}

	// Check results count
	if len(sarif.Runs[0].Results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(sarif.Runs[0].Results))
	}

	// Check level mappings
	levelTests := map[string]string{
		"E1001": "error",
		"W2001": "warning",
		"I1001": "note",
	}

	for _, result := range sarif.Runs[0].Results {
		expected, ok := levelTests[result.RuleID]
		if !ok {
			t.Errorf("Unexpected rule ID: %s", result.RuleID)
			continue
		}
		if result.Level != expected {
			t.Errorf("Expected level %s for %s, got %s", expected, result.RuleID, result.Level)
		}
	}
}

func TestWriteSARIF_Schema(t *testing.T) {
	var buf bytes.Buffer
	err := WriteSARIF(&buf, []lint.Match{}, "1.0.0")
	if err != nil {
		t.Fatalf("WriteSARIF failed: %v", err)
	}

	var sarif SARIF
	if err := json.Unmarshal(buf.Bytes(), &sarif); err != nil {
		t.Fatalf("Failed to parse SARIF output: %v", err)
	}

	if sarif.Schema == "" {
		t.Error("Expected schema to be set")
	}

	if sarif.Runs[0].Tool.Driver.Name != "cfn-lint-go" {
		t.Errorf("Expected driver name cfn-lint-go, got %s", sarif.Runs[0].Tool.Driver.Name)
	}
}

func TestWriteSARIF_UniqueRules(t *testing.T) {
	// Two matches with same rule should result in one rule definition
	matches := []lint.Match{
		{
			Rule: lint.MatchRule{
				ID:               "E1001",
				ShortDescription: "Test rule",
			},
			Location: lint.MatchLocation{
				Filename: "template.yaml",
				Start:    lint.MatchPosition{LineNumber: 10, ColumnNumber: 1},
				End:      lint.MatchPosition{LineNumber: 10, ColumnNumber: 10},
			},
			Level:   "Error",
			Message: "First error",
		},
		{
			Rule: lint.MatchRule{
				ID:               "E1001",
				ShortDescription: "Test rule",
			},
			Location: lint.MatchLocation{
				Filename: "template.yaml",
				Start:    lint.MatchPosition{LineNumber: 20, ColumnNumber: 1},
				End:      lint.MatchPosition{LineNumber: 20, ColumnNumber: 10},
			},
			Level:   "Error",
			Message: "Second error",
		},
	}

	var buf bytes.Buffer
	err := WriteSARIF(&buf, matches, "1.0.0")
	if err != nil {
		t.Fatalf("WriteSARIF failed: %v", err)
	}

	var sarif SARIF
	if err := json.Unmarshal(buf.Bytes(), &sarif); err != nil {
		t.Fatalf("Failed to parse SARIF output: %v", err)
	}

	// Should have 1 rule but 2 results
	if len(sarif.Runs[0].Tool.Driver.Rules) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(sarif.Runs[0].Tool.Driver.Rules))
	}
	if len(sarif.Runs[0].Results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(sarif.Runs[0].Results))
	}
}
