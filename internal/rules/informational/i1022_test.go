package informational

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestI1022_SimpleJoinWithEmptyDelimiter(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !Join
        - ''
        - - 'my-bucket-'
          - !Ref AWS::AccountId
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I1022{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match suggesting Fn::Sub for simple join with empty delimiter")
	}
}

func TestI1022_JoinWithDelimiter(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !Join
        - '-'
        - - 'my'
          - 'bucket'
          - !Ref AWS::AccountId
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I1022{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match suggesting Fn::Sub for simple join with delimiter")
	}
}

func TestI1022_ComplexJoin(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !Join
        - ''
        - - 'my-bucket-'
          - !Join
            - '-'
            - - 'nested'
              - 'join'
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I1022{}
	matches := rule.Match(tmpl)

	// The outer join has a nested join, so it should not trigger (too complex)
	// The inner join is simple and should trigger a suggestion
	if len(matches) != 1 {
		t.Errorf("Expected 1 match for the simple inner join, got %d", len(matches))
	}
}

func TestI1022_SubFunction(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !Sub 'my-bucket-${AWS::AccountId}'
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &I1022{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for Fn::Sub usage, got %d", len(matches))
	}
}

func TestI1022_Metadata(t *testing.T) {
	rule := &I1022{}

	if rule.ID() != "I1022" {
		t.Errorf("Expected ID I1022, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
