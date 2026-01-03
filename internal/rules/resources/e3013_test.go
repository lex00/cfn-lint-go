package resources

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3013_Metadata(t *testing.T) {
	rule := &E3013{}

	if rule.ID() != "E3013" {
		t.Errorf("Expected ID E3013, got %s", rule.ID())
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
	if len(rule.Tags()) == 0 {
		t.Error("Tags should not be empty")
	}
}

func TestE3013_ValidAliases(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyDistribution:
    Type: AWS::CloudFront::Distribution
    Properties:
      DistributionConfig:
        Aliases:
          - example.com
          - www.example.com
          - "*.example.com"
        DefaultCacheBehavior:
          TargetOriginId: myOrigin
          ViewerProtocolPolicy: allow-all
        Enabled: true
        Origins:
          - Id: myOrigin
            DomainName: example.s3.amazonaws.com
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3013{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid aliases, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  %s", m.Message)
		}
	}
}

func TestE3013_InvalidAlias(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyDistribution:
    Type: AWS::CloudFront::Distribution
    Properties:
      DistributionConfig:
        Aliases:
          - invalid_domain
          - "not a domain!"
        DefaultCacheBehavior:
          TargetOriginId: myOrigin
          ViewerProtocolPolicy: allow-all
        Enabled: true
        Origins:
          - Id: myOrigin
            DomainName: example.s3.amazonaws.com
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3013{}
	matches := rule.Match(tmpl)

	if len(matches) != 2 {
		t.Errorf("Expected 2 matches for invalid aliases, got %d", len(matches))
	}
}

func TestE3013_NoAliases(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyDistribution:
    Type: AWS::CloudFront::Distribution
    Properties:
      DistributionConfig:
        DefaultCacheBehavior:
          TargetOriginId: myOrigin
          ViewerProtocolPolicy: allow-all
        Enabled: true
        Origins:
          - Id: myOrigin
            DomainName: example.s3.amazonaws.com
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3013{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when no aliases, got %d", len(matches))
	}
}

func TestE3013_SkipIntrinsicFunctions(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyDistribution:
    Type: AWS::CloudFront::Distribution
    Properties:
      DistributionConfig:
        Aliases: !Ref DomainNames
        DefaultCacheBehavior:
          TargetOriginId: myOrigin
          ViewerProtocolPolicy: allow-all
        Enabled: true
        Origins:
          - Id: myOrigin
            DomainName: example.s3.amazonaws.com
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3013{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when using intrinsic functions, got %d", len(matches))
	}
}
