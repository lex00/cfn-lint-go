package resources

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3633_ValidKinesisWithStartingPosition(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyEventSourceMapping:
    Type: AWS::Lambda::EventSourceMapping
    Properties:
      EventSourceArn: arn:aws:kinesis:us-east-1:123456789012:stream/my-stream
      FunctionName: my-function
      StartingPosition: LATEST
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3633{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid Kinesis with StartingPosition, got %d", len(matches))
	}
}

func TestE3633_InvalidKinesisMissingStartingPosition(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyEventSourceMapping:
    Type: AWS::Lambda::EventSourceMapping
    Properties:
      EventSourceArn: arn:aws:kinesis:us-east-1:123456789012:stream/my-stream
      FunctionName: my-function
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3633{}
	matches := rule.Match(tmpl)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for Kinesis missing StartingPosition, got %d", len(matches))
	}
	if len(matches) > 0 && !containsString(matches[0].Message, "StartingPosition") {
		t.Errorf("Expected error about StartingPosition, got: %s", matches[0].Message)
	}
}

func TestE3633_ValidDynamoDBWithStartingPosition(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyEventSourceMapping:
    Type: AWS::Lambda::EventSourceMapping
    Properties:
      EventSourceArn: arn:aws:dynamodb:us-east-1:123456789012:table/my-table/stream/2023-01-01T00:00:00.000
      FunctionName: my-function
      StartingPosition: TRIM_HORIZON
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3633{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid DynamoDB with StartingPosition, got %d", len(matches))
	}
}

func TestE3633_InvalidDynamoDBMissingStartingPosition(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyEventSourceMapping:
    Type: AWS::Lambda::EventSourceMapping
    Properties:
      EventSourceArn: arn:aws:dynamodb:us-east-1:123456789012:table/my-table/stream/2023-01-01T00:00:00.000
      FunctionName: my-function
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3633{}
	matches := rule.Match(tmpl)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for DynamoDB missing StartingPosition, got %d", len(matches))
	}
}

func TestE3633_SQSDoesNotRequireStartingPosition(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyEventSourceMapping:
    Type: AWS::Lambda::EventSourceMapping
    Properties:
      EventSourceArn: arn:aws:sqs:us-east-1:123456789012:my-queue
      FunctionName: my-function
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3633{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for SQS without StartingPosition, got %d", len(matches))
	}
}

func TestE3633_IntrinsicFunction(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyEventSourceMapping:
    Type: AWS::Lambda::EventSourceMapping
    Properties:
      EventSourceArn:
        Fn::GetAtt: [MyKinesisStream, Arn]
      FunctionName: my-function
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3633{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for intrinsic function, got %d", len(matches))
	}
}

func TestE3633_Metadata(t *testing.T) {
	rule := &E3633{}

	if rule.ID() != "E3633" {
		t.Errorf("Expected ID E3633, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("Expected non-empty ShortDesc")
	}

	if rule.Description() == "" {
		t.Error("Expected non-empty Description")
	}

	if rule.Source() == "" {
		t.Error("Expected non-empty Source")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Expected non-empty Tags")
	}
}
