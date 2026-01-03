package informational

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&I6011{})
}

// I6011 checks if output names are approaching the length limit.
type I6011 struct{}

func (r *I6011) ID() string { return "I6011" }

func (r *I6011) ShortDesc() string {
	return "Output name approaching length limit"
}

func (r *I6011) Description() string {
	return "Checks if output names are approaching the CloudFormation limit of 255 characters. This rule provides an informational message when an output name exceeds 80% of the limit (204 characters)."
}

func (r *I6011) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cloudformation-limits.html"
}

func (r *I6011) Tags() []string {
	return []string{"outputs", "limits", "naming"}
}

// MaxOutputNameLength is the CloudFormation limit for output name length.
const MaxOutputNameLength = 255

// OutputNameLengthWarningThreshold is 80% of the maximum output name length.
const OutputNameLengthWarningThreshold = int(float64(MaxOutputNameLength) * 0.8)

func (r *I6011) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for name := range tmpl.Outputs {
		nameLen := len(name)
		if nameLen >= OutputNameLengthWarningThreshold {
			percentage := int(float64(nameLen) / float64(MaxOutputNameLength) * 100)
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Output name '%s' is %d characters (%d%% of %d character limit). Consider using a shorter name.", name, nameLen, percentage, MaxOutputNameLength),
				Path:    []string{"Outputs", name},
			})
		}
	}

	return matches
}
