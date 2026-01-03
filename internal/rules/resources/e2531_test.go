package resources

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE2531_ValidCurrentRuntimePython311(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Runtime: python3.11
      Handler: index.handler
      Code:
        S3Bucket: my-bucket
        S3Key: function.zip
      Role: arn:aws:iam::123456789012:role/lambda-role
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2531{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for current Python 3.11 runtime, got %d", len(matches))
	}
}

func TestE2531_ValidCurrentRuntimeNodeJS20(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Runtime: nodejs20.x
      Handler: index.handler
      Code:
        S3Bucket: my-bucket
        S3Key: function.zip
      Role: arn:aws:iam::123456789012:role/lambda-role
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2531{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for current Node.js 20.x runtime, got %d", len(matches))
	}
}

func TestE2531_ValidCurrentRuntimeJava21(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Runtime: java21
      Handler: com.example.Handler
      Code:
        S3Bucket: my-bucket
        S3Key: function.jar
      Role: arn:aws:iam::123456789012:role/lambda-role
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2531{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for current Java 21 runtime, got %d", len(matches))
	}
}

func TestE2531_InvalidDeprecatedPython27(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Runtime: python2.7
      Handler: index.handler
      Code:
        S3Bucket: my-bucket
        S3Key: function.zip
      Role: arn:aws:iam::123456789012:role/lambda-role
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2531{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for deprecated Python 2.7 runtime")
	}

	if len(matches) > 0 && !strings.Contains(matches[0].Message, "python2.7") {
		t.Errorf("Expected error message to mention 'python2.7', got '%s'", matches[0].Message)
	}
}

func TestE2531_InvalidDeprecatedPython36(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Runtime: python3.6
      Handler: index.handler
      Code:
        S3Bucket: my-bucket
        S3Key: function.zip
      Role: arn:aws:iam::123456789012:role/lambda-role
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2531{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for deprecated Python 3.6 runtime")
	}
}

func TestE2531_InvalidDeprecatedNodeJS810(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Runtime: nodejs8.10
      Handler: index.handler
      Code:
        S3Bucket: my-bucket
        S3Key: function.zip
      Role: arn:aws:iam::123456789012:role/lambda-role
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2531{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for deprecated Node.js 8.10 runtime")
	}
}

func TestE2531_InvalidDeprecatedNodeJS16(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Runtime: nodejs16.x
      Handler: index.handler
      Code:
        S3Bucket: my-bucket
        S3Key: function.zip
      Role: arn:aws:iam::123456789012:role/lambda-role
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2531{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for deprecated Node.js 16.x runtime")
	}
}

func TestE2531_InvalidDeprecatedRuby25(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Runtime: ruby2.5
      Handler: index.handler
      Code:
        S3Bucket: my-bucket
        S3Key: function.zip
      Role: arn:aws:iam::123456789012:role/lambda-role
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2531{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for deprecated Ruby 2.5 runtime")
	}
}

func TestE2531_InvalidDeprecatedGo1x(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Runtime: go1.x
      Handler: main
      Code:
        S3Bucket: my-bucket
        S3Key: function.zip
      Role: arn:aws:iam::123456789012:role/lambda-role
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2531{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for deprecated go1.x runtime")
	}
}

func TestE2531_InvalidDeprecatedDotNetCore21(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Runtime: dotnetcore2.1
      Handler: MyApp::MyApp.Function::Handler
      Code:
        S3Bucket: my-bucket
        S3Key: function.zip
      Role: arn:aws:iam::123456789012:role/lambda-role
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2531{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for deprecated dotnetcore2.1 runtime")
	}
}

func TestE2531_ValidNoRuntime(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Handler: index.handler
      Code:
        S3Bucket: my-bucket
        S3Key: function.zip
      Role: arn:aws:iam::123456789012:role/lambda-role
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2531{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when Runtime is not specified, got %d", len(matches))
	}
}

func TestE2531_SkipsIntrinsicFunctions(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Runtime: !Ref RuntimeParameter
      Handler: index.handler
      Code:
        S3Bucket: my-bucket
        S3Key: function.zip
      Role: arn:aws:iam::123456789012:role/lambda-role
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2531{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when Runtime is an intrinsic function, got %d", len(matches))
	}
}

func TestE2531_Metadata(t *testing.T) {
	rule := &E2531{}

	if rule.ID() != "E2531" {
		t.Errorf("Expected ID E2531, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	if rule.Description() == "" {
		t.Error("Description should not be empty")
	}

	if rule.Source() == "" {
		t.Error("Source should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
