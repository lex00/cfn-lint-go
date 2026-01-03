package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW3690_CurrentEngineVersion(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyDBCluster:
    Type: AWS::RDS::DBCluster
    Properties:
      Engine: aurora-mysql
      EngineVersion: "8.0.mysql_aurora.3.04.0"
      MasterUsername: admin
      MasterUserPassword: password123
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W3690{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for current engine version, got %d: %v", len(matches), matches)
	}
}

func TestW3690_DeprecatedEngineVersion(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyDBCluster:
    Type: AWS::RDS::DBCluster
    Properties:
      Engine: aurora-postgresql
      EngineVersion: "9.6"
      MasterUsername: admin
      MasterUserPassword: password123
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W3690{}
	matches := rule.Match(parsed)

	if len(matches) == 0 {
		t.Errorf("Expected matches for deprecated engine version, got 0")
	}
}
