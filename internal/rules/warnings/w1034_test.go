package warnings

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestW1034_ValidFindInMap(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Mappings:
  RegionMap:
    us-east-1:
      AMI: ami-12345678
    us-west-2:
      AMI: ami-87654321
Resources:
  MyInstance:
    Type: AWS::EC2::Instance
    Properties:
      ImageId: !FindInMap [RegionMap, !Ref "AWS::Region", AMI]
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W1034{}
	matches := rule.Match(parsed)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid FindInMap, got %d: %v", len(matches), matches)
	}
}

func TestW1034_FindInMapWithDefaultValue(t *testing.T) {
	tmpl := `
AWSTemplateFormatVersion: "2010-09-09"
Mappings:
  RegionMap:
    us-east-1:
      AMI: ami-12345678
Resources:
  MyInstance:
    Type: AWS::EC2::Instance
    Properties:
      ImageId: !FindInMap
        - RegionMap
        - !Ref "AWS::Region"
        - AMI
        - DefaultValue: ami-default
`
	parsed, err := template.Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &W1034{}
	matches := rule.Match(parsed)

	// Should not error on FindInMap with default value
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for FindInMap with default, got %d: %v", len(matches), matches)
	}
}
