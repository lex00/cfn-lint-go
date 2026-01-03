package informational

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestI3100_LegacyInstanceType(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyInstance:
    Type: AWS::EC2::Instance
    Properties:
      InstanceType: t2.micro
      ImageId: ami-12345
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I3100{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for legacy t2 instance type")
	}
}

func TestI3100_CurrentGenInstanceType(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyInstance:
    Type: AWS::EC2::Instance
    Properties:
      InstanceType: t3.micro
      ImageId: ami-12345
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I3100{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for current gen t3 instance type, got %d", len(matches))
	}
}

func TestI3100_LaunchTemplate(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyLaunchTemplate:
    Type: AWS::EC2::LaunchTemplate
    Properties:
      LaunchTemplateData:
        InstanceType: m4.large
        ImageId: ami-12345
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I3100{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for legacy m4 instance type in launch template")
	}
}

func TestI3100_Metadata(t *testing.T) {
	rule := &I3100{}

	if rule.ID() != "I3100" {
		t.Errorf("Expected ID I3100, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
