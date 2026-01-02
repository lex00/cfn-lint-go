package resources

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3020_NoExclusionViolation(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyDB:
    Type: AWS::RDS::DBInstance
    Properties:
      DBInstanceClass: db.t3.micro
      Engine: mysql
      MasterUsername: admin
      MasterUserPassword: mypassword
      AllocatedStorage: 20
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3020{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when no exclusion violation, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE3020_ExclusionViolation(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyDB:
    Type: AWS::RDS::DBInstance
    Properties:
      DBInstanceClass: db.t3.micro
      Engine: mysql
      DBSnapshotIdentifier: my-snapshot
      MasterUsername: admin
      AllocatedStorage: 20
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3020{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for exclusion violation")
	} else {
		if !strings.Contains(matches[0].Message, "cannot be used with") {
			t.Errorf("Error message should mention 'cannot be used with': %s", matches[0].Message)
		}
		if !strings.Contains(matches[0].Message, "DBSnapshotIdentifier") {
			t.Errorf("Error message should mention DBSnapshotIdentifier: %s", matches[0].Message)
		}
	}
}

func TestE3020_SnapshotWithoutExcludedProps(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyDB:
    Type: AWS::RDS::DBInstance
    Properties:
      DBInstanceClass: db.t3.micro
      DBSnapshotIdentifier: my-snapshot
      AllocatedStorage: 20
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3020{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when using snapshot without excluded props, got %d", len(matches))
	}
}

func TestE3020_UnknownResourceSkipped(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyCustom:
    Type: AWS::Unknown::Resource
    Properties:
      PropA: valueA
      PropB: valueB
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3020{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for unknown resource type, got %d", len(matches))
	}
}

func TestE3020_Metadata(t *testing.T) {
	rule := &E3020{}

	if rule.ID() != "E3020" {
		t.Errorf("Expected ID E3020, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
