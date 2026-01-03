package resources

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE2529_ValidLogGroupWithNoFilters(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyLogGroup:
    Type: AWS::Logs::LogGroup
    Properties:
      LogGroupName: /aws/lambda/my-function
      RetentionInDays: 7
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2529{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for LogGroup with no SubscriptionFilters, got %d", len(matches))
	}
}

func TestE2529_ValidLogGroupWithOneFilter(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyLogGroup:
    Type: AWS::Logs::LogGroup
    Properties:
      LogGroupName: /aws/lambda/my-function
      SubscriptionFilters:
        - FilterName: MyFilter
          DestinationArn: arn:aws:lambda:us-east-1:123456789012:function:my-function
          FilterPattern: "[...]"
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2529{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for LogGroup with 1 SubscriptionFilter, got %d", len(matches))
	}
}

func TestE2529_ValidLogGroupWithTwoFilters(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyLogGroup:
    Type: AWS::Logs::LogGroup
    Properties:
      LogGroupName: /aws/lambda/my-function
      SubscriptionFilters:
        - FilterName: MyFilter1
          DestinationArn: arn:aws:lambda:us-east-1:123456789012:function:my-function1
          FilterPattern: "[...]"
        - FilterName: MyFilter2
          DestinationArn: arn:aws:lambda:us-east-1:123456789012:function:my-function2
          FilterPattern: "[...]"
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2529{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for LogGroup with 2 SubscriptionFilters (at limit), got %d", len(matches))
	}
}

func TestE2529_InvalidLogGroupWithThreeFilters(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyLogGroup:
    Type: AWS::Logs::LogGroup
    Properties:
      LogGroupName: /aws/lambda/my-function
      SubscriptionFilters:
        - FilterName: MyFilter1
          DestinationArn: arn:aws:lambda:us-east-1:123456789012:function:my-function1
          FilterPattern: "[...]"
        - FilterName: MyFilter2
          DestinationArn: arn:aws:lambda:us-east-1:123456789012:function:my-function2
          FilterPattern: "[...]"
        - FilterName: MyFilter3
          DestinationArn: arn:aws:lambda:us-east-1:123456789012:function:my-function3
          FilterPattern: "[...]"
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2529{}
	matches := rule.Match(tmpl)

	if len(matches) == 0 {
		t.Error("Expected match for LogGroup with 3 SubscriptionFilters (exceeds limit)")
	}

	if len(matches) > 0 {
		expectedMsg := "LogGroup 'MyLogGroup' has 3 SubscriptionFilters, exceeding the limit of 2"
		if matches[0].Message != expectedMsg {
			t.Errorf("Expected message '%s', got '%s'", expectedMsg, matches[0].Message)
		}
	}
}

func TestE2529_SkipsIntrinsicFunctions(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyLogGroup:
    Type: AWS::Logs::LogGroup
    Properties:
      LogGroupName: /aws/lambda/my-function
      SubscriptionFilters: !Ref Filters
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E2529{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches when SubscriptionFilters is an intrinsic function, got %d", len(matches))
	}
}

func TestE2529_Metadata(t *testing.T) {
	rule := &E2529{}

	if rule.ID() != "E2529" {
		t.Errorf("Expected ID E2529, got %s", rule.ID())
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

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Tags should not be empty")
	}
}
