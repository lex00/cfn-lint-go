package informational

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&I2011{})
}

// I2011 checks if parameter names are approaching the length limit.
type I2011 struct{}

func (r *I2011) ID() string { return "I2011" }

func (r *I2011) ShortDesc() string {
	return "Parameter name approaching length limit"
}

func (r *I2011) Description() string {
	return "Checks if parameter names are approaching the CloudFormation limit of 255 characters. This rule provides an informational message when a parameter name exceeds 80% of the limit (204 characters)."
}

func (r *I2011) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cloudformation-limits.html"
}

func (r *I2011) Tags() []string {
	return []string{"parameters", "limits", "naming"}
}

// MaxParameterNameLength is the CloudFormation limit for parameter name length.
const MaxParameterNameLength = 255

// ParameterNameLengthWarningThreshold is 80% of the maximum parameter name length.
const ParameterNameLengthWarningThreshold = int(float64(MaxParameterNameLength) * 0.8)

func (r *I2011) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for name := range tmpl.Parameters {
		nameLen := len(name)
		if nameLen >= ParameterNameLengthWarningThreshold {
			percentage := int(float64(nameLen) / float64(MaxParameterNameLength) * 100)
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Parameter name '%s' is %d characters (%d%% of %d character limit). Consider using a shorter name.", name, nameLen, percentage, MaxParameterNameLength),
				Path:    []string{"Parameters", name},
			})
		}
	}

	return matches
}
