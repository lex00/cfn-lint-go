package template

import (
	"testing"
)

func TestParse(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Description: Test template

Parameters:
  Environment:
    Type: String
    Default: dev

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: my-bucket

  MyRole:
    Type: AWS::IAM::Role
    DependsOn: MyBucket
    Properties:
      RoleName: my-role

Outputs:
  BucketArn:
    Value: !GetAtt MyBucket.Arn
`

	tmpl, err := Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if tmpl.AWSTemplateFormatVersion != "2010-09-09" {
		t.Errorf("Expected version 2010-09-09, got %s", tmpl.AWSTemplateFormatVersion)
	}

	if tmpl.Description != "Test template" {
		t.Errorf("Expected description 'Test template', got %s", tmpl.Description)
	}

	if len(tmpl.Parameters) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(tmpl.Parameters))
	}

	if _, ok := tmpl.Parameters["Environment"]; !ok {
		t.Error("Expected parameter 'Environment' not found")
	}

	if len(tmpl.Resources) != 2 {
		t.Errorf("Expected 2 resources, got %d", len(tmpl.Resources))
	}

	bucket, ok := tmpl.Resources["MyBucket"]
	if !ok {
		t.Error("Expected resource 'MyBucket' not found")
	}
	if bucket.Type != "AWS::S3::Bucket" {
		t.Errorf("Expected type AWS::S3::Bucket, got %s", bucket.Type)
	}

	role, ok := tmpl.Resources["MyRole"]
	if !ok {
		t.Error("Expected resource 'MyRole' not found")
	}
	if len(role.DependsOn) != 1 || role.DependsOn[0] != "MyBucket" {
		t.Errorf("Expected DependsOn [MyBucket], got %v", role.DependsOn)
	}

	if len(tmpl.Outputs) != 1 {
		t.Errorf("Expected 1 output, got %d", len(tmpl.Outputs))
	}
}

func TestHasResource(t *testing.T) {
	yaml := `
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, _ := Parse([]byte(yaml))

	if !tmpl.HasResource("MyBucket") {
		t.Error("Expected HasResource('MyBucket') to be true")
	}

	if tmpl.HasResource("NonExistent") {
		t.Error("Expected HasResource('NonExistent') to be false")
	}
}

func TestHasParameter(t *testing.T) {
	yaml := `
Parameters:
  MyParam:
    Type: String
Resources:
  Dummy:
    Type: AWS::CloudFormation::WaitConditionHandle
`
	tmpl, _ := Parse([]byte(yaml))

	if !tmpl.HasParameter("MyParam") {
		t.Error("Expected HasParameter('MyParam') to be true")
	}

	if tmpl.HasParameter("NonExistent") {
		t.Error("Expected HasParameter('NonExistent') to be false")
	}
}
