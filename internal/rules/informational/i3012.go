package informational

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&I3012{})
}

// I3012 checks if resource names are approaching the length limit.
type I3012 struct{}

func (r *I3012) ID() string { return "I3012" }

func (r *I3012) ShortDesc() string {
	return "Resource name approaching length limit"
}

func (r *I3012) Description() string {
	return "Checks if resource names are approaching the CloudFormation limit of 255 characters. This rule provides an informational message when a resource name exceeds 80% of the limit (204 characters)."
}

func (r *I3012) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cloudformation-limits.html"
}

func (r *I3012) Tags() []string {
	return []string{"resources", "limits", "naming"}
}

// MaxResourceNameLength is the CloudFormation limit for resource name length.
const MaxResourceNameLength = 255

// ResourceNameLengthWarningThreshold is 80% of the maximum resource name length.
const ResourceNameLengthWarningThreshold = int(float64(MaxResourceNameLength) * 0.8)

func (r *I3012) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for name := range tmpl.Resources {
		nameLen := len(name)
		if nameLen >= ResourceNameLengthWarningThreshold {
			percentage := int(float64(nameLen) / float64(MaxResourceNameLength) * 100)
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Resource name '%s' is %d characters (%d%% of %d character limit). Consider using a shorter name.", name, nameLen, percentage, MaxResourceNameLength),
				Path:    []string{"Resources", name},
			})
		}
	}

	return matches
}
