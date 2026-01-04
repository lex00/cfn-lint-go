package informational

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestI3101_Metadata(t *testing.T) {
	rule := &I3101{}

	if rule.ID() != "I3101" {
		t.Errorf("Expected ID I3101, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	if rule.Description() == "" {
		t.Error("Description should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}

func TestI3101_SAMFunction(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: index.handler
      Runtime: nodejs18.x
      CodeUri: ./src
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I3101{}
	matches := rule.Match(tmpl)

	// Should have informational message about SAM resource expansion
	if len(matches) != 1 {
		t.Errorf("Expected 1 match for SAM Function, got %d", len(matches))
	}

	if len(matches) > 0 && matches[0].Path[1] != "MyFunction" {
		t.Errorf("Expected match for MyFunction, got %v", matches[0].Path)
	}
}

func TestI3101_MultipleSAMResources(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: index.handler
      Runtime: nodejs18.x
      CodeUri: ./src
  MyApi:
    Type: AWS::Serverless::Api
    Properties:
      StageName: prod
  MyTable:
    Type: AWS::Serverless::SimpleTable
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I3101{}
	matches := rule.Match(tmpl)

	// Should have informational message for each SAM resource
	if len(matches) != 3 {
		t.Errorf("Expected 3 matches for 3 SAM resources, got %d", len(matches))
	}
}

func TestI3101_NonSAMTemplate(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Handler: index.handler
      Runtime: nodejs18.x
      Code:
        S3Bucket: my-bucket
        S3Key: code.zip
      Role: arn:aws:iam::123456789012:role/LambdaRole
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I3101{}
	matches := rule.Match(tmpl)

	// Should not match non-SAM templates
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for non-SAM template, got %d", len(matches))
	}
}

func TestI3101_SAMResourceTypes(t *testing.T) {
	// Test each SAM resource type
	testCases := []struct {
		name         string
		resourceType string
	}{
		{"Function", "AWS::Serverless::Function"},
		{"Api", "AWS::Serverless::Api"},
		{"HttpApi", "AWS::Serverless::HttpApi"},
		{"SimpleTable", "AWS::Serverless::SimpleTable"},
		{"LayerVersion", "AWS::Serverless::LayerVersion"},
		{"Application", "AWS::Serverless::Application"},
		{"StateMachine", "AWS::Serverless::StateMachine"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyResource:
    Type: ` + tc.resourceType + `
`
			tmpl, err := template.Parse([]byte(yaml))
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			rule := &I3101{}
			matches := rule.Match(tmpl)

			if len(matches) != 1 {
				t.Errorf("Expected 1 match for %s, got %d", tc.resourceType, len(matches))
			}
		})
	}
}
