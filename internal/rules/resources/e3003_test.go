package resources

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3003_ValidLambda(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Role: arn:aws:iam::123456789012:role/MyRole
      Code:
        S3Bucket: my-bucket
        S3Key: code.zip
      Runtime: python3.9
      Handler: index.handler
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3003{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid Lambda, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE3003_MissingLambdaRole(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Code:
        S3Bucket: my-bucket
        S3Key: code.zip
      Runtime: python3.9
      Handler: index.handler
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3003{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for Lambda missing Role")
	}
}

func TestE3003_MissingIAMRolePolicy(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyRole:
    Type: AWS::IAM::Role
    Properties:
      RoleName: MyRole
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3003{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for IAM Role missing AssumeRolePolicyDocument")
	}
}

func TestE3003_ValidSecurityGroup(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MySG:
    Type: AWS::EC2::SecurityGroup
    Properties:
      GroupDescription: My security group
      VpcId: vpc-12345678
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3003{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid SecurityGroup, got %d", len(matches))
	}
}

func TestE3003_MissingSecurityGroupDescription(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MySG:
    Type: AWS::EC2::SecurityGroup
    Properties:
      VpcId: vpc-12345678
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3003{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for SecurityGroup missing GroupDescription")
	}
}

func TestE3003_S3BucketNoRequired(t *testing.T) {
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

	rule := &E3003{}
	matches := rule.Match(tmpl)

	// S3::Bucket has no required properties
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for S3 Bucket (no required props), got %d", len(matches))
	}
}

func TestE3003_UnknownResourceType(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyCustom:
    Type: AWS::Custom::Resource
    Properties:
      SomeProperty: value
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3003{}
	matches := rule.Match(tmpl)

	// Unknown resource types are skipped
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for unknown resource type, got %d", len(matches))
	}
}

func TestE3003_Metadata(t *testing.T) {
	rule := &E3003{}

	if rule.ID() != "E3003" {
		t.Errorf("Expected ID E3003, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
