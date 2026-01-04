package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW3102_Metadata(t *testing.T) {
	rule := &W3102{}

	if rule.ID() != "W3102" {
		t.Errorf("Expected ID W3102, got %s", rule.ID())
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

func TestW3102_MissingStageName(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyApi:
    Type: AWS::Serverless::Api
    Properties:
      Name: MyAPI
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &W3102{}
	matches := rule.Match(tmpl)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for missing StageName, got %d", len(matches))
	}

	if len(matches) > 0 && matches[0].Path[1] != "MyApi" {
		t.Errorf("Expected match for MyApi, got %v", matches[0].Path)
	}
}

func TestW3102_WithStageName(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyApi:
    Type: AWS::Serverless::Api
    Properties:
      Name: MyAPI
      StageName: prod
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &W3102{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when StageName is set, got %d", len(matches))
	}
}

func TestW3102_NonSAMApi(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyApi:
    Type: AWS::ApiGateway::RestApi
    Properties:
      Name: MyAPI
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &W3102{}
	matches := rule.Match(tmpl)

	// Should not match non-SAM APIs
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for non-SAM API, got %d", len(matches))
	}
}

func TestW3102_HttpApi(t *testing.T) {
	// HttpApi doesn't require StageName
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyApi:
    Type: AWS::Serverless::HttpApi
    Properties:
      Name: MyHttpAPI
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &W3102{}
	matches := rule.Match(tmpl)

	// HttpApi is different from Api, should not trigger this rule
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for HttpApi, got %d", len(matches))
	}
}
