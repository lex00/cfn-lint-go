package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW3691_CurrentEngineVersion(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyDBInstance:
    Type: AWS::RDS::DBInstance
    Properties:
      DBInstanceClass: db.t3.micro
      Engine: mysql
      EngineVersion: "8.0.35"
      MasterUsername: admin
      MasterUserPassword: password123
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W3691{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for current engine version, got %d: %v", len(matches), matches)
	}
}

func TestW3691_DeprecatedEngineVersion(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyDBInstance:
    Type: AWS::RDS::DBInstance
    Properties:
      DBInstanceClass: db.t3.micro
      Engine: mysql
      EngineVersion: "5.6"
      MasterUsername: admin
      MasterUserPassword: password123
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W3691{}
	matches := rule.Match(parsed)

	if len(matches) == 0 {
		t.Errorf("Expected matches for deprecated engine version, got 0")
	}
}
