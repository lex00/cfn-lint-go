package functions

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE1022_ValidJoin(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyResource:
    Type: AWS::EC2::Instance
    Properties:
      Tags:
        - Key: Name
          Value: !Join ["-", ["prefix", "suffix"]]
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1022{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid Join, got %d: %v", len(matches), matches)
	}
}

func TestE1022_InvalidJoinNotArray(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyResource:
    Type: AWS::EC2::Instance
    Properties:
      Tags:
        - Key: Name
          Value:
            Fn::Join: "invalid"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1022{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for invalid Join, got %d", len(matches))
	}
}

func TestE1022_InvalidJoinWrongCount(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyResource:
    Type: AWS::EC2::Instance
    Properties:
      Tags:
        - Key: Name
          Value:
            Fn::Join:
              - "-"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1022{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for invalid Join (wrong count), got %d", len(matches))
	}
}

func TestE1022_InvalidJoinDelimiterNotString(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyResource:
    Type: AWS::EC2::Instance
    Properties:
      Tags:
        - Key: Name
          Value:
            Fn::Join:
              - 123
              - ["a", "b"]
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1022{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for invalid Join (delimiter not string), got %d", len(matches))
	}
}
