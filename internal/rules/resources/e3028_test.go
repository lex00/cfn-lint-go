package resources

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3028_Metadata(t *testing.T) {
	rule := &E3028{}

	if rule.ID() != "E3028" {
		t.Errorf("Expected ID E3028, got %s", rule.ID())
	}
	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}
	if len(rule.Tags()) == 0 {
		t.Error("Tags should not be empty")
	}
}

func TestE3028_ValidMetadata(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyInstance:
    Type: AWS::EC2::Instance
    Metadata:
      AWS::CloudFormation::Init:
        config:
          packages:
            yum:
              httpd: []
    Properties:
      ImageId: ami-12345678
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3028{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid metadata, got %d", len(matches))
	}
}

func TestE3028_InvalidMetadata(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyInstance:
    Type: AWS::EC2::Instance
    Metadata: "invalid"
    Properties:
      ImageId: ami-12345678
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3028{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for invalid metadata (not an object)")
	}
}
