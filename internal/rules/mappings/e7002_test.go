package mappings

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE7002_ValidMappingNameLength(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Mappings:
  RegionMap:
    us-east-1:
      AMI: ami-12345678
Resources:
  MyResource:
    Type: AWS::S3::Bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E7002{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid mapping name length, got %d: %v", len(matches), matches)
	}
}

func TestE7002_MappingNameTooLong(t *testing.T) {
	longName := strings.Repeat("A", 256)
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Mappings:
  ` + longName + `:
    us-east-1:
      AMI: ami-12345678
Resources:
  MyResource:
    Type: AWS::S3::Bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E7002{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for mapping name too long, got %d", len(matches))
	}
}
