package resources

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3021_RequirementSatisfied(t *testing.T) {
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

	rule := &E3021{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when requirement satisfied, got %d", len(matches))
		for _, m := range matches {
			t.Logf("  Match: %s", m.Message)
		}
	}
}

func TestE3021_RequirementMissing(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyDB:
    Type: AWS::RDS::DBInstance
    Properties:
      DBInstanceClass: db.t3.micro
      Engine: mysql
      MasterUsername: admin
      AllocatedStorage: 20
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3021{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for missing required dependency")
	} else {
		if !strings.Contains(matches[0].Message, "requires") {
			t.Errorf("Error message should mention 'requires': %s", matches[0].Message)
		}
		if !strings.Contains(matches[0].Message, "MasterUserPassword") {
			t.Errorf("Error message should mention MasterUserPassword: %s", matches[0].Message)
		}
	}
}

func TestE3021_TriggerPropertyNotPresent(t *testing.T) {
	// If MasterUsername is not present, MasterUserPassword is not required
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

	rule := &E3021{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when trigger property not present, got %d", len(matches))
	}
}

func TestE3021_SQSContentBasedDeduplication(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyQueue:
    Type: AWS::SQS::Queue
    Properties:
      QueueName: my-queue.fifo
      ContentBasedDeduplication: true
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3021{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match - ContentBasedDeduplication requires FifoQueue")
	}
}

func TestE3021_UnknownResourceSkipped(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyCustom:
    Type: AWS::Unknown::Resource
    Properties:
      PropA: valueA
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	rule := &E3021{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for unknown resource type, got %d", len(matches))
	}
}

func TestE3021_Metadata(t *testing.T) {
	rule := &E3021{}

	if rule.ID() != "E3021" {
		t.Errorf("Expected ID E3021, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
