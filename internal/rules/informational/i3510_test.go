package informational

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestI3510_MatchingServices(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyRole:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument:
        Statement:
          - Effect: Allow
            Principal:
              Service: lambda.amazonaws.com
            Action: sts:AssumeRole
      Policies:
        - PolicyName: MyPolicy
          PolicyDocument:
            Statement:
              - Effect: Allow
                Action: s3:GetObject
                Resource: arn:aws:s3:::my-bucket/*
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I3510{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for matching services, got %d", len(matches))
	}
}

func TestI3510_MismatchedServices(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyRole:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument:
        Statement:
          - Effect: Allow
            Principal:
              Service: lambda.amazonaws.com
            Action: sts:AssumeRole
      Policies:
        - PolicyName: MyPolicy
          PolicyDocument:
            Statement:
              - Effect: Allow
                Action: s3:GetObject
                Resource: arn:aws:dynamodb:us-east-1:123456789012:table/MyTable
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I3510{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for mismatched services (s3 action with dynamodb resource)")
	}
}

func TestI3510_WildcardResource(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyRole:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument:
        Statement:
          - Effect: Allow
            Principal:
              Service: lambda.amazonaws.com
            Action: sts:AssumeRole
      Policies:
        - PolicyName: MyPolicy
          PolicyDocument:
            Statement:
              - Effect: Allow
                Action: s3:GetObject
                Resource: '*'
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I3510{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for wildcard resource, got %d", len(matches))
	}
}

func TestI3510_Metadata(t *testing.T) {
	rule := &I3510{}

	if rule.ID() != "I3510" {
		t.Errorf("Expected ID I3510, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
