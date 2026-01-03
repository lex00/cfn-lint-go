package rulessection

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE1701_ValidAssertion(t *testing.T) {
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

	rule := &E1701{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid assertion, got %d: %v", len(matches), matches)
	}
}

func TestE1701_AssertionWithoutDescription(t *testing.T) {
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
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E1701{}
	matches := rule.Match(tmpl)

	// AssertDescription is optional, so this should be valid
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for assertion without description, got %d", len(matches))
	}
}

func TestE1701_MissingAssert(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Rules:
  InvalidRule:
    Assertions:
      - AssertDescription: Missing Assert property
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E1701{}
	matches := rule.Match(tmpl)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for missing Assert, got %d", len(matches))
		return
	}

	if matches[0].Message != "Rule 'InvalidRule' Assertion[0] is missing required property 'Assert'" {
		t.Errorf("Unexpected error message: %s", matches[0].Message)
	}
}

func TestE1701_InvalidProperty(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Rules:
  InvalidRule:
    Assertions:
      - Assert:
          Fn::Equals:
            - a
            - b
        InvalidProperty: value
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E1701{}
	matches := rule.Match(tmpl)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for invalid property, got %d", len(matches))
		return
	}

	if matches[0].Path[4] != "InvalidProperty" {
		t.Errorf("Expected error on InvalidProperty, got path: %v", matches[0].Path)
	}
}

func TestE1701_MultipleAssertions(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  Environment:
    Type: String
  InstanceType:
    Type: String
Rules:
  MultiAssertRule:
    Assertions:
      - Assert:
          Fn::Contains:
            - [dev, prod]
            - !Ref Environment
        AssertDescription: Environment must be dev or prod
      - Assert:
          Fn::Contains:
            - [t3.micro, t3.small]
            - !Ref InstanceType
        AssertDescription: InstanceType must be t3.micro or t3.small
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E1701{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid multiple assertions, got %d", len(matches))
	}
}

func TestE1701_Metadata(t *testing.T) {
	rule := &E1701{}

	if rule.ID() != "E1701" {
		t.Errorf("Expected ID E1701, got %s", rule.ID())
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
