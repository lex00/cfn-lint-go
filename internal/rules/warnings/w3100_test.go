package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW3100_Metadata(t *testing.T) {
	rule := &W3100{}

	if rule.ID() != "W3100" {
		t.Errorf("Expected ID W3100, got %s", rule.ID())
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

func TestW3100_MissingMemorySize(t *testing.T) {
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

	rule := &W3100{}
	matches := rule.Match(tmpl)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for missing MemorySize, got %d", len(matches))
	}

	if len(matches) > 0 && matches[0].Path[1] != "MyFunction" {
		t.Errorf("Expected match for MyFunction, got %v", matches[0].Path)
	}
}

func TestW3100_WithMemorySize(t *testing.T) {
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
      MemorySize: 256
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &W3100{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when MemorySize is set, got %d", len(matches))
	}
}

func TestW3100_NonSAMFunction(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
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

	rule := &W3100{}
	matches := rule.Match(tmpl)

	// Should not match non-SAM functions
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for non-SAM function, got %d", len(matches))
	}
}

func TestW3100_GlobalMemorySize(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Globals:
  Function:
    MemorySize: 512
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

	rule := &W3100{}
	matches := rule.Match(tmpl)

	// When Globals has MemorySize, function-level is not required
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when Globals has MemorySize, got %d", len(matches))
	}
}
