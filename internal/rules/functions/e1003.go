// Package functions contains intrinsic function validation rules (E1xxx).
package functions

import (
	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1003{})
}

// E1003 checks that Description length is within limits.
type E1003 struct{}

func (r *E1003) ID() string { return "E1003" }

func (r *E1003) ShortDesc() string {
	return "Validate the max size of a description"
}

func (r *E1003) Description() string {
	return "Check if the size of the template description is less than the upper limit (1024 bytes)."
}

func (r *E1003) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cloudformation-limits.html"
}

func (r *E1003) Tags() []string {
	return []string{"template", "description", "limits"}
}

func (r *E1003) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check template description length (max 1024 bytes)
	if len(tmpl.Description) > 1024 {
		matches = append(matches, rules.Match{
			Message: "Template Description must not exceed 1024 bytes",
			Path:    []string{"Description"},
		})
	}

	return matches
}
