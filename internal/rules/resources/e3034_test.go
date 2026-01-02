package resources

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3034_ValidNumberRange(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      FunctionName: my-function
      Runtime: python3.12
      Handler: index.handler
      Role: arn:aws:iam::123456789012:role/MyRole
      MemorySize: 256
      Timeout: 30
      Code:
        S3Bucket: my-bucket
        S3Key: code.zip
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3034{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid number range, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE3034_MemorySizeTooLow(t *testing.T) {
	// Lambda MemorySize min is 128
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      FunctionName: my-function
      Runtime: python3.12
      Handler: index.handler
      Role: arn:aws:iam::123456789012:role/MyRole
      MemorySize: 64
      Code:
        S3Bucket: my-bucket
        S3Key: code.zip
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3034{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for MemorySize below minimum 128")
	} else {
		if !strings.Contains(matches[0].Message, "less than minimum") {
			t.Errorf("Error message should mention minimum: %s", matches[0].Message)
		}
	}
}

func TestE3034_MemorySizeTooHigh(t *testing.T) {
	// Lambda MemorySize max is 10240
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      FunctionName: my-function
      Runtime: python3.12
      Handler: index.handler
      Role: arn:aws:iam::123456789012:role/MyRole
      MemorySize: 20000
      Code:
        S3Bucket: my-bucket
        S3Key: code.zip
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3034{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for MemorySize above maximum 10240")
	} else {
		if !strings.Contains(matches[0].Message, "exceeds maximum") {
			t.Errorf("Error message should mention maximum: %s", matches[0].Message)
		}
	}
}

func TestE3034_TimeoutTooHigh(t *testing.T) {
	// Lambda Timeout max is 900
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      FunctionName: my-function
      Runtime: python3.12
      Handler: index.handler
      Role: arn:aws:iam::123456789012:role/MyRole
      Timeout: 1000
      Code:
        S3Bucket: my-bucket
        S3Key: code.zip
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3034{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for Timeout above maximum 900")
	}
}

func TestE3034_SQSDelaySecondsRange(t *testing.T) {
	// SQS DelaySeconds range is 0-900
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyQueue:
    Type: AWS::SQS::Queue
    Properties:
      QueueName: my-queue
      DelaySeconds: 1000
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3034{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for DelaySeconds above maximum 900")
	}
}

func TestE3034_IAMMaxSessionDurationRange(t *testing.T) {
	// IAM Role MaxSessionDuration range is 3600-43200
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyRole:
    Type: AWS::IAM::Role
    Properties:
      RoleName: my-role
      AssumeRolePolicyDocument:
        Version: "2012-10-17"
        Statement: []
      MaxSessionDuration: 1000
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3034{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for MaxSessionDuration below minimum 3600")
	}
}

func TestE3034_NonNumericSkipped(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      FunctionName: my-function
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

	rule := &E3034{}
	matches := rule.Match(tmpl)

	// String properties should be skipped
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for non-numeric properties, got %d", len(matches))
	}
}

func TestE3034_UnknownResourceSkipped(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyCustom:
    Type: AWS::Unknown::Resource
    Properties:
      SomeNumber: -999
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3034{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for unknown resource type, got %d", len(matches))
	}
}

func TestE3034_Metadata(t *testing.T) {
	rule := &E3034{}

	if rule.ID() != "E3034" {
		t.Errorf("Expected ID E3034, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
