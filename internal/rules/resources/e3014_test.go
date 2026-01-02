package resources

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3014_NoConflict(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyInstance:
    Type: AWS::EC2::Instance
    Properties:
      ImageId: ami-12345678
      InstanceType: t2.micro
      SecurityGroupIds:
        - sg-12345678
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3014{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when no conflict, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE3014_MutuallyExclusiveConflict(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyInstance:
    Type: AWS::EC2::Instance
    Properties:
      ImageId: ami-12345678
      InstanceType: t2.micro
      SecurityGroups:
        - default
      SecurityGroupIds:
        - sg-12345678
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3014{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for mutually exclusive properties")
	} else {
		if !strings.Contains(matches[0].Message, "mutually exclusive") {
			t.Errorf("Error message should mention mutually exclusive: %s", matches[0].Message)
		}
		if !strings.Contains(matches[0].Message, "SecurityGroups") {
			t.Errorf("Error message should mention SecurityGroups: %s", matches[0].Message)
		}
	}
}

func TestE3014_SecretsManagerConflict(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MySecret:
    Type: AWS::SecretsManager::Secret
    Properties:
      Name: my-secret
      SecretString: mysecretvalue
      GenerateSecretString:
        PasswordLength: 32
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3014{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for SecretString and GenerateSecretString conflict")
	}
}

func TestE3014_UnknownResourceSkipped(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyCustom:
    Type: AWS::Unknown::Resource
    Properties:
      PropA: valueA
      PropB: valueB
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3014{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for unknown resource type, got %d", len(matches))
	}
}

func TestE3014_Metadata(t *testing.T) {
	rule := &E3014{}

	if rule.ID() != "E3014" {
		t.Errorf("Expected ID E3014, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
