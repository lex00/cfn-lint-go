package sam

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestIsSAMTemplate_WithTransformString(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	if !IsSAMTemplate(tmpl) {
		t.Error("Expected IsSAMTemplate to return true for template with AWS::Serverless-2016-10-31 transform")
	}
}

func TestIsSAMTemplate_WithTransformArray(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Transform:
  - AWS::Serverless-2016-10-31
  - AWS::LanguageExtensions
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	if !IsSAMTemplate(tmpl) {
		t.Error("Expected IsSAMTemplate to return true for template with AWS::Serverless-2016-10-31 in transform array")
	}
}

func TestIsSAMTemplate_WithServerlessFunction(t *testing.T) {
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

	if !IsSAMTemplate(tmpl) {
		t.Error("Expected IsSAMTemplate to return true for template with AWS::Serverless::Function")
	}
}

func TestIsSAMTemplate_WithServerlessApi(t *testing.T) {
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

	if !IsSAMTemplate(tmpl) {
		t.Error("Expected IsSAMTemplate to return true for template with AWS::Serverless::Api")
	}
}

func TestIsSAMTemplate_WithServerlessHttpApi(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyHttpApi:
    Type: AWS::Serverless::HttpApi
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	if !IsSAMTemplate(tmpl) {
		t.Error("Expected IsSAMTemplate to return true for template with AWS::Serverless::HttpApi")
	}
}

func TestIsSAMTemplate_WithServerlessSimpleTable(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyTable:
    Type: AWS::Serverless::SimpleTable
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	if !IsSAMTemplate(tmpl) {
		t.Error("Expected IsSAMTemplate to return true for template with AWS::Serverless::SimpleTable")
	}
}

func TestIsSAMTemplate_WithServerlessLayerVersion(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyLayer:
    Type: AWS::Serverless::LayerVersion
    Properties:
      ContentUri: ./layer
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	if !IsSAMTemplate(tmpl) {
		t.Error("Expected IsSAMTemplate to return true for template with AWS::Serverless::LayerVersion")
	}
}

func TestIsSAMTemplate_WithServerlessApplication(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyApp:
    Type: AWS::Serverless::Application
    Properties:
      Location: ./nested
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	if !IsSAMTemplate(tmpl) {
		t.Error("Expected IsSAMTemplate to return true for template with AWS::Serverless::Application")
	}
}

func TestIsSAMTemplate_WithServerlessStateMachine(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyStateMachine:
    Type: AWS::Serverless::StateMachine
    Properties:
      Definition:
        StartAt: Hello
        States:
          Hello:
            Type: Pass
            End: true
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	if !IsSAMTemplate(tmpl) {
		t.Error("Expected IsSAMTemplate to return true for template with AWS::Serverless::StateMachine")
	}
}

func TestIsSAMTemplate_WithServerlessConnector(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyConnector:
    Type: AWS::Serverless::Connector
    Properties:
      Source:
        Id: MyFunction
      Destination:
        Id: MyTable
      Permissions:
        - Read
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	if !IsSAMTemplate(tmpl) {
		t.Error("Expected IsSAMTemplate to return true for template with AWS::Serverless::Connector")
	}
}

func TestIsSAMTemplate_WithServerlessGraphQLApi(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyGraphQLApi:
    Type: AWS::Serverless::GraphQLApi
    Properties:
      SchemaUri: ./schema.graphql
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	if !IsSAMTemplate(tmpl) {
		t.Error("Expected IsSAMTemplate to return true for template with AWS::Serverless::GraphQLApi")
	}
}

func TestIsSAMTemplate_WithServerlessResourceButNoTransform(t *testing.T) {
	// This is an invalid SAM template (missing transform) but we should still detect it
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
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

	if !IsSAMTemplate(tmpl) {
		t.Error("Expected IsSAMTemplate to return true for template with AWS::Serverless resource even without transform")
	}
}

func TestIsSAMTemplate_RegularCloudFormation(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Runtime: python3.9
      Handler: index.handler
      Code:
        S3Bucket: my-bucket
        S3Key: code.zip
      Role: !GetAtt MyRole.Arn
  MyRole:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument:
        Version: '2012-10-17'
        Statement:
          - Effect: Allow
            Principal:
              Service: lambda.amazonaws.com
            Action: sts:AssumeRole
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	if IsSAMTemplate(tmpl) {
		t.Error("Expected IsSAMTemplate to return false for regular CloudFormation template")
	}
}

func TestIsSAMTemplate_WithOtherTransform(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::LanguageExtensions
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	if IsSAMTemplate(tmpl) {
		t.Error("Expected IsSAMTemplate to return false for template with non-SAM transform")
	}
}

func TestIsSAMTemplate_EmptyTemplate(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources: {}
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	if IsSAMTemplate(tmpl) {
		t.Error("Expected IsSAMTemplate to return false for empty template")
	}
}

func TestIsSAMTemplate_NilTemplate(t *testing.T) {
	if IsSAMTemplate(nil) {
		t.Error("Expected IsSAMTemplate to return false for nil template")
	}
}

func TestIsSAMResourceType(t *testing.T) {
	samTypes := []string{
		"AWS::Serverless::Function",
		"AWS::Serverless::Api",
		"AWS::Serverless::HttpApi",
		"AWS::Serverless::SimpleTable",
		"AWS::Serverless::LayerVersion",
		"AWS::Serverless::Application",
		"AWS::Serverless::StateMachine",
		"AWS::Serverless::Connector",
		"AWS::Serverless::GraphQLApi",
	}

	for _, samType := range samTypes {
		if !IsSAMResourceType(samType) {
			t.Errorf("Expected IsSAMResourceType(%q) to return true", samType)
		}
	}
}

func TestIsSAMResourceType_NonSAMTypes(t *testing.T) {
	nonSAMTypes := []string{
		"AWS::Lambda::Function",
		"AWS::S3::Bucket",
		"AWS::IAM::Role",
		"AWS::ApiGateway::RestApi",
		"AWS::DynamoDB::Table",
		"",
	}

	for _, nonSAMType := range nonSAMTypes {
		if IsSAMResourceType(nonSAMType) {
			t.Errorf("Expected IsSAMResourceType(%q) to return false", nonSAMType)
		}
	}
}
