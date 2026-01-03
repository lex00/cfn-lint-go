package resources

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE2533_ValidZipPackageWithRuntime(t *testing.T) {
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
      PackageType: Zip
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2533{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid Zip package with Runtime, got %d", len(matches))
	}
}

func TestE2533_ValidCustomRuntime(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Runtime: provided.al2023
      Handler: bootstrap
      Code:
        S3Bucket: my-bucket
        S3Key: function.zip
      Role: arn:aws:iam::123456789012:role/lambda-role
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2533{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid custom runtime, got %d", len(matches))
	}
}

func TestE2533_InvalidImagePackageWithRuntime(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Runtime: python3.11
      Code:
        ImageUri: 123456789012.dkr.ecr.us-east-1.amazonaws.com/my-function:latest
      Role: arn:aws:iam::123456789012:role/lambda-role
      PackageType: Image
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2533{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for Image package with Runtime specified (incompatible)")
	}

	if len(matches) > 0 && !strings.Contains(matches[0].Message, "PackageType 'Image'") {
		t.Errorf("Expected error message to mention PackageType Image, got '%s'", matches[0].Message)
	}
}

func TestE2533_InvalidUnrecognizedRuntime(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Runtime: python4.0
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

	rule := &E2533{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for unrecognized runtime 'python4.0'")
	}

	if len(matches) > 0 && !strings.Contains(matches[0].Message, "unrecognized runtime") {
		t.Errorf("Expected error message to mention unrecognized runtime, got '%s'", matches[0].Message)
	}
}

func TestE2533_InvalidUnrecognizedRuntimeNodeJS100(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Runtime: nodejs100.x
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

	rule := &E2533{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for unrecognized runtime 'nodejs100.x'")
	}
}

func TestE2533_ValidNoRuntime(t *testing.T) {
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

	rule := &E2533{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when Runtime is not specified, got %d", len(matches))
	}
}

func TestE2533_ValidImagePackageWithoutRuntime(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Code:
        ImageUri: 123456789012.dkr.ecr.us-east-1.amazonaws.com/my-function:latest
      Role: arn:aws:iam::123456789012:role/lambda-role
      PackageType: Image
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2533{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for Image package without Runtime, got %d", len(matches))
	}
}

func TestE2533_SkipsIntrinsicFunctions(t *testing.T) {
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

	rule := &E2533{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when Runtime is an intrinsic function, got %d", len(matches))
	}
}

func TestE2533_ValidAllKnownRuntimes(t *testing.T) {
	knownRuntimes := []string{
		"python3.11",
		"python3.12",
		"nodejs18.x",
		"nodejs20.x",
		"java11",
		"java17",
		"java21",
		"dotnet6",
		"dotnet8",
		"ruby3.2",
		"provided.al2",
		"provided.al2023",
	}

	for _, runtime := range knownRuntimes {
		yaml := fmt.Sprintf(`
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Runtime: %s
      Handler: index.handler
      Code:
        S3Bucket: my-bucket
        S3Key: function.zip
      Role: arn:aws:iam::123456789012:role/lambda-role
`, runtime)

		tmpl, err := template.Parse([]byte(yaml))
		if err != nil {
			t.Fatalf("Failed to parse template for runtime %s: %v", runtime, err)
		}

		rule := &E2533{}
		matches := rule.Match(tmpl)

		if len(matches) != 0 {
			t.Errorf("Expected 0 matches for known runtime '%s', got %d: %v", runtime, len(matches), matches)
		}
	}
}

func TestE2533_Metadata(t *testing.T) {
	rule := &E2533{}

	if rule.ID() != "E2533" {
		t.Errorf("Expected ID E2533, got %s", rule.ID())
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
