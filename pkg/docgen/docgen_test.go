package docgen

import (
	"bytes"
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

// mockRule implements rules.Rule for testing
type mockRule struct {
	id          string
	shortDesc   string
	description string
	source      string
	tags        []string
}

func (r *mockRule) ID() string          { return r.id }
func (r *mockRule) ShortDesc() string   { return r.shortDesc }
func (r *mockRule) Description() string { return r.description }
func (r *mockRule) Source() string      { return r.source }
func (r *mockRule) Tags() []string      { return r.tags }
func (r *mockRule) Match(_ *template.Template) []rules.Match {
	return nil
}

func TestGenerateRulesMarkdown(t *testing.T) {
	testRules := []rules.Rule{
		&mockRule{
			id:          "E1001",
			shortDesc:   "Ref to undefined resource",
			description: "Validates that Ref function references a defined resource or parameter",
			source:      "https://docs.aws.amazon.com/",
			tags:        []string{"functions", "ref"},
		},
		&mockRule{
			id:          "E3001",
			shortDesc:   "Resource configuration error",
			description: "Validates resource configuration",
			source:      "",
			tags:        []string{"resources"},
		},
		&mockRule{
			id:          "W2001",
			shortDesc:   "Unused parameter",
			description: "Finds parameters that are not used",
			source:      "",
			tags:        []string{"parameters"},
		},
		&mockRule{
			id:          "I1001",
			shortDesc:   "Template size info",
			description: "Informational about template size",
			source:      "",
			tags:        []string{"template"},
		},
	}

	var buf bytes.Buffer
	err := GenerateRulesMarkdown(&buf, testRules)
	if err != nil {
		t.Fatalf("GenerateRulesMarkdown() error = %v", err)
	}

	output := buf.String()

	// Check header
	if !strings.Contains(output, "# Rules Reference") {
		t.Error("expected output to contain header '# Rules Reference'")
	}

	// Check that it contains auto-generated notice
	if !strings.Contains(output, "auto-generated") {
		t.Error("expected output to contain auto-generated notice")
	}

	// Check each rule category header exists
	if !strings.Contains(output, "E1xxx") {
		t.Error("expected output to contain E1xxx category")
	}
	if !strings.Contains(output, "E3xxx") {
		t.Error("expected output to contain E3xxx category")
	}
	if !strings.Contains(output, "W2xxx") {
		t.Error("expected output to contain W2xxx category")
	}
	if !strings.Contains(output, "I1xxx") {
		t.Error("expected output to contain I1xxx category")
	}

	// Check rule IDs appear in tables
	if !strings.Contains(output, "E1001") {
		t.Error("expected output to contain rule E1001")
	}
	if !strings.Contains(output, "E3001") {
		t.Error("expected output to contain rule E3001")
	}
	if !strings.Contains(output, "W2001") {
		t.Error("expected output to contain rule W2001")
	}
	if !strings.Contains(output, "I1001") {
		t.Error("expected output to contain rule I1001")
	}

	// Check descriptions appear
	if !strings.Contains(output, "Ref to undefined resource") {
		t.Error("expected output to contain short description")
	}
}

func TestGenerateRulesMarkdownEmpty(t *testing.T) {
	var buf bytes.Buffer
	err := GenerateRulesMarkdown(&buf, nil)
	if err != nil {
		t.Fatalf("GenerateRulesMarkdown() error = %v", err)
	}

	output := buf.String()

	// Should still have header
	if !strings.Contains(output, "# Rules Reference") {
		t.Error("expected output to contain header '# Rules Reference'")
	}

	// Should indicate no rules
	if !strings.Contains(output, "No rules") && !strings.Contains(output, "0 rules") {
		t.Error("expected output to indicate no rules")
	}
}

func TestRuleCategorization(t *testing.T) {
	tests := []struct {
		ruleID   string
		prefix   string
		category string
	}{
		{"E0001", "E0xxx", "Template Errors"},
		{"E1001", "E1xxx", "Functions"},
		{"E2001", "E2xxx", "Parameters"},
		{"E3001", "E3xxx", "Resources"},
		{"E4001", "E4xxx", "Metadata"},
		{"E5001", "E5xxx", "Modules"},
		{"E6001", "E6xxx", "Outputs"},
		{"E7001", "E7xxx", "Mappings"},
		{"E8001", "E8xxx", "Conditions"},
		{"W1001", "W1xxx", "Template Warnings"},
		{"W2001", "W2xxx", "Parameter Warnings"},
		{"W3001", "W3xxx", "Resource Warnings"},
		{"I1001", "I1xxx", "Template Informational"},
		{"I2001", "I2xxx", "Parameter Informational"},
		{"I3001", "I3xxx", "Resource Informational"},
	}

	for _, tt := range tests {
		t.Run(tt.ruleID, func(t *testing.T) {
			prefix := GetRulePrefix(tt.ruleID)
			if prefix != tt.prefix {
				t.Errorf("GetRulePrefix(%s) = %s, want %s", tt.ruleID, prefix, tt.prefix)
			}

			category := GetRuleCategory(tt.ruleID)
			if category != tt.category {
				t.Errorf("GetRuleCategory(%s) = %s, want %s", tt.ruleID, category, tt.category)
			}
		})
	}
}
