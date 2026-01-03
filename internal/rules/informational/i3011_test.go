package informational

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestI3011_StatefulResourceWithoutPolicies(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyDatabase:
    Type: AWS::RDS::DBInstance
    Properties:
      Engine: postgres
      DBInstanceClass: db.t3.micro
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I3011{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for stateful resource without deletion policies")
	}
}

func TestI3011_StatefulResourceWithPolicies(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyDatabase:
    Type: AWS::RDS::DBInstance
    DeletionPolicy: Snapshot
    UpdateReplacePolicy: Snapshot
    Properties:
      Engine: postgres
      DBInstanceClass: db.t3.micro
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I3011{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for stateful resource with both policies, got %d", len(matches))
	}
}

func TestI3011_StatelessResource(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Runtime: python3.11
      Handler: index.handler
      Code:
        ZipFile: |
          def handler(event, context):
            return 'Hello'
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I3011{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for stateless resource, got %d", len(matches))
	}
}

func TestI3011_S3BucketWithoutPolicies(t *testing.T) {
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

	rule := &I3011{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for S3 bucket without deletion policies")
	}
}

func TestI3011_Metadata(t *testing.T) {
	rule := &I3011{}

	if rule.ID() != "I3011" {
		t.Errorf("Expected ID I3011, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
