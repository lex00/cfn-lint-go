package informational

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&I2010{})
}

// I2010 checks if the parameter count is approaching the CloudFormation limit.
type I2010 struct{}

func (r *I2010) ID() string { return "I2010" }

func (r *I2010) ShortDesc() string {
	return "Parameter count approaching limit"
}

func (r *I2010) Description() string {
	return "Checks if the template is approaching the maximum of 200 parameters. This rule provides an informational message when the parameter count exceeds 80% of the limit (160 parameters)."
}

func (r *I2010) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cloudformation-limits.html"
}

func (r *I2010) Tags() []string {
	return []string{"parameters", "limits"}
}

// MaxParameters is the CloudFormation limit for parameters per template.
const MaxParametersInfo = 200

// ParameterCountWarningThreshold is 80% of the maximum parameter count.
const ParameterCountWarningThreshold = int(float64(MaxParametersInfo) * 0.8)

func (r *I2010) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	count := len(tmpl.Parameters)
	if count >= ParameterCountWarningThreshold && count < MaxParametersInfo {
		percentage := int(float64(count) / float64(MaxParametersInfo) * 100)
		matches = append(matches, rules.Match{
			Message: fmt.Sprintf("Template has %d parameters (%d%% of %d limit). Consider reducing the number of parameters or using nested stacks.", count, percentage, MaxParametersInfo),
			Line:    1,
			Column:  1,
			Path:    []string{"Parameters"},
		})
	}

	return matches
}
