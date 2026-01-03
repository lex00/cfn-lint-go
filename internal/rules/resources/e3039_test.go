package resources

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3039_Metadata(t *testing.T) {
	rule := &E3039{}

	if rule.ID() != "E3039" {
		t.Errorf("Expected ID E3039, got %s", rule.ID())
	}
	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}
	if len(rule.Tags()) == 0 {
		t.Error("Tags should not be empty")
	}
}

func TestE3039_ValidTable(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyTable:
    Type: AWS::DynamoDB::Table
    Properties:
      TableName: my-table
      AttributeDefinitions:
        - AttributeName: id
          AttributeType: S
        - AttributeName: timestamp
          AttributeType: N
      KeySchema:
        - AttributeName: id
          KeyType: HASH
        - AttributeName: timestamp
          KeyType: RANGE
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3039{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid table, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  %s", m.Message)
		}
	}
}

func TestE3039_MissingAttributeDefinition(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyTable:
    Type: AWS::DynamoDB::Table
    Properties:
      TableName: my-table
      AttributeDefinitions:
        - AttributeName: id
          AttributeType: S
      KeySchema:
        - AttributeName: id
          KeyType: HASH
        - AttributeName: timestamp
          KeyType: RANGE
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3039{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for missing attribute definition")
	}
}

func TestE3039_UnusedAttributeDefinition(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyTable:
    Type: AWS::DynamoDB::Table
    Properties:
      TableName: my-table
      AttributeDefinitions:
        - AttributeName: id
          AttributeType: S
        - AttributeName: unused
          AttributeType: S
      KeySchema:
        - AttributeName: id
          KeyType: HASH
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3039{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for unused attribute definition")
	}
}

func TestE3039_WithGSI(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyTable:
    Type: AWS::DynamoDB::Table
    Properties:
      TableName: my-table
      AttributeDefinitions:
        - AttributeName: id
          AttributeType: S
        - AttributeName: gsi_key
          AttributeType: S
      KeySchema:
        - AttributeName: id
          KeyType: HASH
      GlobalSecondaryIndexes:
        - IndexName: MyGSI
          KeySchema:
            - AttributeName: gsi_key
              KeyType: HASH
          Projection:
            ProjectionType: ALL
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3039{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid table with GSI, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  %s", m.Message)
		}
	}
}
