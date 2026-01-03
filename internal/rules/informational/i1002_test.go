package informational

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestI1002_SmallTemplate(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Description: Small template
Parameters:
  Param1:
    Type: String
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I1002{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for small template, got %d", len(matches))
	}
}

func TestI1002_LargeTemplate(t *testing.T) {
	// Generate template with many resources to exceed threshold
	var resources []string
	for i := 0; i < 150; i++ {
		resources = append(resources, fmt.Sprintf(`  Bucket%d:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !Sub 'my-bucket-%d-${AWS::AccountId}'`, i, i))
	}

	yaml := fmt.Sprintf(`
AWSTemplateFormatVersion: '2010-09-09'
Description: Large template approaching size limit
Resources:
%s
`, strings.Join(resources, "\n"))

	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I1002{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for large template approaching size limit")
	}
}

func TestI1002_Metadata(t *testing.T) {
	rule := &I1002{}

	if rule.ID() != "I1002" {
		t.Errorf("Expected ID I1002, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
