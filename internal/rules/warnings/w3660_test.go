package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW3660_SingleApiResource(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyApi:
    Type: AWS::ApiGateway::RestApi
    Properties:
      Name: MyApi
  MyResource:
    Type: AWS::ApiGateway::Resource
    Properties:
      RestApiId: !Ref MyApi
      ParentId: !GetAtt MyApi.RootResourceId
      PathPart: items
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W3660{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for single API resource, got %d: %v", len(matches), matches)
	}
}

func TestW3660_DeploymentWithoutDependsOn(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyApi:
    Type: AWS::ApiGateway::RestApi
    Properties:
      Name: MyApi
  MyResource:
    Type: AWS::ApiGateway::Resource
    Properties:
      RestApiId: !Ref MyApi
      ParentId: !GetAtt MyApi.RootResourceId
      PathPart: items
  MyDeployment:
    Type: AWS::ApiGateway::Deployment
    Properties:
      RestApiId: !Ref MyApi
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W3660{}
	matches := rule.Match(parsed)

	if len(matches) == 0 {
		t.Errorf("Expected matches for deployment without DependsOn, got 0")
	}
}
