package informational

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestI6010_LowOutputCount(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
Outputs:
  BucketName:
    Value: !Ref MyBucket
  BucketArn:
    Value: !GetAtt MyBucket.Arn
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I6010{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for low output count, got %d", len(matches))
	}
}

func TestI6010_ApproachingLimit(t *testing.T) {
	// Generate template with 165 outputs (82.5% of limit)
	var outputs []string
	for i := 0; i < 165; i++ {
		outputs = append(outputs, fmt.Sprintf("  Output%d:\n    Value: !Ref MyBucket", i))
	}

	yaml := fmt.Sprintf(`
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
Outputs:
%s
`, strings.Join(outputs, "\n"))

	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I6010{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for output count approaching limit")
	}
}

func TestI6010_Metadata(t *testing.T) {
	rule := &I6010{}

	if rule.ID() != "I6010" {
		t.Errorf("Expected ID I6010, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
