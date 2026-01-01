package conditions

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE8002_ValidConditions(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  Environment:
    Type: String

Conditions:
  IsProd: !Equals [!Ref Environment, prod]

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Condition: IsProd
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E8002{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid conditions, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE8002_UndefinedResourceCondition(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Condition: NonExistent
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E8002{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for undefined condition on resource")
	}
}

func TestE8002_UndefinedOutputCondition(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket

Outputs:
  BucketName:
    Value: !Ref MyBucket
    Condition: NonExistent
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E8002{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for undefined condition on output")
	}
}

func TestE8002_UndefinedFnIfCondition(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Conditions:
  IsProd: !Equals ["a", "a"]

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !If [NonExistent, prod-bucket, dev-bucket]
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E8002{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for undefined condition in Fn::If")
	}
}

func TestE8002_NestedConditionReference(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Conditions:
  IsProd: !Equals ["a", "a"]
  IsEnabled: !And
    - !Condition IsProd
    - !Condition NonExistent
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E8002{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for undefined condition in nested Condition")
	}
}

func TestE8002_NoConditions(t *testing.T) {
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

	rule := &E8002{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when no conditions, got %d", len(matches))
	}
}

func TestE8002_Metadata(t *testing.T) {
	rule := &E8002{}

	if rule.ID() != "E8002" {
		t.Errorf("Expected ID E8002, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
