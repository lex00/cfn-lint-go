package informational

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&I7010{})
}

// I7010 checks if the mapping count is approaching the CloudFormation limit.
type I7010 struct{}

func (r *I7010) ID() string { return "I7010" }

func (r *I7010) ShortDesc() string {
	return "Mapping count approaching limit"
}

func (r *I7010) Description() string {
	return "Checks if the template is approaching the maximum of 200 mappings. This rule provides an informational message when the mapping count exceeds 80% of the limit (160 mappings)."
}

func (r *I7010) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cloudformation-limits.html"
}

func (r *I7010) Tags() []string {
	return []string{"mappings", "limits"}
}

// MaxMappings is the CloudFormation limit for mappings per template.
const MaxMappings = 200

// MappingCountWarningThreshold is 80% of the maximum mapping count.
const MappingCountWarningThreshold = int(float64(MaxMappings) * 0.8)

func (r *I7010) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	count := len(tmpl.Mappings)
	if count >= MappingCountWarningThreshold && count < MaxMappings {
		percentage := int(float64(count) / float64(MaxMappings) * 100)
		matches = append(matches, rules.Match{
			Message: fmt.Sprintf("Template has %d mappings (%d%% of %d limit). Consider reducing the number of mappings or using nested stacks.", count, percentage, MaxMappings),
			Line:    1,
			Column:  1,
			Path:    []string{"Mappings"},
		})
	}

	return matches
}
