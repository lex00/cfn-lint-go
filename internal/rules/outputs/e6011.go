// Package outputs contains output validation rules (E6xxx).
package outputs

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E6011{})
}

// E6011 checks that output names don't exceed the length limit.
type E6011 struct{}

func (r *E6011) ID() string { return "E6011" }

func (r *E6011) ShortDesc() string {
	return "Output name length error"
}

func (r *E6011) Description() string {
	return "Checks that output names don't exceed 255 characters."
}

func (r *E6011) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cloudformation-limits.html"
}

func (r *E6011) Tags() []string {
	return []string{"outputs", "limits"}
}

const maxOutputNameLength = 255

func (r *E6011) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for outName, out := range tmpl.Outputs {
		if len(outName) > maxOutputNameLength {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Output name '%s' exceeds maximum length of %d characters (got %d)", outName, maxOutputNameLength, len(outName)),
				Line:    out.Node.Line,
				Column:  out.Node.Column,
				Path:    []string{"Outputs", outName},
			})
		}
	}

	return matches
}
