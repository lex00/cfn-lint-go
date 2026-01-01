package resources

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3005_ValidDependsOn(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  ResourceA:
    Type: AWS::S3::Bucket
  ResourceB:
    Type: AWS::S3::Bucket
    DependsOn: ResourceA
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3005{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid DependsOn, got %d: %v", len(matches), matches)
	}
}

func TestE3005_UndefinedDependsOn(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  ResourceA:
    Type: AWS::S3::Bucket
    DependsOn: NonExistentResource
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3005{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for undefined DependsOn, got %d", len(matches))
	}
}

func TestE3005_MultipleDependsOn(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  ResourceA:
    Type: AWS::S3::Bucket
  ResourceB:
    Type: AWS::S3::Bucket
  ResourceC:
    Type: AWS::S3::Bucket
    DependsOn:
      - ResourceA
      - NonExistentResource
      - ResourceB
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3005{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for one undefined in list, got %d", len(matches))
	}
}

func TestE3005_SelfDependency(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  ResourceA:
    Type: AWS::S3::Bucket
    DependsOn: ResourceA
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3005{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for self-dependency, got %d", len(matches))
	}
}

func TestE3005_MultipleUndefined(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  ResourceA:
    Type: AWS::S3::Bucket
    DependsOn:
      - NonExistent1
      - NonExistent2
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3005{}
	matches := rule.Match(parsed)

	if len(matches) != 2 {
		t.Errorf("Expected 2 matches for multiple undefined, got %d", len(matches))
	}
}
