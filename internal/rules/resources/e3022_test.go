package resources

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3022_Metadata(t *testing.T) {
	rule := &E3022{}

	if rule.ID() != "E3022" {
		t.Errorf("Expected ID E3022, got %s", rule.ID())
	}
	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}
	if len(rule.Tags()) == 0 {
		t.Error("Tags should not be empty")
	}
}

func TestE3022_UniqueSubnets(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  Association1:
    Type: AWS::EC2::SubnetRouteTableAssociation
    Properties:
      SubnetId: subnet-12345
      RouteTableId: rtb-11111
  Association2:
    Type: AWS::EC2::SubnetRouteTableAssociation
    Properties:
      SubnetId: subnet-67890
      RouteTableId: rtb-22222
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3022{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for unique subnets, got %d", len(matches))
	}
}

func TestE3022_DuplicateSubnets(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  Association1:
    Type: AWS::EC2::SubnetRouteTableAssociation
    Properties:
      SubnetId: subnet-12345
      RouteTableId: rtb-11111
  Association2:
    Type: AWS::EC2::SubnetRouteTableAssociation
    Properties:
      SubnetId: subnet-12345
      RouteTableId: rtb-22222
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3022{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for duplicate subnet associations")
	}
}
