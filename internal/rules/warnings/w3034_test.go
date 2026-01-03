package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW3034_ValidParameterRange(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  Timeout:
    Type: Number
    MinValue: 1
    MaxValue: 900
    Default: 30
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      FunctionName: MyFunction
      Runtime: python3.9
      Handler: index.handler
      Timeout: !Ref Timeout
      Code:
        S3Bucket: my-bucket
        S3Key: code.zip
      Role: !GetAtt MyRole.Arn
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
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W3034{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid parameter range, got %d: %v", len(matches), matches)
	}
}

func TestW3034_ParameterExceedsRange(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  Timeout:
    Type: Number
    MinValue: 0
    MaxValue: 1000
    Default: 500
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      FunctionName: MyFunction
      Runtime: python3.9
      Handler: index.handler
      Timeout: !Ref Timeout
      Code:
        S3Bucket: my-bucket
        S3Key: code.zip
      Role: !GetAtt MyRole.Arn
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
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W3034{}
	matches := rule.Match(parsed)

	// Should warn about parameter allowing values outside Lambda timeout range
	if len(matches) == 0 {
		t.Errorf("Expected matches for parameter exceeding range, got 0")
	}
}
