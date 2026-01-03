package informational

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestI3042_HardcodedARN(t *testing.T) {
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
                Resource: arn:aws:s3:us-east-1:123456789012:bucket/my-bucket/*
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I3042{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for hardcoded ARN with account ID and region")
	}
}

func TestI3042_PseudoParameterARN(t *testing.T) {
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
                Resource: !Sub 'arn:${AWS::Partition}:s3:${AWS::Region}:${AWS::AccountId}:bucket/my-bucket/*'
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I3042{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for ARN using pseudo parameters, got %d", len(matches))
	}
}

func TestI3042_NonARNResource(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: my-bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I3042{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for resource without ARNs, got %d", len(matches))
	}
}

func TestI3042_Metadata(t *testing.T) {
	rule := &I3042{}

	if rule.ID() != "I3042" {
		t.Errorf("Expected ID I3042, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
