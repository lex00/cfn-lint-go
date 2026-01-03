package informational

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestI3010_LowResourceCount(t *testing.T) {
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

	rule := &I3010{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for low resource count, got %d", len(matches))
	}
}

func TestI3010_ApproachingLimit(t *testing.T) {
	// Generate template with 410 resources (82% of limit)
	var resources []string
	for i := 0; i < 410; i++ {
		resources = append(resources, fmt.Sprintf("  Bucket%d:\n    Type: AWS::S3::Bucket", i))
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

	rule := &I3010{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for resource count approaching limit")
	}
}

func TestI3010_Metadata(t *testing.T) {
	rule := &I3010{}

	if rule.ID() != "I3010" {
		t.Errorf("Expected ID I3010, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
