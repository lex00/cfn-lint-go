package resources

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3019_Metadata(t *testing.T) {
	rule := &E3019{}

	if rule.ID() != "E3019" {
		t.Errorf("Expected ID E3019, got %s", rule.ID())
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
	if len(rule.Tags()) == 0 {
		t.Error("Tags should not be empty")
	}
}

func TestE3019_UniqueIdentifiers(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  Bucket1:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: my-bucket-1
  Bucket2:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: my-bucket-2
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3019{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for unique identifiers, got %d", len(matches))
	}
}

func TestE3019_DuplicateIdentifiers(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  Bucket1:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: duplicate-bucket
  Bucket2:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: duplicate-bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3019{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for duplicate identifiers")
	}
}

func TestE3019_NoIdentifiers(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  Bucket1:
    Type: AWS::S3::Bucket
  Bucket2:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3019{}
	matches := rule.Match(tmpl)

	// No explicit identifiers, CloudFormation will auto-generate
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when no identifiers specified, got %d", len(matches))
	}
}
