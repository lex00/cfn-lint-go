package functions

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE1029_ValidSubVariable(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyResource:
    Type: AWS::EC2::Instance
    Properties:
      UserData: !Sub "echo ${AWS::StackName}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1029{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid Sub variable, got %d: %v", len(matches), matches)
	}
}

func TestE1029_InvalidVariableWithoutSub(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyResource:
    Type: AWS::EC2::Instance
    Properties:
      UserData: "echo ${AWS::StackName}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1029{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for variable without Sub, got %d", len(matches))
	}
}

func TestE1029_LiteralDollarSign(t *testing.T) {
	// $$ is escaped in Fn::Sub, but outside Sub context it's just a regular string
	// This test verifies we don't flag regular shell variable syntax
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyResource:
    Type: AWS::EC2::Instance
    Properties:
      UserData: "echo $HOME"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1029{}
	matches := rule.Match(parsed)

	// $HOME is not ${...} syntax, so should not match
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for shell variable syntax, got %d: %v", len(matches), matches)
	}
}

func TestE1029_ValidLiteralNotation(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyResource:
    Type: AWS::EC2::Instance
    Properties:
      UserData: "echo ${!Literal}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1029{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for literal notation, got %d: %v", len(matches), matches)
	}
}
