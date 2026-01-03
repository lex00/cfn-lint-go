package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW1035_SelectWithDynamicIndex(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Parameters:
  Index:
    Type: Number
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !Select [!Ref Index, [bucket1, bucket2, bucket3]]
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W1035{}
	matches := rule.Match(parsed)

	// Dynamic index should not trigger warning about static selection
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for dynamic Select, got %d: %v", len(matches), matches)
	}
}

func TestW1035_SelectWithGetAZs(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  MySubnet:
    Type: AWS::EC2::Subnet
    Properties:
      AvailabilityZone: !Select [0, !GetAZs ""]
      CidrBlock: 10.0.0.0/24
      VpcId: !Ref MyVPC
  MyVPC:
    Type: AWS::EC2::VPC
    Properties:
      CidrBlock: 10.0.0.0/16
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W1035{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for Select with GetAZs, got %d: %v", len(matches), matches)
	}
}
