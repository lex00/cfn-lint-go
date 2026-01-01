package functions

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE1019_ValidSub(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  Environment:
    Type: String
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !Sub "my-bucket-${Environment}-${AWS::Region}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1019{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid Sub, got %d: %v", len(matches), matches)
	}
}

func TestE1019_UndefinedVariable(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !Sub "my-bucket-${UndefinedParam}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1019{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for undefined variable, got %d", len(matches))
	}
}

func TestE1019_ValidGetAttSyntax(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyRole:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument:
        Version: "2012-10-17"
        Statement: []
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Role: !Sub "${MyRole.Arn}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1019{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid GetAtt syntax in Sub, got %d: %v", len(matches), matches)
	}
}

func TestE1019_InvalidGetAttSyntax(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Role: !Sub "${NonExistentRole.Arn}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1019{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for undefined resource in GetAtt syntax, got %d", len(matches))
	}
}

func TestE1019_SubWithVariableMap(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName:
        Fn::Sub:
          - "my-bucket-${CustomVar}"
          - CustomVar: my-value
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1019{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for Sub with variable map, got %d: %v", len(matches), matches)
	}
}

func TestE1019_PseudoParameters(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !Sub "bucket-${AWS::AccountId}-${AWS::Region}-${AWS::StackName}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1019{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for pseudo-parameters, got %d: %v", len(matches), matches)
	}
}
