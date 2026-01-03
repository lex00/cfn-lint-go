package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW1051_ValidSecretsManagerRef(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyDB:
    Type: AWS::RDS::DBInstance
    Properties:
      DBInstanceClass: db.t3.micro
      Engine: mysql
      MasterUsername: admin
      MasterUserPassword: "{{resolve:secretsmanager:MySecret:SecretString:password}}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W1051{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid Secrets Manager ref, got %d: %v", len(matches), matches)
	}
}

func TestW1051_SecretsManagerARNInRef(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyDB:
    Type: AWS::RDS::DBInstance
    Properties:
      DBInstanceClass: db.t3.micro
      Engine: mysql
      MasterUsername: admin
      MasterUserPassword: "{{resolve:secretsmanager:arn:aws:secretsmanager:us-east-1:123456789012:secret:MySecret:SecretString:password}}"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W1051{}
	matches := rule.Match(parsed)

	// Should warn about using full ARN instead of secret name
	if len(matches) == 0 {
		t.Errorf("Expected matches for Secrets Manager ARN in dynamic ref, got 0")
	}
}
