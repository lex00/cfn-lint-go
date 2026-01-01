package functions

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE1010_ValidGetAtt(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
  MyRole:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument:
        Version: "2012-10-17"
        Statement:
          - Effect: Allow
            Principal:
              Service: lambda.amazonaws.com
            Action: sts:AssumeRole
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Role: !GetAtt MyRole.Arn
      Code:
        S3Bucket: !GetAtt MyBucket.DomainName
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1010{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid GetAtt, got %d: %v", len(matches), matches)
	}
}

func TestE1010_UndefinedResource(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Role: !GetAtt NonExistentRole.Arn
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1010{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for undefined GetAtt resource, got %d", len(matches))
	}
}

func TestE1010_GetAttArrayFormat(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Role:
        Fn::GetAtt:
          - NonExistentRole
          - Arn
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1010{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for undefined GetAtt in array format, got %d", len(matches))
	}
}

func TestE1010_OutputGetAtt(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
Outputs:
  BucketArn:
    Value: !GetAtt NonExistentBucket.Arn
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1010{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for undefined GetAtt in output, got %d", len(matches))
	}
}
