package informational

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestI3037_NoDuplicates(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyInstance:
    Type: AWS::EC2::Instance
    Properties:
      SecurityGroupIds:
        - sg-12345
        - sg-67890
      ImageId: ami-12345
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I3037{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for list without duplicates, got %d", len(matches))
	}
}

func TestI3037_WithDuplicates(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyInstance:
    Type: AWS::EC2::Instance
    Properties:
      SecurityGroupIds:
        - sg-12345
        - sg-67890
        - sg-12345
      ImageId: ami-12345
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I3037{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for list with duplicate values")
	}
}

func TestI3037_NonCheckedResource(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I3037{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for unchecked resource type, got %d", len(matches))
	}
}

func TestI3037_Metadata(t *testing.T) {
	rule := &I3037{}

	if rule.ID() != "I3037" {
		t.Errorf("Expected ID I3037, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
