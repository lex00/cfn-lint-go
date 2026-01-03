package informational

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestI3013_LogGroupWithoutRetention(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyLogGroup:
    Type: AWS::Logs::LogGroup
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I3013{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for log group without retention period")
	}
}

func TestI3013_LogGroupWithRetention(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyLogGroup:
    Type: AWS::Logs::LogGroup
    Properties:
      RetentionInDays: 7
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I3013{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for log group with retention period, got %d", len(matches))
	}
}

func TestI3013_NonRetentionResource(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Runtime: python3.11
      Handler: index.handler
      Code:
        ZipFile: |
          def handler(event, context):
            return 'Hello'
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I3013{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for resource without retention concept, got %d", len(matches))
	}
}

func TestI3013_Metadata(t *testing.T) {
	rule := &I3013{}

	if rule.ID() != "I3013" {
		t.Errorf("Expected ID I3013, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
