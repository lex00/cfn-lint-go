package resources

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3030_ValidLambdaRuntime(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      FunctionName: MyFunc
      Runtime: python3.12
      Handler: index.handler
      Role: arn:aws:iam::123456789012:role/MyRole
      Code:
        S3Bucket: my-bucket
        S3Key: code.zip
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3030{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid runtime, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE3030_InvalidLambdaRuntime(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      FunctionName: MyFunc
      Runtime: python99
      Handler: index.handler
      Role: arn:aws:iam::123456789012:role/MyRole
      Code:
        S3Bucket: my-bucket
        S3Key: code.zip
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3030{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for invalid runtime 'python99'")
	} else {
		// Verify error message mentions the invalid value
		if !strings.Contains(matches[0].Message, "python99") {
			t.Errorf("Error message should contain 'python99': %s", matches[0].Message)
		}
	}
}

func TestE3030_ValidEC2VolumeType(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyVolume:
    Type: AWS::EC2::Volume
    Properties:
      AvailabilityZone: us-east-1a
      Size: 100
      VolumeType: gp3
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3030{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid VolumeType, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE3030_InvalidEC2VolumeType(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyVolume:
    Type: AWS::EC2::Volume
    Properties:
      AvailabilityZone: us-east-1a
      Size: 100
      VolumeType: super-fast-ssd
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3030{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for invalid VolumeType 'super-fast-ssd'")
	}
}

func TestE3030_ValidECSLaunchType(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyService:
    Type: AWS::ECS::Service
    Properties:
      ServiceName: my-service
      LaunchType: FARGATE
      TaskDefinition: arn:aws:ecs:us-east-1:123456789012:task-definition/my-task:1
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3030{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid LaunchType, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE3030_InvalidECSLaunchType(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyService:
    Type: AWS::ECS::Service
    Properties:
      ServiceName: my-service
      LaunchType: KUBERNETES
      TaskDefinition: arn:aws:ecs:us-east-1:123456789012:task-definition/my-task:1
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3030{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for invalid LaunchType 'KUBERNETES'")
	}
}

func TestE3030_IntrinsicFunctionSkipped(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  Runtime:
    Type: String
    Default: python3.12
Resources:
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      FunctionName: MyFunc
      Runtime: !Ref Runtime
      Handler: index.handler
      Role: arn:aws:iam::123456789012:role/MyRole
      Code:
        S3Bucket: my-bucket
        S3Key: code.zip
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3030{}
	matches := rule.Match(tmpl)

	// Intrinsic functions should be skipped (not string values)
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches (intrinsic skipped), got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE3030_NonEnumPropertySkipped(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: any-bucket-name-is-valid
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3030{}
	matches := rule.Match(tmpl)

	// BucketName is not an enum property
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for non-enum property, got %d", len(matches))
	}
}

func TestE3030_CustomResourceSkipped(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyCustom:
    Type: Custom::MyResource
    Properties:
      SomeEnumLikeProperty: invalid-value
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3030{}
	matches := rule.Match(tmpl)

	// Custom resources are not validated
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for custom resource, got %d", len(matches))
	}
}

func TestE3030_ValidS3StorageClass(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: my-bucket
      LifecycleConfiguration:
        Rules:
          - Id: MoveToGlacier
            Status: Enabled
            Transitions:
              - StorageClass: GLACIER
                TransitionInDays: 90
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3030{}
	matches := rule.Match(tmpl)

	// Note: StorageClass is nested, so top-level property check won't catch it
	// This test verifies no false positives on complex structures
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE3030_Metadata(t *testing.T) {
	rule := &E3030{}

	if rule.ID() != "E3030" {
		t.Errorf("Expected ID E3030, got %s", rule.ID())
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

func TestExtractServiceName(t *testing.T) {
	tests := []struct {
		resourceType string
		expected     string
	}{
		{"AWS::Lambda::Function", "lambda"},
		{"AWS::S3::Bucket", "s3"},
		{"AWS::EC2::Instance", "ec2"},
		{"AWS::ECS::Service", "ecs"},
		{"Custom::MyResource", ""},
		{"Invalid", ""},
		{"", ""},
	}

	for _, tt := range tests {
		result := extractServiceName(tt.resourceType)
		if result != tt.expected {
			t.Errorf("extractServiceName(%q) = %q, expected %q", tt.resourceType, result, tt.expected)
		}
	}
}
