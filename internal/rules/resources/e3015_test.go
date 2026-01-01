package resources

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3015_ValidCondition(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  Environment:
    Type: String
Conditions:
  IsProd:
    Fn::Equals:
      - !Ref Environment
      - prod
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Condition: IsProd
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3015{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid condition, got %d: %v", len(matches), matches)
	}
}

func TestE3015_UndefinedCondition(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Condition: NonExistentCondition
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3015{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for undefined condition, got %d", len(matches))
	}
}

func TestE3015_NoCondition(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3015{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for no condition, got %d: %v", len(matches), matches)
	}
}

func TestE3015_MultipleResources(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Conditions:
  IsProd:
    Fn::Equals:
      - true
      - true
Resources:
  ValidResource:
    Type: AWS::S3::Bucket
    Condition: IsProd
  InvalidResource1:
    Type: AWS::S3::Bucket
    Condition: NonExistent1
  InvalidResource2:
    Type: AWS::S3::Bucket
    Condition: NonExistent2
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3015{}
	matches := rule.Match(parsed)

	if len(matches) != 2 {
		t.Errorf("Expected 2 matches for multiple undefined conditions, got %d", len(matches))
	}
}
