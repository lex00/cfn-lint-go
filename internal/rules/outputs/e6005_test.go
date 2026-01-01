package outputs

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE6005_ValidCondition(t *testing.T) {
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
Outputs:
  BucketName:
    Condition: IsProd
    Value: !Ref MyBucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E6005{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid condition, got %d: %v", len(matches), matches)
	}
}

func TestE6005_UndefinedCondition(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
Outputs:
  BucketName:
    Condition: NonExistentCondition
    Value: !Ref MyBucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E6005{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for undefined condition, got %d", len(matches))
	}
}

func TestE6005_NoCondition(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
Outputs:
  BucketName:
    Value: !Ref MyBucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E6005{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for no condition, got %d: %v", len(matches), matches)
	}
}

func TestE6005_MultipleOutputs(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Conditions:
  IsProd:
    Fn::Equals:
      - true
      - true
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
Outputs:
  ValidOutput:
    Condition: IsProd
    Value: !Ref MyBucket
  InvalidOutput1:
    Condition: NonExistent1
    Value: value1
  InvalidOutput2:
    Condition: NonExistent2
    Value: value2
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E6005{}
	matches := rule.Match(parsed)

	if len(matches) != 2 {
		t.Errorf("Expected 2 matches for multiple undefined conditions, got %d", len(matches))
	}
}
