package resources

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3033_ValidStringLength(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      FunctionName: my-function
      Description: A short description
      Runtime: python3.12
      Handler: index.handler
      Role: arn:aws:iam::123456789012:role/MyRole
      Code:
        S3Bucket: my-bucket
        S3Key: code.zip
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3033{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid string length, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE3033_StringTooLong(t *testing.T) {
	// Lambda FunctionName max is 64 characters
	longName := "this-function-name-is-way-too-long-and-exceeds-the-maximum-allowed-length-of-64-characters"
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      FunctionName: ` + longName + `
      Runtime: python3.12
      Handler: index.handler
      Role: arn:aws:iam::123456789012:role/MyRole
      Code:
        S3Bucket: my-bucket
        S3Key: code.zip
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3033{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for function name exceeding 64 characters")
	} else {
		if !strings.Contains(matches[0].Message, "maximum") {
			t.Errorf("Error message should mention maximum: %s", matches[0].Message)
		}
	}
}

func TestE3033_StringTooShort(t *testing.T) {
	// S3 BucketName min is 3 characters
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: ab
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3033{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for bucket name less than 3 characters")
	} else {
		if !strings.Contains(matches[0].Message, "minimum") {
			t.Errorf("Error message should mention minimum: %s", matches[0].Message)
		}
	}
}

func TestE3033_ValidS3BucketName(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: my-valid-bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3033{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid bucket name length, got %d", len(matches))
	}
}

func TestE3033_LambdaDescriptionTooLong(t *testing.T) {
	// Lambda Description max is 256 characters
	longDesc := strings.Repeat("x", 300)
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      FunctionName: my-function
      Description: ` + longDesc + `
      Runtime: python3.12
      Handler: index.handler
      Role: arn:aws:iam::123456789012:role/MyRole
      Code:
        S3Bucket: my-bucket
        S3Key: code.zip
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3033{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for description exceeding 256 characters")
	}
}

func TestE3033_UnknownResourceSkipped(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyCustom:
    Type: AWS::Unknown::Resource
    Properties:
      SomeProp: ab
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3033{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for unknown resource type, got %d", len(matches))
	}
}

func TestE3033_Metadata(t *testing.T) {
	rule := &E3033{}

	if rule.ID() != "E3033" {
		t.Errorf("Expected ID E3033, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
