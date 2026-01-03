package informational

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&I7002{})
}

// I7002 checks if mapping names are approaching the length limit.
type I7002 struct{}

func (r *I7002) ID() string { return "I7002" }

func (r *I7002) ShortDesc() string {
	return "Mapping name approaching length limit"
}

func (r *I7002) Description() string {
	return "Checks if mapping names are approaching the CloudFormation limit of 255 characters. This rule provides an informational message when a mapping name exceeds 80% of the limit (204 characters)."
}

func (r *I7002) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cloudformation-limits.html"
}

func (r *I7002) Tags() []string {
	return []string{"mappings", "limits", "naming"}
}

// MaxMappingNameLength is the CloudFormation limit for mapping name length.
const MaxMappingNameLength = 255

// MappingNameLengthWarningThreshold is 80% of the maximum mapping name length.
const MappingNameLengthWarningThreshold = int(float64(MaxMappingNameLength) * 0.8)

func (r *I7002) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for name := range tmpl.Mappings {
		nameLen := len(name)
		if nameLen >= MappingNameLengthWarningThreshold {
			percentage := int(float64(nameLen) / float64(MaxMappingNameLength) * 100)
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Mapping name '%s' is %d characters (%d%% of %d character limit). Consider using a shorter name.", name, nameLen, percentage, MaxMappingNameLength),
				Path:    []string{"Mappings", name},
			})
		}
	}

	return matches
}
