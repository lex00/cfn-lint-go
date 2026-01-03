package resources

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3601_ValidDefinition(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyStateMachine:
    Type: AWS::StepFunctions::StateMachine
    Properties:
      RoleArn: arn:aws:iam::123456789012:role/MyRole
      Definition:
        StartAt: FirstState
        States:
          FirstState:
            Type: Pass
            Result: Done
            End: true
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3601{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid definition, got %d: %v", len(matches), matches)
	}
}

func TestE3601_MissingStates(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyStateMachine:
    Type: AWS::StepFunctions::StateMachine
    Properties:
      RoleArn: arn:aws:iam::123456789012:role/MyRole
      Definition:
        StartAt: FirstState
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3601{}
	matches := rule.Match(tmpl)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for missing States, got %d", len(matches))
	}
	if len(matches) > 0 && !containsString(matches[0].Message, "States") {
		t.Errorf("Expected error about missing States, got: %s", matches[0].Message)
	}
}

func TestE3601_MissingStartAt(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyStateMachine:
    Type: AWS::StepFunctions::StateMachine
    Properties:
      RoleArn: arn:aws:iam::123456789012:role/MyRole
      Definition:
        States:
          FirstState:
            Type: Pass
            Result: Done
            End: true
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3601{}
	matches := rule.Match(tmpl)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for missing StartAt, got %d", len(matches))
	}
	if len(matches) > 0 && !containsString(matches[0].Message, "StartAt") {
		t.Errorf("Expected error about missing StartAt, got: %s", matches[0].Message)
	}
}

func TestE3601_ValidDefinitionString(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyStateMachine:
    Type: AWS::StepFunctions::StateMachine
    Properties:
      RoleArn: arn:aws:iam::123456789012:role/MyRole
      DefinitionString: '{"StartAt":"FirstState","States":{"FirstState":{"Type":"Pass","Result":"Done","End":true}}}'
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3601{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid DefinitionString, got %d: %v", len(matches), matches)
	}
}

func TestE3601_InvalidDefinitionString(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyStateMachine:
    Type: AWS::StepFunctions::StateMachine
    Properties:
      RoleArn: arn:aws:iam::123456789012:role/MyRole
      DefinitionString: 'invalid json'
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3601{}
	matches := rule.Match(tmpl)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for invalid JSON, got %d", len(matches))
	}
}

func TestE3601_IntrinsicFunction(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyStateMachine:
    Type: AWS::StepFunctions::StateMachine
    Properties:
      RoleArn: arn:aws:iam::123456789012:role/MyRole
      Definition:
        Fn::Sub: |
          {
            "StartAt": "FirstState",
            "States": {
              "FirstState": {
                "Type": "Pass",
                "Result": "Done",
                "End": true
              }
            }
          }
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3601{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for intrinsic function, got %d", len(matches))
	}
}

func TestE3601_Metadata(t *testing.T) {
	rule := &E3601{}

	if rule.ID() != "E3601" {
		t.Errorf("Expected ID E3601, got %s", rule.ID())
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

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && findInString(s, substr)
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
