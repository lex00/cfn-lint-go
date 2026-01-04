package sam

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestTransform_BasicSAMFunction(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyFunction:
    Type: AWS::Serverless::Function
    Properties:
      Runtime: python3.9
      Handler: index.handler
      CodeUri: ./src
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	result, err := Transform(tmpl, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil TransformResult")
	}

	if result.Template == nil {
		t.Fatal("Expected non-nil transformed template")
	}

	// The transformed template should have Lambda function instead of SAM function
	hasLambda := false
	for _, res := range result.Template.Resources {
		if res.Type == "AWS::Lambda::Function" {
			hasLambda = true
			break
		}
	}
	if !hasLambda {
		t.Error("Expected transformed template to contain AWS::Lambda::Function")
	}

	// Should not have SAM resources anymore
	for name, res := range result.Template.Resources {
		if IsSAMResourceType(res.Type) {
			t.Errorf("Transformed template should not contain SAM resource type, found %s: %s", name, res.Type)
		}
	}
}

func TestTransform_WithSourceMap(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyFunction:
    Type: AWS::Serverless::Function
    Properties:
      Runtime: python3.9
      Handler: index.handler
      CodeUri: ./src
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	result, err := Transform(tmpl, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	if result.SourceMap == nil {
		t.Fatal("Expected non-nil SourceMap")
	}

	// Check that generated resources map back to original SAM resource
	// The Lambda function should map back to MyFunction
	foundMapping := false
	for cfnName := range result.Template.Resources {
		loc, ok := result.SourceMap.GetResourceLocation(cfnName)
		if ok && loc.OriginalResource == "MyFunction" {
			foundMapping = true
			break
		}
	}
	if !foundMapping {
		t.Error("Expected at least one CFN resource to map back to MyFunction")
	}
}

func TestTransform_NonSAMTemplate(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	result, err := Transform(tmpl, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Non-SAM templates should pass through unchanged
	if result.Template == nil {
		t.Fatal("Expected non-nil template")
	}

	if len(result.Template.Resources) != 1 {
		t.Errorf("Expected 1 resource, got %d", len(result.Template.Resources))
	}

	if _, ok := result.Template.Resources["MyBucket"]; !ok {
		t.Error("Expected MyBucket resource to be preserved")
	}
}

func TestTransform_NilTemplate(t *testing.T) {
	result, err := Transform(nil, nil)
	if err == nil {
		t.Error("Expected error for nil template")
	}
	if result != nil {
		t.Error("Expected nil result for nil template")
	}
}

func TestTransform_WithOptions(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyFunction:
    Type: AWS::Serverless::Function
    Properties:
      Runtime: python3.9
      Handler: index.handler
      CodeUri: ./src
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	opts := &TransformOptions{
		Region:    "us-west-2",
		AccountID: "123456789012",
		StackName: "test-stack",
	}

	result, err := Transform(tmpl, opts)
	if err != nil {
		t.Fatalf("Transform with options failed: %v", err)
	}

	if result.Template == nil {
		t.Fatal("Expected non-nil transformed template")
	}
}

func TestTransform_SAMApi(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyApi:
    Type: AWS::Serverless::Api
    Properties:
      StageName: prod
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	result, err := Transform(tmpl, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have API Gateway resources
	hasApiGateway := false
	for _, res := range result.Template.Resources {
		if res.Type == "AWS::ApiGateway::RestApi" {
			hasApiGateway = true
			break
		}
	}
	if !hasApiGateway {
		t.Error("Expected transformed template to contain AWS::ApiGateway::RestApi")
	}
}

func TestTransform_SAMSimpleTable(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyTable:
    Type: AWS::Serverless::SimpleTable
    Properties:
      PrimaryKey:
        Name: id
        Type: String
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	result, err := Transform(tmpl, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have DynamoDB table
	hasDynamoDB := false
	for _, res := range result.Template.Resources {
		if res.Type == "AWS::DynamoDB::Table" {
			hasDynamoDB = true
			break
		}
	}
	if !hasDynamoDB {
		t.Error("Expected transformed template to contain AWS::DynamoDB::Table")
	}
}

func TestTransformResult_HasWarnings(t *testing.T) {
	result := &TransformResult{
		Warnings: []string{"Warning 1", "Warning 2"},
	}

	if len(result.Warnings) != 2 {
		t.Errorf("Expected 2 warnings, got %d", len(result.Warnings))
	}
}

func TestTransformOptions_Defaults(t *testing.T) {
	opts := DefaultTransformOptions()

	if opts.Region == "" {
		t.Error("Expected default region to be set")
	}
	if opts.AccountID == "" {
		t.Error("Expected default account ID to be set")
	}
}
