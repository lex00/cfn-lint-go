package parameters

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE2012_ValidBasicTypes(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  StringParam:
    Type: String
  NumberParam:
    Type: Number
  ListParam:
    Type: CommaDelimitedList
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2012{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid basic types, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE2012_ValidAWSTypes(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  AZ:
    Type: AWS::EC2::AvailabilityZone::Name
  AMI:
    Type: AWS::EC2::Image::Id
  Instance:
    Type: AWS::EC2::Instance::Id
  KeyPair:
    Type: AWS::EC2::KeyPair::KeyName
  SG:
    Type: AWS::EC2::SecurityGroup::Id
  Subnet:
    Type: AWS::EC2::Subnet::Id
  VPC:
    Type: AWS::EC2::VPC::Id
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2012{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid AWS types, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE2012_ValidListTypes(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  StringList:
    Type: List<String>
  NumberList:
    Type: List<Number>
  AZList:
    Type: List<AWS::EC2::AvailabilityZone::Name>
  SGList:
    Type: List<AWS::EC2::SecurityGroup::Id>
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2012{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid list types, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE2012_ValidSSMTypes(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  SSMName:
    Type: AWS::SSM::Parameter::Name
  SSMString:
    Type: AWS::SSM::Parameter::Value<String>
  SSMList:
    Type: AWS::SSM::Parameter::Value<List<String>>
  SSMAMI:
    Type: AWS::SSM::Parameter::Value<AWS::EC2::Image::Id>
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2012{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid SSM types, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE2012_InvalidType(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  InvalidParam:
    Type: InvalidType
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2012{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for invalid parameter type")
	}
}

func TestE2012_InvalidAWSType(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  InvalidParam:
    Type: AWS::Invalid::Type
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2012{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for invalid AWS parameter type")
	}
}

func TestE2012_InvalidListType(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  InvalidParam:
    Type: List<InvalidType>
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2012{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for invalid list parameter type")
	}
}

func TestE2012_EmptyType(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  NoTypeParam:
    Default: value
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E2012{}
	matches := rule.Match(tmpl)

	// Empty type should not trigger E2012 (that's E2001's job)
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for empty type (E2012 only validates non-empty types), got %d", len(matches))
	}
}

func TestE2012_NoParameters(t *testing.T) {
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

	rule := &E2012{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when no parameters, got %d", len(matches))
	}
}

func TestE2012_Metadata(t *testing.T) {
	rule := &E2012{}

	if rule.ID() != "E2012" {
		t.Errorf("Expected ID E2012, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	if rule.Description() == "" {
		t.Error("Description should not be empty")
	}

	if rule.Source() == "" {
		t.Error("Source should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}

	// Check for ssm tag
	hasSSMTag := false
	for _, tag := range tags {
		if tag == "ssm" {
			hasSSMTag = true
			break
		}
	}
	if !hasSSMTag {
		t.Error("Expected 'ssm' tag")
	}
}
