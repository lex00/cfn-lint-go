package resources

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE3615_ValidPeriods(t *testing.T) {
	validPeriods := []int{10, 30, 60, 120, 180, 300, 600, 3600}

	for _, period := range validPeriods {
		testYaml := createAlarmYAML(period)
		tmpl, err := template.Parse([]byte(testYaml))
		if err != nil {
			t.Fatalf("Failed to parse template for period %d: %v", period, err)
		}

		rule := &E3615{}
		matches := rule.Match(tmpl)

		if len(matches) != 0 {
			t.Errorf("Expected 0 matches for valid period %d, got %d: %v", period, len(matches), matches)
		}
	}
}

func TestE3615_InvalidPeriod(t *testing.T) {
	yaml := createAlarmYAML(45)

	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3615{}
	matches := rule.Match(tmpl)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for invalid period 45, got %d", len(matches))
	}
	if len(matches) > 0 && !containsString(matches[0].Message, "Period") {
		t.Errorf("Expected error about Period, got: %s", matches[0].Message)
	}
}

func TestE3615_InvalidPeriod90(t *testing.T) {
	yaml := createAlarmYAML(90)

	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3615{}
	matches := rule.Match(tmpl)

	if len(matches) != 1 {
		t.Errorf("Expected 1 match for invalid period 90, got %d", len(matches))
	}
}

func TestE3615_IntrinsicFunction(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyAlarm:
    Type: AWS::CloudWatch::Alarm
    Properties:
      ComparisonOperator: GreaterThanThreshold
      EvaluationPeriods: 1
      MetricName: CPUUtilization
      Namespace: AWS/EC2
      Period:
        Ref: AlarmPeriod
      Statistic: Average
      Threshold: 80
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	rule := &E3615{}
	matches := rule.Match(tmpl)

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for intrinsic function, got %d", len(matches))
	}
}

func TestE3615_Metadata(t *testing.T) {
	rule := &E3615{}

	if rule.ID() != "E3615" {
		t.Errorf("Expected ID E3615, got %s", rule.ID())
	}

	if rule.ShortDesc() == "" {
		t.Error("Expected non-empty ShortDesc")
	}

	if rule.Description() == "" {
		t.Error("Expected non-empty Description")
	}

	if rule.Source() == "" {
		t.Error("Expected non-empty Source")
	}

	tags := rule.Tags()
	if len(tags) == 0 {
		t.Error("Expected non-empty Tags")
	}
}

func createAlarmYAML(period int) string {
	return `AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyAlarm:
    Type: AWS::CloudWatch::Alarm
    Properties:
      ComparisonOperator: GreaterThanThreshold
      EvaluationPeriods: 1
      MetricName: CPUUtilization
      Namespace: AWS/EC2
      Period: ` + itoa(period) + `
      Statistic: Average
      Threshold: 80
`
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}

	negative := n < 0
	if negative {
		n = -n
	}

	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}

	if negative {
		digits = append([]byte{'-'}, digits...)
	}

	return string(digits)
}
