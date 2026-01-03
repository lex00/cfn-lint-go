package parameters

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE2004_SensitiveWithNoEcho(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  DatabasePassword:
    Type: String
    NoEcho: true
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2004{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for sensitive parameter with NoEcho=true, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE2004_SensitiveWithoutNoEcho(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  DatabasePassword:
    Type: String
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2004{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for sensitive parameter without NoEcho")
	}
}

func TestE2004_SensitiveKeywords(t *testing.T) {
	// Test each sensitive keyword
	keywords := []string{
		"Password",
		"Secret",
		"ApiKey",
		"Api_Key",
		"Token",
		"Credential",
		"Passphrase",
		"PrivateKey",
		"Private_Key",
	}

	for _, keyword := range keywords {
		t.Run(keyword, func(t *testing.T) {
			yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  My` + keyword + `:
    Type: String
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
			tmpl, err := template.Parse([]byte(yaml))
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			rule := &E2004{}
			matches := rule.Match(tmpl)

			if len(matches) == 0 {
				t.Errorf("Expected match for parameter containing '%s' without NoEcho", keyword)
			}
		})
	}
}

func TestE2004_NonSensitiveParameter(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  Environment:
    Type: String
    Default: dev
  InstanceType:
    Type: String
    Default: t2.micro
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2004{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for non-sensitive parameters, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE2004_CaseInsensitive(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  DATABASE_PASSWORD:
    Type: String
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2004{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for uppercase sensitive parameter name")
	}
}

func TestE2004_NoParameters(t *testing.T) {
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

	rule := &E2004{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when no parameters, got %d", len(matches))
	}
}

func TestE2004_NoEchoFalse(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  SecretKey:
    Type: String
    NoEcho: false
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2004{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for sensitive parameter with NoEcho=false")
	}
}

func TestE2004_Metadata(t *testing.T) {
	rule := &E2004{}

	if rule.ID() != "E2004" {
		t.Errorf("Expected ID E2004, got %s", rule.ID())
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

	// Check for security tag
	hasSecurityTag := false
	for _, tag := range tags {
		if tag == "security" {
			hasSecurityTag = true
			break
		}
	}
	if !hasSecurityTag {
		t.Error("Expected 'security' tag")
	}
}
