package mappings

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE7010_UnderLimit(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Mappings:
  Map1:
    us-east-1:
      AMI: ami-12345
  Map2:
    us-east-1:
      AMI: ami-67890

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E7010{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for under limit, got %d", len(matches))
	}
}

func TestE7010_OverLimit(t *testing.T) {
	// Generate template with 201 mappings
	var mappings []string
	for i := 0; i <= MaxMappings; i++ {
		mappings = append(mappings, fmt.Sprintf("  Map%d:\n    us-east-1:\n      Value: val%d", i, i))
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

	rule := &E7010{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for exceeding mapping limit")
	}
}

func TestE7010_ExactlyAtLimit(t *testing.T) {
	// Generate template with exactly 200 mappings
	var mappings []string
	for i := 0; i < MaxMappings; i++ {
		mappings = append(mappings, fmt.Sprintf("  Map%d:\n    us-east-1:\n      Value: val%d", i, i))
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

	rule := &E7010{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches at exactly limit, got %d", len(matches))
	}
}

func TestE7010_NoMappings(t *testing.T) {
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

	rule := &E7010{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when no mappings, got %d", len(matches))
	}
}

func TestE7010_Metadata(t *testing.T) {
	rule := &E7010{}

	if rule.ID() != "E7010" {
		t.Errorf("Expected ID E7010, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
