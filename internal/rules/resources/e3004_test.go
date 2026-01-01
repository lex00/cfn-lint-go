package resources

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3004_NoCycle(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  ResourceA:
    Type: AWS::S3::Bucket
  ResourceB:
    Type: AWS::S3::Bucket
    DependsOn: ResourceA
  ResourceC:
    Type: AWS::S3::Bucket
    DependsOn: ResourceB
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3004{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for no cycle, got %d: %v", len(matches), matches)
	}
}

func TestE3004_SimpleCycle(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  ResourceA:
    Type: AWS::S3::Bucket
    DependsOn: ResourceB
  ResourceB:
    Type: AWS::S3::Bucket
    DependsOn: ResourceA
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3004{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for simple cycle, got %d", len(matches))
	}
}

func TestE3004_ThreeNodeCycle(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  ResourceA:
    Type: AWS::S3::Bucket
    DependsOn: ResourceC
  ResourceB:
    Type: AWS::S3::Bucket
    DependsOn: ResourceA
  ResourceC:
    Type: AWS::S3::Bucket
    DependsOn: ResourceB
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3004{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for 3-node cycle, got %d", len(matches))
	}
}

func TestE3004_CycleViaRef(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  ResourceA:
    Type: AWS::CloudFormation::WaitConditionHandle
    Properties:
      Metadata: !Ref ResourceB
  ResourceB:
    Type: AWS::CloudFormation::WaitConditionHandle
    Properties:
      Metadata: !Ref ResourceA
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3004{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for cycle via Ref, got %d", len(matches))
	}
}

func TestE3004_CycleViaGetAtt(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  ResourceA:
    Type: AWS::CloudFormation::WaitConditionHandle
    Properties:
      Metadata: !GetAtt ResourceB.Data
  ResourceB:
    Type: AWS::CloudFormation::WaitConditionHandle
    Properties:
      Metadata: !GetAtt ResourceA.Data
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3004{}
	matches := rule.Match(parsed)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for cycle via GetAtt, got %d", len(matches))
	}
}

func TestE3004_RefToParameterNotCycle(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  ParamA:
    Type: String
Resources:
  ResourceA:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !Ref ParamA
  ResourceB:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !Ref ParamA
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3004{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for refs to parameters, got %d: %v", len(matches), matches)
	}
}
