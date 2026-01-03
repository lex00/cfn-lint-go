package resources

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3026_Metadata(t *testing.T) {
	rule := &E3026{}

	if rule.ID() != "E3026" {
		t.Errorf("Expected ID E3026, got %s", rule.ID())
	}
	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}
	if len(rule.Tags()) == 0 {
		t.Error("Tags should not be empty")
	}
}

func TestE3026_ClusterModeWithFailover(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyReplicationGroup:
    Type: AWS::ElastiCache::ReplicationGroup
    Properties:
      ReplicationGroupDescription: Test
      ClusterMode: enabled
      AutomaticFailoverEnabled: true
      NumNodeGroups: 2
      ReplicasPerNodeGroup: 1
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3026{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid cluster mode with failover, got %d", len(matches))
	}
}

func TestE3026_ClusterModeWithoutFailover(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyReplicationGroup:
    Type: AWS::ElastiCache::ReplicationGroup
    Properties:
      ReplicationGroupDescription: Test
      ClusterMode: enabled
      AutomaticFailoverEnabled: false
      NumNodeGroups: 2
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3026{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for cluster mode without automatic failover")
	}
}
