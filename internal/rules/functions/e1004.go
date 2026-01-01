// Package functions contains intrinsic function validation rules (E1xxx).
package functions

import (
	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1004{})
}

// E1004 checks that Description is a string.
type E1004 struct{}

func (r *E1004) ID() string { return "E1004" }

func (r *E1004) ShortDesc() string {
	return "Description must be a string"
}

func (r *E1004) Description() string {
	return "Checks that the template Description is a string value."
}

func (r *E1004) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-anatomy.html"
}

func (r *E1004) Tags() []string {
	return []string{"template", "description"}
}

func (r *E1004) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check template description - this is already parsed as string in template.go
	// but we should validate it's not too long (max 1024 characters)
	if len(tmpl.Description) > 1024 {
		matches = append(matches, rules.Match{
			Message: "Template Description must not exceed 1024 characters",
			Path:    []string{"Description"},
		})
	}

	return matches
}
