package conditions

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE8001_ValidCondition(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Conditions:
  IsProd:
    !Equals [!Ref Environment, prod]

Parameters:
  Environment:
    Type: String

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E8001{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid condition, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE8001_ValidConditionFunctions(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Conditions:
  Cond1:
    !Equals [a, b]
  Cond2:
    !And
      - !Equals [a, a]
      - !Equals [b, b]
  Cond3:
    !Or
      - !Equals [a, a]
      - !Equals [b, b]
  Cond4:
    !Not
      - !Equals [a, b]
  Cond5:
    Condition: Cond1

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E8001{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid condition functions, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE8001_InvalidConditionScalar(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Conditions:
  BadCondition: just-a-string

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E8001{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for condition with scalar value")
	}
}

func TestE8001_InvalidConditionFunction(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Conditions:
  BadCondition:
    Fn::Join:
      - ","
      - [a, b]

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E8001{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for invalid condition function")
	}
}

func TestE8001_NoConditions(t *testing.T) {
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

	rule := &E8001{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when no conditions, got %d", len(matches))
	}
}

func TestE8001_Metadata(t *testing.T) {
	rule := &E8001{}

	if rule.ID() != "E8001" {
		t.Errorf("Expected ID E8001, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
