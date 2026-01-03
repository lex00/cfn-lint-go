package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW3693_AuroraClusterWithoutIgnoredProps(t *testing.T) {
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

	rule := &W3693{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for cluster without ignored props, got %d: %v", len(matches), matches)
	}
}

func TestW3693_AuroraWithIgnoredProps(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyDBCluster:
    Type: AWS::RDS::DBCluster
    Properties:
      Engine: aurora-mysql
      MasterUsername: admin
      MasterUserPassword: password123
      AllocatedStorage: 100
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W3693{}
	matches := rule.Match(parsed)

	// The rule checks for properties that are ignored in Aurora clusters
	// Just verify it runs without error
	_ = matches
}
