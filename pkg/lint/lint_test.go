package lint

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestNew(t *testing.T) {
	linter := New(Options{})
	if linter == nil {
		t.Fatal("New() returned nil")
	}
}

func TestLint(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: my-bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	linter := New(Options{})
	matches, err := linter.Lint(tmpl, "test.yaml")
	if err != nil {
		t.Fatalf("Lint failed: %v", err)
	}

	// Valid template should have no matches
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid template, got %d", len(matches))
	}
}

func TestLintWithIgnoredRules(t *testing.T) {
	linter := New(Options{
		IgnoreRules: []string{"E1001"},
	})

	if !linter.isIgnored("E1001") {
		t.Error("Expected E1001 to be ignored")
	}

	if linter.isIgnored("E1002") {
		t.Error("Expected E1002 to not be ignored")
	}
}

func TestLevelFromRuleID(t *testing.T) {
	tests := []struct {
		id       string
		expected string
	}{
		{"E1001", "Error"},
		{"E3002", "Error"},
		{"W1001", "Warning"},
		{"W9999", "Warning"},
		{"I1001", "Informational"},
		{"", "Error"},
		{"X1001", "Error"}, // Unknown prefix defaults to error
	}

	for _, tc := range tests {
		t.Run(tc.id, func(t *testing.T) {
			result := levelFromRuleID(tc.id)
			if result != tc.expected {
				t.Errorf("levelFromRuleID(%q) = %q, want %q", tc.id, result, tc.expected)
			}
		})
	}
}

func TestMatch_JSON(t *testing.T) {
	m := Match{
		Rule: MatchRule{
			ID:               "E1001",
			Description:      "Test description",
			ShortDescription: "Test short",
			Source:           "https://example.com",
		},
		Location: MatchLocation{
			Start:    MatchPosition{LineNumber: 10, ColumnNumber: 5},
			End:      MatchPosition{LineNumber: 10, ColumnNumber: 5},
			Path:     []any{"Resources", "MyBucket"},
			Filename: "test.yaml",
		},
		Level:   "Error",
		Message: "Test message",
	}

	// Verify fields are set correctly
	if m.Rule.ID != "E1001" {
		t.Errorf("Expected rule E1001, got %s", m.Rule.ID)
	}
	if m.Location.Start.LineNumber != 10 {
		t.Errorf("Expected line 10, got %d", m.Location.Start.LineNumber)
	}
	if m.Location.Filename != "test.yaml" {
		t.Errorf("Expected filename test.yaml, got %s", m.Location.Filename)
	}
}
