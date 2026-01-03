package resources

import (
	"fmt"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3027_Metadata(t *testing.T) {
	rule := &E3027{}

	if rule.ID() != "E3027" {
		t.Errorf("Expected ID E3027, got %s", rule.ID())
	}
	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}
	if len(rule.Tags()) == 0 {
		t.Error("Tags should not be empty")
	}
}

func TestE3027_ValidRateExpression(t *testing.T) {
	testCases := []string{
		"rate(5 minutes)",
		"rate(1 hour)",
		"rate(7 days)",
	}

	for _, scheduleExpr := range testCases {
		yaml := fmt.Sprintf(`
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyRule:
    Type: AWS::Events::Rule
    Properties:
      ScheduleExpression: "%s"
      State: ENABLED
`, scheduleExpr)

		tmpl, err := template.Parse([]byte(yaml))
		if err != nil {
			t.Fatalf("Failed to parse template: %v", err)
		}

		rule := &E3027{}
		matches := rule.Match(tmpl)

		if len(matches) != 0 {
			t.Errorf("Expected 0 matches for valid rate expression '%s', got %d", scheduleExpr, len(matches))
		}
	}
}

func TestE3027_ValidCronExpression(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyRule:
    Type: AWS::Events::Rule
    Properties:
      ScheduleExpression: "cron(0 12 * * ? *)"
      State: ENABLED
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3027{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for valid cron expression, got %d", len(matches))
	}
}

func TestE3027_InvalidScheduleExpression(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyRule:
    Type: AWS::Events::Rule
    Properties:
      ScheduleExpression: "invalid"
      State: ENABLED
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3027{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for invalid schedule expression")
	}
}
