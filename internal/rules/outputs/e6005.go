// Package outputs contains output validation rules (E6xxx).
package outputs

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E6005{})
}

// E6005 checks that output Condition references a defined condition.
type E6005 struct{}

func (r *E6005) ID() string { return "E6005" }

func (r *E6005) ShortDesc() string {
	return "Output condition references undefined condition"
}

func (r *E6005) Description() string {
	return "Checks that output Condition attribute references a defined condition."
}

func (r *E6005) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/outputs-section-structure.html"
}

func (r *E6005) Tags() []string {
	return []string{"outputs", "conditions"}
}

func (r *E6005) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for outName, out := range tmpl.Outputs {
		if out.Condition != "" {
			if _, ok := tmpl.Conditions[out.Condition]; !ok {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Output '%s' references undefined condition '%s'", outName, out.Condition),
					Line:    out.Node.Line,
					Column:  out.Node.Column,
					Path:    []string{"Outputs", outName, "Condition"},
				})
			}
		}
	}

	return matches
}
