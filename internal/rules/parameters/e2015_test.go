package parameters

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE2015_ValidDefault(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  Environment:
    Type: String
    Default: prod
    AllowedValues:
      - dev
      - staging
      - prod
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2015{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid default, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE2015_InvalidAllowedValues(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  Environment:
    Type: String
    Default: invalid
    AllowedValues:
      - dev
      - staging
      - prod
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2015{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for default not in AllowedValues")
	}
}

func TestE2015_InvalidPattern(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  BucketName:
    Type: String
    Default: INVALID_BUCKET
    AllowedPattern: "^[a-z0-9-]+$"
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2015{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for default not matching AllowedPattern")
	}
}

func TestE2015_NumberBelowMin(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  InstanceCount:
    Type: Number
    Default: 0
    MinValue: 1
    MaxValue: 10
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2015{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for default below MinValue")
	}
}

func TestE2015_NumberAboveMax(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  InstanceCount:
    Type: Number
    Default: 100
    MinValue: 1
    MaxValue: 10
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2015{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for default above MaxValue")
	}
}

func TestE2015_StringTooShort(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  Password:
    Type: String
    Default: "ab"
    MinLength: 8
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2015{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for default shorter than MinLength")
	}
}

func TestE2015_StringTooLong(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  ShortCode:
    Type: String
    Default: "abcdefghij"
    MaxLength: 5
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2015{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for default longer than MaxLength")
	}
}

func TestE2015_NoDefault(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  Environment:
    Type: String
    AllowedValues:
      - dev
      - prod
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2015{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when no default, got %d", len(matches))
	}
}

func TestE2015_Metadata(t *testing.T) {
	rule := &E2015{}

	if rule.ID() != "E2015" {
		t.Errorf("Expected ID E2015, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
