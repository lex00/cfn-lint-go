package parameters

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE2014_ConstraintDescriptionWithAllowedPattern(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  Environment:
    Type: String
    AllowedPattern: ^(dev|staging|prod)$
    ConstraintDescription: Must be dev, staging, or prod
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2014{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid ConstraintDescription with AllowedPattern, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE2014_ConstraintDescriptionWithAllowedValues(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  Environment:
    Type: String
    AllowedValues:
      - dev
      - staging
      - prod
    ConstraintDescription: Must be one of the allowed values
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2014{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid ConstraintDescription with AllowedValues, got %d", len(matches))
	}
}

func TestE2014_ConstraintDescriptionWithMinMaxLength(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  Username:
    Type: String
    MinLength: 3
    MaxLength: 20
    ConstraintDescription: Username must be between 3 and 20 characters
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2014{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid ConstraintDescription with MinLength/MaxLength, got %d", len(matches))
	}
}

func TestE2014_ConstraintDescriptionWithMinMaxValue(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  InstanceCount:
    Type: Number
    MinValue: 1
    MaxValue: 10
    ConstraintDescription: Must be between 1 and 10
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2014{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid ConstraintDescription with MinValue/MaxValue, got %d", len(matches))
	}
}

func TestE2014_ConstraintDescriptionWithoutConstraints(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  Environment:
    Type: String
    ConstraintDescription: This description has no matching constraint
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2014{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for ConstraintDescription without constraints")
	}
}

func TestE2014_NoConstraintDescription(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  Environment:
    Type: String
    AllowedValues:
      - dev
      - prod
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2014{}
	matches := rule.Match(tmpl)

	// No ConstraintDescription = no error for E2014
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when no ConstraintDescription, got %d", len(matches))
	}
}

func TestE2014_NoParameters(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2014{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when no parameters, got %d", len(matches))
	}
}

func TestE2014_MultipleParameters(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  GoodParam:
    Type: String
    AllowedPattern: ^[a-z]+$
    ConstraintDescription: Valid constraint description
  BadParam:
    Type: String
    ConstraintDescription: No constraint for this
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2014{}
	matches := rule.Match(tmpl)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for bad parameter, got %d", len(matches))
	}
}

func TestE2014_Metadata(t *testing.T) {
	rule := &E2014{}

	if rule.ID() != "E2014" {
		t.Errorf("Expected ID E2014, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	if rule.Description() == "" {
		t.Error("Description should not be empty")
	}

	if rule.Source() == "" {
		t.Error("Source should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}

	// Check for constraints tag
	hasConstraintsTag := false
	for _, tag := range tags {
		if tag == "constraints" {
			hasConstraintsTag = true
			break
		}
	}
	if !hasConstraintsTag {
		t.Error("Expected 'constraints' tag")
	}
}
