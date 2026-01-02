package resources

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3040_ValidProperties(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: my-bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3040{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid properties, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE3040_ReadOnlyProperty(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: my-bucket
      Arn: arn:aws:s3:::my-bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3040{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for read-only property 'Arn'")
	} else {
		if !strings.Contains(matches[0].Message, "read-only") {
			t.Errorf("Error message should mention read-only: %s", matches[0].Message)
		}
		if !strings.Contains(matches[0].Message, "Arn") {
			t.Errorf("Error message should mention property name: %s", matches[0].Message)
		}
	}
}

func TestE3040_MultipleReadOnlyProperties(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: my-bucket
      Arn: arn:aws:s3:::my-bucket
      DomainName: my-bucket.s3.amazonaws.com
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3040{}
	matches := rule.Match(tmpl)

	if len(matches) < 2 {
		t.Errorf("Expected at least 2 matches for multiple read-only properties, got %d", len(matches))
	}
}

func TestE3040_UnknownResourceSkipped(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyCustom:
    Type: AWS::Unknown::Resource
    Properties:
      Arn: some-arn
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3040{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for unknown resource type, got %d", len(matches))
	}
}

func TestE3040_Metadata(t *testing.T) {
	rule := &E3040{}

	if rule.ID() != "E3040" {
		t.Errorf("Expected ID E3040, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
