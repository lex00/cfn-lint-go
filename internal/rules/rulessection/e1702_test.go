package rulessection

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE1702_ValidRuleCondition(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  Environment:
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
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E1702{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid RuleCondition, got %d: %v", len(matches), matches)
	}
}

func TestE1702_NoRuleCondition(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  Environment:
    Type: String
Rules:
  TestRule:
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

	rule := &E1702{}
	matches := rule.Match(tmpl)

	// RuleCondition is optional, so this should be valid
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when RuleCondition is absent, got %d", len(matches))
	}
}

func TestE1702_ValidRuleFunctions(t *testing.T) {
	functions := []string{
		"Fn::And",
		"Fn::Contains",
		"Fn::EachMemberEquals",
		"Fn::EachMemberIn",
		"Fn::Equals",
		"Fn::If",
		"Fn::Not",
		"Fn::Or",
		"Fn::RefAll",
		"Fn::ValueOf",
		"Fn::ValueOfAll",
	}

	for _, fn := range functions {
		yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  Environment:
    Type: String
Rules:
  TestRule:
    RuleCondition:
      ` + fn + `:
        - test
        - test
    Assertions:
      - Assert:
          Fn::Equals:
            - a
            - b
`
		tmpl, err := template.Parse([]byte(yaml))
		if err != nil {
			t.Fatalf("Failed to parse for %s: %v", fn, err)
		}

		rule := &E1702{}
		matches := rule.Match(tmpl)

		if len(matches) != 0 {
			t.Errorf("Expected 0 matches for valid function %s, got %d", fn, len(matches))
		}
	}
}

func TestE1702_InvalidFunction(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Rules:
  InvalidRule:
    RuleCondition:
      Fn::GetAtt:
        - MyResource
        - Arn
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

	rule := &E1702{}
	matches := rule.Match(tmpl)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for invalid function, got %d", len(matches))
		return
	}

	if matches[0].Path[2] != "RuleCondition" {
		t.Errorf("Expected error on RuleCondition, got path: %v", matches[0].Path)
	}
}

func TestE1702_RuleConditionWithFnEquals(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  Environment:
    Type: String
  UseSSL:
    Type: String
Rules:
  SSLRule:
    RuleCondition:
      Fn::Equals:
        - !Ref UseSSL
        - 'Yes'
    Assertions:
      - Assert:
          Fn::Not:
            - Fn::Equals:
                - !Ref SSLCertificate
                - ''
        AssertDescription: SSL Certificate is required when UseSSL is Yes
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E1702{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid RuleCondition with Fn::Equals, got %d", len(matches))
	}
}

func TestE1702_Metadata(t *testing.T) {
	rule := &E1702{}

	if rule.ID() != "E1702" {
		t.Errorf("Expected ID E1702, got %s", rule.ID())
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
