package outputs

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE6002_ValidOutput(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket

Outputs:
  BucketName:
    Description: The bucket name
    Value: !Ref MyBucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E6002{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid output, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE6002_MissingValue(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket

Outputs:
  BucketName:
    Description: The bucket name
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E6002{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for output missing Value")
	}
}

func TestE6002_MultipleOutputs(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket

Outputs:
  ValidOutput:
    Value: !Ref MyBucket
  MissingValue1:
    Description: No value here
  MissingValue2:
    Description: Also no value
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E6002{}
	matches := rule.Match(tmpl)

	if len(matches) != 2 {
		t.Errorf("Expected 2 matches for 2 outputs missing Value, got %d", len(matches))
	}
}

func TestE6002_NoOutputs(t *testing.T) {
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

	rule := &E6002{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when no outputs, got %d", len(matches))
	}
}

func TestE6002_Metadata(t *testing.T) {
	rule := &E6002{}

	if rule.ID() != "E6002" {
		t.Errorf("Expected ID E6002, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
