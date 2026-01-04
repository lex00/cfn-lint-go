package errors

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE0010_Metadata(t *testing.T) {
	rule := &E0010{}

	if rule.ID() != "E0010" {
		t.Errorf("Expected ID E0010, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	if rule.Description() == "" {
		t.Error("Description should not be empty")
	}

	if rule.Source() == "" {
		t.Error("Source should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}

func TestE0010_ValidTemplate(t *testing.T) {
	// E0010 is triggered by SAM transform failures which happen
	// at transform time, not during rule matching.
	// This test just verifies the rule doesn't match on valid templates.
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: index.handler
      Runtime: nodejs18.x
      CodeUri: ./src
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E0010{}
	matches := rule.Match(tmpl)

	// This rule is used for documentation; actual SAM transform errors
	// are caught at transform time by the linter
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid template, got %d", len(matches))
	}
}
