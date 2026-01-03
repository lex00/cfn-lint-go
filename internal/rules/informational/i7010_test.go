package informational

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestI7010_LowMappingCount(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Mappings:
  RegionMap:
    us-east-1:
      AMI: ami-12345
  InstanceTypeMap:
    us-east-1:
      Type: t3.micro
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I7010{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for low mapping count, got %d", len(matches))
	}
}

func TestI7010_ApproachingLimit(t *testing.T) {
	// Generate template with 165 mappings (82.5% of limit)
	var mappings []string
	for i := 0; i < 165; i++ {
		mappings = append(mappings, fmt.Sprintf("  Mapping%d:\n    us-east-1:\n      AMI: ami-12345", i))
	}

	yaml := fmt.Sprintf(`
AWSTemplateFormatVersion: '2010-09-09'
Mappings:
%s
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`, strings.Join(mappings, "\n"))

	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I7010{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for mapping count approaching limit")
	}
}

func TestI7010_Metadata(t *testing.T) {
	rule := &I7010{}

	if rule.ID() != "I7010" {
		t.Errorf("Expected ID I7010, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
