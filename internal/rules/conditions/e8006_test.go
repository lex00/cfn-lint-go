package conditions

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE8006_ValidOr(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Conditions:
  EitherTrue:
    !Or
      - !Equals [a, a]
      - !Equals [b, c]

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E8006{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid Fn::Or, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE8006_TooFewConditions(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Conditions:
  BadOr:
    Fn::Or:
      - !Equals [a, a]

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E8006{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for Fn::Or with too few conditions")
	}
}

func TestE8006_TooManyConditions(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Conditions:
  BadOr:
    Fn::Or:
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

	rule := &E8006{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for Fn::Or with too many conditions")
	}
}

func TestE8006_ExactlyTenConditions(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Conditions:
  MaxOr:
    Fn::Or:
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

	rule := &E8006{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for Fn::Or with exactly 10 conditions, got %d", len(matches))
	}
}

func TestE8006_Metadata(t *testing.T) {
	rule := &E8006{}

	if rule.ID() != "E8006" {
		t.Errorf("Expected ID E8006, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
