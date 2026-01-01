package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW1011_DynamicReference(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyDB:
    Type: AWS::RDS::DBInstance
    Properties:
      MasterUserPassword: "{{resolve:secretsmanager:MySecret}}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W1011{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for dynamic reference, got %d: %v", len(matches), matches)
	}
}

func TestW1011_HardcodedPassword(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyDB:
    Type: AWS::RDS::DBInstance
    Properties:
      MasterUserPassword: "MySecretPassword123"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W1011{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for hardcoded password, got %d", len(matches))
	}
}

func TestW1011_RefIsOK(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  DBPassword:
    Type: String
    NoEcho: true
Resources:
  MyDB:
    Type: AWS::RDS::DBInstance
    Properties:
      MasterUserPassword: !Ref DBPassword
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W1011{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for Ref to parameter, got %d: %v", len(matches), matches)
	}
}
