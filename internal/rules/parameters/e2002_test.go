package parameters

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE2002_ValidTypes(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  StringParam:
    Type: String
  NumberParam:
    Type: Number
  ListParam:
    Type: List<Number>
  CommaParam:
    Type: CommaDelimitedList
  VpcParam:
    Type: AWS::EC2::VPC::Id
  SubnetListParam:
    Type: List<AWS::EC2::Subnet::Id>
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2002{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid types, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE2002_InvalidType(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  BadParam:
    Type: InvalidType
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2002{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for invalid parameter type")
	}
}

func TestE2002_MissingType(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  NoTypeParam:
    Default: value
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2002{}
	matches := rule.Match(tmpl)

	// E2001 handles missing Type, E2002 should skip
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for missing type (handled by E2001), got %d", len(matches))
	}
}

func TestE2002_Metadata(t *testing.T) {
	rule := &E2002{}

	if rule.ID() != "E2002" {
		t.Errorf("Expected ID E2002, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
