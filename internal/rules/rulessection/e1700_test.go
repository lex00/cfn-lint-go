package rulessection

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE1700_ValidRule(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  Environment:
    Type: String
Rules:
  ValidRule:
    Assertions:
      - Assert:
          Fn::Contains:
            - [dev, prod]
            - !Ref Environment
        AssertDescription: Environment must be dev or prod
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E1700{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid rule, got %d: %v", len(matches), matches)
	}
}

func TestE1700_RuleWithCondition(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  Environment:
    Type: String
  InstanceType:
    Type: String
Rules:
  TestRule:
    RuleCondition:
      Fn::Equals:
        - !Ref Environment
        - test
    Assertions:
      - Assert:
          Fn::Contains:
            - [t3.micro]
            - !Ref InstanceType
        AssertDescription: Test environment must use t3.micro
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E1700{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid rule with condition, got %d: %v", len(matches), matches)
	}
}

func TestE1700_MissingAssertions(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Rules:
  InvalidRule:
    RuleCondition:
      Fn::Equals:
        - !Ref Environment
        - test
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E1700{}
	matches := rule.Match(tmpl)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for missing Assertions, got %d", len(matches))
		return
	}

	if matches[0].Message != "Rule 'InvalidRule' is missing required property 'Assertions'" {
		t.Errorf("Unexpected error message: %s", matches[0].Message)
	}
}

func TestE1700_InvalidProperty(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Rules:
  InvalidRule:
    InvalidProperty: value
    Assertions:
      - Assert:
          Fn::Equals:
            - a
            - b
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E1700{}
	matches := rule.Match(tmpl)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for invalid property, got %d", len(matches))
		return
	}

	if matches[0].Path[2] != "InvalidProperty" {
		t.Errorf("Expected error on InvalidProperty, got path: %v", matches[0].Path)
	}
}

func TestE1700_Metadata(t *testing.T) {
	rule := &E1700{}

	if rule.ID() != "E1700" {
		t.Errorf("Expected ID E1700, got %s", rule.ID())
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
