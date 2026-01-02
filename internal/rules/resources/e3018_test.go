package resources

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3018_OneOfSatisfied(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyStack:
    Type: AWS::CloudFormation::Stack
    Properties:
      TemplateURL: https://s3.amazonaws.com/mybucket/mytemplate.yaml
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3018{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when oneOf satisfied, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE3018_OneOfMissing(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyStack:
    Type: AWS::CloudFormation::Stack
    Properties:
      Parameters:
        Key: Value
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3018{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match when oneOf missing")
	} else {
		if !strings.Contains(matches[0].Message, "exactly one of") {
			t.Errorf("Error message should mention 'exactly one of': %s", matches[0].Message)
		}
	}
}

func TestE3018_MultipleOneOf(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyStack:
    Type: AWS::CloudFormation::Stack
    Properties:
      TemplateBody: |
        AWSTemplateFormatVersion: '2010-09-09'
        Resources: {}
      TemplateURL: https://s3.amazonaws.com/mybucket/mytemplate.yaml
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3018{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match when multiple oneOf properties present")
	} else {
		if !strings.Contains(matches[0].Message, "only one is allowed") {
			t.Errorf("Error message should mention 'only one is allowed': %s", matches[0].Message)
		}
	}
}

func TestE3018_CloudWatchAlarmMetrics(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyAlarm:
    Type: AWS::CloudWatch::Alarm
    Properties:
      AlarmName: my-alarm
      MetricName: CPUUtilization
      Namespace: AWS/EC2
      Statistic: Average
      Period: 300
      EvaluationPeriods: 2
      Threshold: 80
      ComparisonOperator: GreaterThanThreshold
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3018{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when MetricName specified, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE3018_UnknownResourceSkipped(t *testing.T) {
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

	rule := &E3018{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for unknown resource type, got %d", len(matches))
	}
}

func TestE3018_Metadata(t *testing.T) {
	rule := &E3018{}

	if rule.ID() != "E3018" {
		t.Errorf("Expected ID E3018, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
