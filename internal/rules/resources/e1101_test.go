package resources

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE1101_ValidProperties(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: my-bucket
      Tags:
        - Key: Env
          Value: prod
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E1101{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid properties, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE1101_UnknownProperty(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: my-bucket
      InvalidProperty: some-value
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E1101{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for unknown property")
	} else {
		if !strings.Contains(matches[0].Message, "unknown property") {
			t.Errorf("Error message should mention 'unknown property': %s", matches[0].Message)
		}
		if !strings.Contains(matches[0].Message, "InvalidProperty") {
			t.Errorf("Error message should mention property name: %s", matches[0].Message)
		}
	}
}

func TestE1101_MultipleUnknownProperties(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: my-bucket
      FakeProperty1: value1
      FakeProperty2: value2
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E1101{}
	matches := rule.Match(tmpl)

	if len(matches) < 2 {
		t.Errorf("Expected at least 2 matches for multiple unknown properties, got %d", len(matches))
	}
}

func TestE1101_UnknownResourceSkipped(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyCustom:
    Type: AWS::Unknown::Resource
    Properties:
      AnyProperty: any-value
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E1101{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for unknown resource type, got %d", len(matches))
	}
}

func TestE1101_LambdaFunction(t *testing.T) {
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
      NotARealProperty: oops
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E1101{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for unknown property on Lambda function")
	}

	// Check that we flagged the unknown property, not the valid ones
	foundUnknown := false
	for _, m := range matches {
		if strings.Contains(m.Message, "NotARealProperty") {
			foundUnknown = true
		}
	}
	if !foundUnknown {
		t.Error("Expected to find match for NotARealProperty")
	}
}

func TestE1101_Metadata(t *testing.T) {
	rule := &E1101{}

	if rule.ID() != "E1101" {
		t.Errorf("Expected ID E1101, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
