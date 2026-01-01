// Package outputs contains output validation rules (E6xxx).
package outputs

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E6010{})
}

// E6010 checks that the template does not exceed the output limit.
type E6010 struct{}

func (r *E6010) ID() string { return "E6010" }

func (r *E6010) ShortDesc() string {
	return "Output limit exceeded"
}

func (r *E6010) Description() string {
	return "Checks that the template does not exceed the maximum of 200 outputs."
}

func (r *E6010) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E6010"
}

func (r *E6010) Tags() []string {
	return []string{"outputs", "limits"}
}

// MaxOutputs is the CloudFormation limit for outputs per template.
const MaxOutputs = 200

func (r *E6010) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	count := len(tmpl.Outputs)
	if count > MaxOutputs {
		matches = append(matches, rules.Match{
			Message: fmt.Sprintf("Template has %d outputs, exceeding the limit of %d", count, MaxOutputs),
			Line:    1,
			Column:  1,
			Path:    []string{"Outputs"},
		})
	}

	return matches
}
