package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW4001_ValidInterface(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Metadata:
  AWS::CloudFormation::Interface:
    ParameterGroups:
      - Label:
          default: "Network Configuration"
        Parameters:
          - VPCId
          - SubnetId
    ParameterLabels:
      VPCId:
        default: "VPC ID"
Parameters:
  VPCId:
    Type: String
  SubnetId:
    Type: String
Resources:
  MyResource:
    Type: AWS::EC2::Instance
    Properties:
      SubnetId: !Ref SubnetId
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W4001{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid interface, got %d: %v", len(matches), matches)
	}
}

func TestW4001_UndefinedParameterInGroup(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Metadata:
  AWS::CloudFormation::Interface:
    ParameterGroups:
      - Label:
          default: "Network Configuration"
        Parameters:
          - VPCId
          - NonExistentParam
Parameters:
  VPCId:
    Type: String
Resources:
  MyResource:
    Type: AWS::EC2::Instance
    Properties:
      SubnetId: subnet-123
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W4001{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for undefined parameter in group, got %d", len(matches))
	}
}

func TestW4001_UndefinedParameterInLabels(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Metadata:
  AWS::CloudFormation::Interface:
    ParameterLabels:
      NonExistentParam:
        default: "Some Label"
Parameters:
  VPCId:
    Type: String
Resources:
  MyResource:
    Type: AWS::EC2::Instance
    Properties:
      SubnetId: subnet-123
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W4001{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for undefined parameter in labels, got %d", len(matches))
	}
}
