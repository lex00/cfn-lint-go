package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW8003_DynamicEquals(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  Environment:
    Type: String
Conditions:
  IsProd:
    Fn::Equals: [!Ref Environment, "prod"]
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Condition: IsProd
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W8003{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for dynamic Fn::Equals, got %d: %v", len(matches), matches)
	}
}

func TestW8003_AlwaysTrue(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Conditions:
  AlwaysTrue:
    Fn::Equals: ["same", "same"]
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Condition: AlwaysTrue
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W8003{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for always true condition, got %d", len(matches))
	}
}

func TestW8003_AlwaysFalse(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Conditions:
  AlwaysFalse:
    Fn::Equals: ["a", "b"]
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Condition: AlwaysFalse
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W8003{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for always false condition, got %d", len(matches))
	}
}
