package informational

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&I6010{})
}

// I6010 checks if the output count is approaching the CloudFormation limit.
type I6010 struct{}

func (r *I6010) ID() string { return "I6010" }

func (r *I6010) ShortDesc() string {
	return "Output count approaching limit"
}

func (r *I6010) Description() string {
	return "Checks if the template is approaching the maximum of 200 outputs. This rule provides an informational message when the output count exceeds 80% of the limit (160 outputs)."
}

func (r *I6010) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cloudformation-limits.html"
}

func (r *I6010) Tags() []string {
	return []string{"outputs", "limits"}
}

// MaxOutputs is the CloudFormation limit for outputs per template.
const MaxOutputs = 200

// OutputCountWarningThreshold is 80% of the maximum output count.
const OutputCountWarningThreshold = int(float64(MaxOutputs) * 0.8)

func (r *I6010) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	count := len(tmpl.Outputs)
	if count >= OutputCountWarningThreshold && count < MaxOutputs {
		percentage := int(float64(count) / float64(MaxOutputs) * 100)
		matches = append(matches, rules.Match{
			Message: fmt.Sprintf("Template has %d outputs (%d%% of %d limit). Consider reducing the number of outputs or using nested stacks.", count, percentage, MaxOutputs),
			Line:    1,
			Column:  1,
			Path:    []string{"Outputs"},
		})
	}

	return matches
}
