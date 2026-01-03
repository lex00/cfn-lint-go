package resources

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3016_Metadata(t *testing.T) {
	rule := &E3016{}

	if rule.ID() != "E3016" {
		t.Errorf("Expected ID E3016, got %s", rule.ID())
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

func TestE3016_ValidUpdatePolicy(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyASG:
    Type: AWS::AutoScaling::AutoScalingGroup
    UpdatePolicy:
      AutoScalingRollingUpdate:
        MaxBatchSize: 2
        MinInstancesInService: 1
        PauseTime: PT5M
    Properties:
      MinSize: 1
      MaxSize: 5
      LaunchTemplate:
        LaunchTemplateId: !Ref MyLaunchTemplate
        Version: !GetAtt MyLaunchTemplate.LatestVersionNumber
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3016{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid UpdatePolicy, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  %s", m.Message)
		}
	}
}

func TestE3016_InvalidUpdatePolicyKey(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyASG:
    Type: AWS::AutoScaling::AutoScalingGroup
    UpdatePolicy:
      InvalidKey:
        SomeProperty: value
    Properties:
      MinSize: 1
      MaxSize: 5
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3016{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for invalid UpdatePolicy key")
	}
}

func TestE3016_UpdatePolicyNotObject(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyASG:
    Type: AWS::AutoScaling::AutoScalingGroup
    UpdatePolicy: "invalid"
    Properties:
      MinSize: 1
      MaxSize: 5
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3016{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for UpdatePolicy not being an object")
	}
}
