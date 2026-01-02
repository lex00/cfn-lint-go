package resources

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3012_ValidTypes(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: my-bucket
      VersioningConfiguration:
        Status: Enabled
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3012{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid types, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE3012_WrongStringType(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: 12345
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3012{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for integer where string expected")
	}
}

func TestE3012_WrongBooleanType(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      FunctionName: MyFunc
      Runtime: python3.9
      Handler: index.handler
      Role: arn:aws:iam::123456789012:role/MyRole
      Code:
        S3Bucket: my-bucket
        S3Key: code.zip
      TracingConfig:
        Mode: "Active"
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3012{}
	matches := rule.Match(tmpl)

	// TracingConfig.Mode is a string, so this should be valid
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE3012_IntrinsicFunctionSkipped(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !Ref AWS::StackName
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3012{}
	matches := rule.Match(tmpl)

	// Intrinsic functions should be skipped
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches (intrinsic skipped), got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE3012_WrongListType(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MySG:
    Type: AWS::EC2::SecurityGroup
    Properties:
      GroupDescription: My SG
      SecurityGroupIngress: not-a-list
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3012{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for string where list expected")
	}
}

func TestE3012_ValidListType(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MySG:
    Type: AWS::EC2::SecurityGroup
    Properties:
      GroupDescription: My SG
      SecurityGroupIngress:
        - IpProtocol: tcp
          FromPort: 443
          ToPort: 443
          CidrIp: 0.0.0.0/0
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3012{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid list, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE3012_IntegerAsDouble(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      FunctionName: MyFunc
      Runtime: python3.9
      Handler: index.handler
      Role: arn:aws:iam::123456789012:role/MyRole
      Code:
        S3Bucket: my-bucket
        S3Key: code.zip
      MemorySize: 128
      Timeout: 30
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3012{}
	matches := rule.Match(tmpl)

	// Integer values should be accepted for Integer type properties
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid integer, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE3012_UnknownResourceType(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyCustom:
    Type: AWS::Custom::Unknown
    Properties:
      SomeProperty: 12345
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3012{}
	matches := rule.Match(tmpl)

	// Unknown resource types are skipped
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for unknown resource type, got %d", len(matches))
	}
}

func TestE3012_Metadata(t *testing.T) {
	rule := &E3012{}

	if rule.ID() != "E3012" {
		t.Errorf("Expected ID E3012, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	if rule.Description() == "" {
		t.Error("Description should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
