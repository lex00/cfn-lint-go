// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3615{})
}

// E3615 validates CloudWatch Alarm period values.
type E3615 struct{}

func (r *E3615) ID() string { return "E3615" }

func (r *E3615) ShortDesc() string {
	return "Validate CloudWatch Alarm period"
}

func (r *E3615) Description() string {
	return "Validates that the Period property in AWS::CloudWatch::Alarm resources is one of the valid values: 10, 30, 60, or any multiple of 60."
}

func (r *E3615) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3615"
}

func (r *E3615) Tags() []string {
	return []string{"resources", "properties", "cloudwatch", "alarm"}
}

func (r *E3615) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::CloudWatch::Alarm" {
			continue
		}

		period, hasPeriod := res.Properties["Period"]
		if !hasPeriod {
			continue
		}

		// Skip intrinsic functions
		if isIntrinsicFunction(period) {
			continue
		}

		// Period can be int or float64 (from YAML/JSON parsing)
		var periodValue int
		switch v := period.(type) {
		case int:
			periodValue = v
		case float64:
			periodValue = int(v)
		default:
			continue // Skip non-numeric values
		}

		// Valid values: 10, 30, 60, or multiples of 60
		isValid := false
		if periodValue == 10 || periodValue == 30 {
			isValid = true
		} else if periodValue >= 60 && periodValue%60 == 0 {
			isValid = true
		}

		if !isValid {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': CloudWatch Alarm Period must be 10, 30, 60, or a multiple of 60. Got: %d",
					resName, periodValue,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties", "Period"},
			})
		}
	}

	return matches
}
