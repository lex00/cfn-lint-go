package resources

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3010_UnderLimit(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  Bucket1:
    Type: AWS::S3::Bucket
  Bucket2:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3010{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for under limit, got %d", len(matches))
	}
}

func TestE3010_OverLimit(t *testing.T) {
	// Generate template with 501 resources
	var resources []string
	for i := 0; i <= MaxResources; i++ {
		resources = append(resources, fmt.Sprintf("  Resource%d:\n    Type: AWS::S3::Bucket", i))
	}

	yaml := fmt.Sprintf(`
AWSTemplateFormatVersion: '2010-09-09'
Resources:
%s
`, strings.Join(resources, "\n"))

	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3010{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for exceeding resource limit")
	}
}

func TestE3010_ExactlyAtLimit(t *testing.T) {
	// Generate template with exactly 500 resources
	var resources []string
	for i := 0; i < MaxResources; i++ {
		resources = append(resources, fmt.Sprintf("  Resource%d:\n    Type: AWS::S3::Bucket", i))
	}

	yaml := fmt.Sprintf(`
AWSTemplateFormatVersion: '2010-09-09'
Resources:
%s
`, strings.Join(resources, "\n"))

	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3010{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches at exactly limit, got %d", len(matches))
	}
}

func TestE3010_Metadata(t *testing.T) {
	rule := &E3010{}

	if rule.ID() != "E3010" {
		t.Errorf("Expected ID E3010, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
