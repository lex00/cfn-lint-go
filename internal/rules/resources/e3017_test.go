package resources

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3017_AnyOfSatisfied(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyRule:
    Type: AWS::Events::Rule
    Properties:
      Name: my-rule
      ScheduleExpression: rate(5 minutes)
      State: ENABLED
      Targets:
        - Id: target1
          Arn: arn:aws:lambda:us-east-1:123456789012:function:my-function
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3017{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when anyOf satisfied, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE3017_AnyOfNotSatisfied(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyRule:
    Type: AWS::Events::Rule
    Properties:
      Name: my-rule
      State: ENABLED
      Targets:
        - Id: target1
          Arn: arn:aws:lambda:us-east-1:123456789012:function:my-function
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3017{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match when anyOf not satisfied")
	} else {
		if !strings.Contains(matches[0].Message, "at least one of") {
			t.Errorf("Error message should mention 'at least one of': %s", matches[0].Message)
		}
	}
}

func TestE3017_BothAnyOfPresent(t *testing.T) {
	// Having both is fine for anyOf (at least one required)
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyRule:
    Type: AWS::Events::Rule
    Properties:
      Name: my-rule
      EventPattern:
        source:
          - aws.ec2
      ScheduleExpression: rate(5 minutes)
      State: ENABLED
      Targets:
        - Id: target1
          Arn: arn:aws:lambda:us-east-1:123456789012:function:my-function
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3017{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when both anyOf present (valid), got %d", len(matches))
	}
}

func TestE3017_UnknownResourceSkipped(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyCustom:
    Type: AWS::Unknown::Resource
    Properties:
      SomeProp: value
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3017{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for unknown resource type, got %d", len(matches))
	}
}

func TestE3017_Metadata(t *testing.T) {
	rule := &E3017{}

	if rule.ID() != "E3017" {
		t.Errorf("Expected ID E3017, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
