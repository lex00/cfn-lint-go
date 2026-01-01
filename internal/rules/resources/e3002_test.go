package resources

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3002_ValidProperties(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: my-bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3002{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid Properties, got %d", len(matches))
	}
}

func TestE3002_PropertiesAsList(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      - Item1
      - Item2
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3002{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for Properties as list")
	}
}

func TestE3002_PropertiesAsScalar(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties: scalar_value
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3002{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for Properties as scalar")
	}
}

func TestE3002_NoProperties(t *testing.T) {
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

	rule := &E3002{}
	matches := rule.Match(tmpl)

	// No Properties is valid (some resources don't require any)
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for no Properties, got %d", len(matches))
	}
}

func TestE3002_Metadata(t *testing.T) {
	rule := &E3002{}

	if rule.ID() != "E3002" {
		t.Errorf("Expected ID E3002, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
