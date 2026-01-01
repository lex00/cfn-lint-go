package mappings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE7001_ValidMapping(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Mappings:
  RegionMap:
    us-east-1:
      AMI: ami-12345678
    us-west-2:
      AMI: ami-87654321

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E7001{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid mapping, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE7001_InvalidMappingName(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Mappings:
  Region@Map:
    us-east-1:
      AMI: ami-12345678

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E7001{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for invalid mapping name with special character")
	}
}

func TestE7001_InvalidTopLevelKey(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Mappings:
  RegionMap:
    us-east-1:
      AMI: ami-12345678
    "invalid key":
      AMI: ami-87654321

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E7001{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for invalid top-level key with space")
	}
}

func TestE7001_EmptyMapping(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Mappings:
  EmptyMap: {}

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E7001{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for empty mapping")
	}
}

func TestE7001_EmptySecondLevel(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Mappings:
  RegionMap:
    us-east-1: {}

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E7001{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for empty second-level mapping")
	}
}

func TestE7001_NoMappings(t *testing.T) {
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

	rule := &E7001{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when no mappings, got %d", len(matches))
	}
}

func TestE7001_Metadata(t *testing.T) {
	rule := &E7001{}

	if rule.ID() != "E7001" {
		t.Errorf("Expected ID E7001, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
