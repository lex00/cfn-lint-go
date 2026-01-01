package metadata

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE4002_ValidMetadata(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Metadata:
  AWS::CloudFormation::Interface:
    ParameterGroups:
      - Label:
          default: Network Configuration
        Parameters:
          - VpcId

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E4002{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid metadata, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE4002_NullValue(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Metadata:
  CustomKey: null

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E4002{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for null value in metadata")
	}
}

func TestE4002_NestedNullValue(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Metadata:
  Level1:
    Level2:
      Level3: null

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E4002{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for nested null value in metadata")
	}
}

func TestE4002_NoMetadata(t *testing.T) {
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

	rule := &E4002{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when no metadata, got %d", len(matches))
	}
}

func TestE4002_Metadata(t *testing.T) {
	rule := &E4002{}

	if rule.ID() != "E4002" {
		t.Errorf("Expected ID E4002, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
