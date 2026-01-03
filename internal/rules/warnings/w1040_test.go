package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW1040_ValidToJsonString(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MyPolicy:
    Type: AWS::IAM::Policy
    Properties:
      PolicyName: MyPolicy
      PolicyDocument:
        Fn::ToJsonString:
          Version: "2012-10-17"
          Statement:
            - Effect: Allow
              Action: s3:GetObject
              Resource: "*"
      Roles:
        - !Ref MyRole
  MyRole:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument:
        Version: "2012-10-17"
        Statement:
          - Effect: Allow
            Principal:
              Service: lambda.amazonaws.com
            Action: sts:AssumeRole
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W1040{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid ToJsonString, got %d: %v", len(matches), matches)
	}
}
