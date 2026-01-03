package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW3663_PermissionWithSourceAccount(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      FunctionName: MyFunction
      Runtime: python3.9
      Handler: index.handler
      Code:
        S3Bucket: my-bucket
        S3Key: code.zip
      Role: !GetAtt MyRole.Arn
  MyPermission:
    Type: AWS::Lambda::Permission
    Properties:
      FunctionName: !Ref MyFunction
      Action: lambda:InvokeFunction
      Principal: s3.amazonaws.com
      SourceAccount: !Ref AWS::AccountId
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

	rule := &W3663{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for permission with SourceAccount, got %d: %v", len(matches), matches)
	}
}

func TestW3663_PermissionWithoutSourceAccount(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      FunctionName: MyFunction
      Runtime: python3.9
      Handler: index.handler
      Code:
        S3Bucket: my-bucket
        S3Key: code.zip
      Role: !GetAtt MyRole.Arn
  MyPermission:
    Type: AWS::Lambda::Permission
    Properties:
      FunctionName: !Ref MyFunction
      Action: lambda:InvokeFunction
      Principal: s3.amazonaws.com
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

	rule := &W3663{}
	matches := rule.Match(parsed)

	if len(matches) == 0 {
		t.Errorf("Expected matches for permission without SourceAccount, got 0")
	}
}
