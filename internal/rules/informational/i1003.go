package informational

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&I1003{})
}

// I1003 checks if the template description is approaching the size limit.
type I1003 struct{}

func (r *I1003) ID() string { return "I1003" }

func (r *I1003) ShortDesc() string {
	return "Description approaching size limit"
}

func (r *I1003) Description() string {
	return "Checks if the template description is approaching the CloudFormation limit of 1,024 bytes. This rule provides an informational message when the description exceeds 80% of the limit (819 bytes)."
}

func (r *I1003) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cloudformation-limits.html"
}

func (r *I1003) Tags() []string {
	return []string{"template", "description", "limits"}
}

// MaxDescriptionSize is the CloudFormation limit for description size in bytes.
const MaxDescriptionSize = 1024

// DescriptionSizeWarningThreshold is 80% of the maximum description size.
const DescriptionSizeWarningThreshold = 819 // 80% of 1024

func (r *I1003) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	if tmpl.Description == "" {
		return matches
	}

	descSize := len(tmpl.Description)
	if descSize > DescriptionSizeWarningThreshold {
		percentage := int(float64(descSize) / float64(MaxDescriptionSize) * 100)
		matches = append(matches, rules.Match{
			Message: fmt.Sprintf("Description is %d bytes (%d%% of %d byte limit). Consider shortening the description.", descSize, percentage, MaxDescriptionSize),
			Line:    1,
			Column:  1,
			Path:    []string{"Description"},
		})
	}

	return matches
}
