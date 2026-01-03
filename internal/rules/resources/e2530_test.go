package resources

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE2530_ValidJava11WithSnapStart(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Runtime: java11
      Handler: com.example.Handler
      Code:
        S3Bucket: my-bucket
        S3Key: function.jar
      Role: arn:aws:iam::123456789012:role/lambda-role
      SnapStart:
        ApplyOn: PublishedVersions
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2530{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for Java 11 with SnapStart, got %d", len(matches))
	}
}

func TestE2530_ValidJava17WithSnapStart(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Runtime: java17
      Handler: com.example.Handler
      Code:
        S3Bucket: my-bucket
        S3Key: function.jar
      Role: arn:aws:iam::123456789012:role/lambda-role
      SnapStart:
        ApplyOn: PublishedVersions
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2530{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for Java 17 with SnapStart, got %d", len(matches))
	}
}

func TestE2530_ValidJava21WithSnapStart(t *testing.T) {
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
      SnapStart:
        ApplyOn: PublishedVersions
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2530{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for Java 21 with SnapStart, got %d", len(matches))
	}
}

func TestE2530_InvalidPython3WithSnapStart(t *testing.T) {
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
      SnapStart:
        ApplyOn: PublishedVersions
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2530{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for Python runtime with SnapStart (incompatible)")
	}

	if len(matches) > 0 && !strings.Contains(matches[0].Message, "python3.11") {
		t.Errorf("Expected error message to mention runtime 'python3.11', got '%s'", matches[0].Message)
	}
}

func TestE2530_InvalidNodeJSWithSnapStart(t *testing.T) {
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
      SnapStart:
        ApplyOn: PublishedVersions
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2530{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for Node.js runtime with SnapStart (incompatible)")
	}
}

func TestE2530_InvalidJava8WithSnapStart(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Runtime: java8.al2
      Handler: com.example.Handler
      Code:
        S3Bucket: my-bucket
        S3Key: function.jar
      Role: arn:aws:iam::123456789012:role/lambda-role
      SnapStart:
        ApplyOn: PublishedVersions
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2530{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for Java 8 runtime with SnapStart (incompatible)")
	}
}

func TestE2530_ValidSnapStartDisabled(t *testing.T) {
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
      SnapStart:
        ApplyOn: None
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2530{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when SnapStart is disabled (ApplyOn: None), got %d", len(matches))
	}
}

func TestE2530_ValidNoSnapStart(t *testing.T) {
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

	rule := &E2530{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when SnapStart is not configured, got %d", len(matches))
	}
}

func TestE2530_InvalidNoRuntimeWithSnapStart(t *testing.T) {
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
      SnapStart:
        ApplyOn: PublishedVersions
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2530{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for SnapStart enabled without Runtime specified")
	}
}

func TestE2530_SkipsIntrinsicFunctions(t *testing.T) {
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
      SnapStart:
        ApplyOn: PublishedVersions
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2530{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when Runtime is an intrinsic function, got %d", len(matches))
	}
}

func TestE2530_Metadata(t *testing.T) {
	rule := &E2530{}

	if rule.ID() != "E2530" {
		t.Errorf("Expected ID E2530, got %s", rule.ID())
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
