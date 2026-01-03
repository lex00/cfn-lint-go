package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW1033_ValidSplit(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      Tags:
        - Key: Parts
          Value: !Select [0, !Split [",", "a,b,c"]]
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W1033{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid Split, got %d: %v", len(matches), matches)
	}
}

func TestW1033_SplitWithRef(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  CsvList:
    Type: String
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      Tags:
        - Key: Parts
          Value: !Select [0, !Split [",", !Ref CsvList]]
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W1033{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for Split with Ref, got %d: %v", len(matches), matches)
	}
}
