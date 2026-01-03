package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW2501_PasswordWithFullValidation(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  DBPassword:
    Type: String
    NoEcho: true
    MinLength: 8
    AllowedPattern: "^[a-zA-Z0-9]+$"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W2501{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for password with full validation, got %d: %v", len(matches), matches)
	}
}

func TestW2501_PasswordWithoutNoEcho(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  DBPassword:
    Type: String
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W2501{}
	matches := rule.Match(parsed)

	if len(matches) == 0 {
		t.Errorf("Expected matches for password without NoEcho, got 0")
	}
}
