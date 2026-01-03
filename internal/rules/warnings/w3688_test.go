package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW3688_NewDBCluster(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyDBCluster:
    Type: AWS::RDS::DBCluster
    Properties:
      Engine: aurora-mysql
      MasterUsername: admin
      MasterUserPassword: password123
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W3688{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for new DB cluster, got %d: %v", len(matches), matches)
	}
}

func TestW3688_RestoreWithIgnoredProps(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyDBCluster:
    Type: AWS::RDS::DBCluster
    Properties:
      Engine: aurora-mysql
      SnapshotIdentifier: my-snapshot
      MasterUsername: admin
      MasterUserPassword: password123
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W3688{}
	matches := rule.Match(parsed)

	if len(matches) == 0 {
		t.Errorf("Expected matches for restore with ignored props, got 0")
	}
}
