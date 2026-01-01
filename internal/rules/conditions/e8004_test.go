package conditions

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE8004_ValidAnd(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Conditions:
  BothTrue:
    !And
      - !Equals [a, a]
      - !Equals [b, b]

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E8004{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid Fn::And, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE8004_TooFewConditions(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Conditions:
  BadAnd:
    Fn::And:
      - !Equals [a, a]

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E8004{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for Fn::And with too few conditions")
	}
}

func TestE8004_TooManyConditions(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Conditions:
  BadAnd:
    Fn::And:
      - !Equals [a, a]
      - !Equals [b, b]
      - !Equals [c, c]
      - !Equals [d, d]
      - !Equals [e, e]
      - !Equals [f, f]
      - !Equals [g, g]
      - !Equals [h, h]
      - !Equals [i, i]
      - !Equals [j, j]
      - !Equals [k, k]

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E8004{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for Fn::And with too many conditions")
	}
}

func TestE8004_ExactlyTenConditions(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Conditions:
  MaxAnd:
    Fn::And:
      - !Equals [a, a]
      - !Equals [b, b]
      - !Equals [c, c]
      - !Equals [d, d]
      - !Equals [e, e]
      - !Equals [f, f]
      - !Equals [g, g]
      - !Equals [h, h]
      - !Equals [i, i]
      - !Equals [j, j]

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E8004{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for Fn::And with exactly 10 conditions, got %d", len(matches))
	}
}

func TestE8004_Metadata(t *testing.T) {
	rule := &E8004{}

	if rule.ID() != "E8004" {
		t.Errorf("Expected ID E8004, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
