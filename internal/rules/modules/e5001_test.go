package modules

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE5001_Metadata(t *testing.T) {
	rule := &E5001{}

	if rule.ID() != "E5001" {
		t.Errorf("Expected ID 'E5001', got '%s'", rule.ID())
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
}

func TestE5001_ValidModule(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyModule:
    Type: MyOrganization::MyService::MyResource::MODULE
    Properties:
      Param1: Value1
      Param2: Value2
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E5001{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid module, got %d: %v", len(matches), matches)
	}
}

func TestE5001_InvalidModuleFormat(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  InvalidModule:
    Type: MyOrganization::MODULE
    Properties:
      Param1: Value1
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E5001{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for invalid module format, got %d", len(matches))
	}

	if len(matches) > 0 && !contains(matches[0].Message, "invalid type format") {
		t.Errorf("Expected error about invalid type format, got: %s", matches[0].Message)
	}
}

func TestE5001_ModuleMissingProperties(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  ModuleNoProps:
    Type: MyOrganization::MyService::MyResource::MODULE
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E5001{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for module missing properties, got %d", len(matches))
	}

	if len(matches) > 0 && !contains(matches[0].Message, "must have Properties") {
		t.Errorf("Expected error about missing Properties, got: %s", matches[0].Message)
	}
}

func TestE5001_NonModuleResource(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E5001{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for non-module resource, got %d", len(matches))
	}
}

func TestE5001_ModuleWithEmptyComponent(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  BadModule:
    Type: MyOrganization::::MODULE
    Properties:
      Param1: Value1
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E5001{}
	matches := rule.Match(parsed)

	if len(matches) == 0 {
		t.Error("Expected at least 1 match for module with empty component")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
