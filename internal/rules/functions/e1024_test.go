package functions

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE1024_ValidCidr(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MySubnet:
    Type: AWS::EC2::Subnet
    Properties:
      CidrBlock: !Select [0, !Cidr ["10.0.0.0/16", 256, 8]]
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1024{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid Cidr, got %d: %v", len(matches), matches)
	}
}

func TestE1024_InvalidCidrNotArray(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MySubnet:
    Type: AWS::EC2::Subnet
    Properties:
      CidrBlock:
        Fn::Cidr: "invalid"
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1024{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for invalid Cidr, got %d", len(matches))
	}
}

func TestE1024_InvalidCidrWrongCount(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MySubnet:
    Type: AWS::EC2::Subnet
    Properties:
      CidrBlock:
        Fn::Cidr:
          - "10.0.0.0/16"
          - 256
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E1024{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for invalid Cidr (wrong count), got %d", len(matches))
	}
}
