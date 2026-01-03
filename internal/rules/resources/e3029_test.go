package resources

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3029_Metadata(t *testing.T) {
	rule := &E3029{}

	if rule.ID() != "E3029" {
		t.Errorf("Expected ID E3029, got %s", rule.ID())
	}
	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}
	if len(rule.Tags()) == 0 {
		t.Error("Tags should not be empty")
	}
}

func TestE3029_ValidAlias(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyRecordSet:
    Type: AWS::Route53::RecordSet
    Properties:
      HostedZoneId: Z1234567890ABC
      Name: example.com
      Type: A
      AliasTarget:
        DNSName: example.cloudfront.net
        HostedZoneId: Z2FDTNDATAQYW2
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3029{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid alias, got %d", len(matches))
	}
}

func TestE3029_AliasWithTTL(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyRecordSet:
    Type: AWS::Route53::RecordSet
    Properties:
      HostedZoneId: Z1234567890ABC
      Name: example.com
      Type: A
      TTL: 300
      AliasTarget:
        DNSName: example.cloudfront.net
        HostedZoneId: Z2FDTNDATAQYW2
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3029{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for alias with TTL")
	}
}

func TestE3029_AliasMissingDNSName(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyRecordSet:
    Type: AWS::Route53::RecordSet
    Properties:
      HostedZoneId: Z1234567890ABC
      Name: example.com
      Type: A
      AliasTarget:
        HostedZoneId: Z2FDTNDATAQYW2
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3029{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for alias missing DNSName")
	}
}

func TestE3029_AliasMissingHostedZoneId(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyRecordSet:
    Type: AWS::Route53::RecordSet
    Properties:
      HostedZoneId: Z1234567890ABC
      Name: example.com
      Type: A
      AliasTarget:
        DNSName: example.cloudfront.net
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3029{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for alias missing HostedZoneId")
	}
}
