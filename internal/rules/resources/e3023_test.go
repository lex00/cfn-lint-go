package resources

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3023_Metadata(t *testing.T) {
	rule := &E3023{}

	if rule.ID() != "E3023" {
		t.Errorf("Expected ID E3023, got %s", rule.ID())
	}
	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}
	if len(rule.Tags()) == 0 {
		t.Error("Tags should not be empty")
	}
}

func TestE3023_ValidRecordSetWithResourceRecords(t *testing.T) {
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
      ResourceRecords:
        - 192.0.2.1
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3023{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid RecordSet, got %d", len(matches))
	}
}

func TestE3023_ValidRecordSetWithAlias(t *testing.T) {
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

	rule := &E3023{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid alias RecordSet, got %d", len(matches))
	}
}

func TestE3023_BothResourceRecordsAndAlias(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyRecordSet:
    Type: AWS::Route53::RecordSet
    Properties:
      HostedZoneId: Z1234567890ABC
      Name: example.com
      Type: A
      ResourceRecords:
        - 192.0.2.1
      AliasTarget:
        DNSName: example.cloudfront.net
        HostedZoneId: Z2FDTNDATAQYW2
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3023{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for RecordSet with both ResourceRecords and AliasTarget")
	}
}

func TestE3023_AliasWithTTL(t *testing.T) {
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

	rule := &E3023{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for alias RecordSet with TTL")
	}
}

func TestE3023_NeitherResourceRecordsNorAlias(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyRecordSet:
    Type: AWS::Route53::RecordSet
    Properties:
      HostedZoneId: Z1234567890ABC
      Name: example.com
      Type: A
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3023{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for RecordSet without ResourceRecords or AliasTarget")
	}
}
