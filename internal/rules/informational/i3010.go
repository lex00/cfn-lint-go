package informational

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&I3010{})
}

// I3010 checks if the resource count is approaching the CloudFormation limit.
type I3010 struct{}

func (r *I3010) ID() string { return "I3010" }

func (r *I3010) ShortDesc() string {
	return "Resource count approaching limit"
}

func (r *I3010) Description() string {
	return "Checks if the template is approaching the maximum of 500 resources. This rule provides an informational message when the resource count exceeds 80% of the limit (400 resources)."
}

func (r *I3010) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cloudformation-limits.html"
}

func (r *I3010) Tags() []string {
	return []string{"resources", "limits"}
}

// MaxResources is the CloudFormation limit for resources per template.
const MaxResources = 500

// ResourceCountWarningThreshold is 80% of the maximum resource count.
const ResourceCountWarningThreshold = int(float64(MaxResources) * 0.8)

func (r *I3010) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	count := len(tmpl.Resources)
	if count >= ResourceCountWarningThreshold && count < MaxResources {
		percentage := int(float64(count) / float64(MaxResources) * 100)
		matches = append(matches, rules.Match{
			Message: fmt.Sprintf("Template has %d resources (%d%% of %d limit). Consider splitting into nested stacks.", count, percentage, MaxResources),
			Line:    1,
			Column:  1,
			Path:    []string{"Resources"},
		})
	}

	return matches
}
