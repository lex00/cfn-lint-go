package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW2030_ValidParameterValue(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  Environment:
    Type: String
    AllowedValues:
      - dev
      - prod
    Default: dev
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W2030{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid parameter default, got %d: %v", len(matches), matches)
	}
}

func TestW2030_ParameterWithMinMax(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  InstanceCount:
    Type: Number
    MinValue: 1
    MaxValue: 10
    Default: 5
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W2030{}
	matches := rule.Match(parsed)

	// W2030 validates parameter value checks, just verify it runs
	_ = matches
}
