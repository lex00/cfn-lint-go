package resources

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3638_ValidPayPerRequest(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyTable:
    Type: AWS::DynamoDB::Table
    Properties:
      TableName: my-table
      BillingMode: PAY_PER_REQUEST
      AttributeDefinitions:
        - AttributeName: id
          AttributeType: S
      KeySchema:
        - AttributeName: id
          KeyType: HASH
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3638{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid PAY_PER_REQUEST, got %d", len(matches))
	}
}

func TestE3638_InvalidPayPerRequestWithProvisionedThroughput(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyTable:
    Type: AWS::DynamoDB::Table
    Properties:
      TableName: my-table
      BillingMode: PAY_PER_REQUEST
      AttributeDefinitions:
        - AttributeName: id
          AttributeType: S
      KeySchema:
        - AttributeName: id
          KeyType: HASH
      ProvisionedThroughput:
        ReadCapacityUnits: 5
        WriteCapacityUnits: 5
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3638{}
	matches := rule.Match(tmpl)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for PAY_PER_REQUEST with ProvisionedThroughput, got %d", len(matches))
	}
	if len(matches) > 0 && !containsString(matches[0].Message, "ProvisionedThroughput") {
		t.Errorf("Expected error about ProvisionedThroughput, got: %s", matches[0].Message)
	}
}

func TestE3638_ValidProvisioned(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyTable:
    Type: AWS::DynamoDB::Table
    Properties:
      TableName: my-table
      BillingMode: PROVISIONED
      AttributeDefinitions:
        - AttributeName: id
          AttributeType: S
      KeySchema:
        - AttributeName: id
          KeyType: HASH
      ProvisionedThroughput:
        ReadCapacityUnits: 5
        WriteCapacityUnits: 5
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3638{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid PROVISIONED with ProvisionedThroughput, got %d", len(matches))
	}
}

func TestE3638_IntrinsicFunction(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyTable:
    Type: AWS::DynamoDB::Table
    Properties:
      TableName: my-table
      BillingMode:
        Ref: BillingModeParameter
      AttributeDefinitions:
        - AttributeName: id
          AttributeType: S
      KeySchema:
        - AttributeName: id
          KeyType: HASH
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3638{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for intrinsic function, got %d", len(matches))
	}
}

func TestE3638_Metadata(t *testing.T) {
	rule := &E3638{}

	if rule.ID() != "E3638" {
		t.Errorf("Expected ID E3638, got %s", rule.ID())
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
