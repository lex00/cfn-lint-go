package functions

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE1021_ValidBase64(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyResource:
    Type: AWS::EC2::Instance
    Properties:
      UserData: !Base64 "#!/bin/bash\necho hello"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1021{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid Base64, got %d: %v", len(matches), matches)
	}
}

func TestE1021_ValidBase64WithSub(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyResource:
    Type: AWS::EC2::Instance
    Properties:
      UserData:
        Fn::Base64:
          Fn::Sub: "#!/bin/bash\necho ${AWS::StackName}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1021{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid Base64 with Sub, got %d: %v", len(matches), matches)
	}
}

func TestE1021_InvalidBase64(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyResource:
    Type: AWS::EC2::Instance
    Properties:
      UserData:
        Fn::Base64:
          - item1
          - item2
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1021{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for invalid Base64, got %d", len(matches))
	}
}
