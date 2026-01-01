// Package parameters contains parameter validation rules (E2xxx).
package parameters

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E2010{})
}

// E2010 checks that the template does not exceed the parameter limit.
type E2010 struct{}

func (r *E2010) ID() string { return "E2010" }

func (r *E2010) ShortDesc() string {
	return "Parameter limit exceeded"
}

func (r *E2010) Description() string {
	return "Checks that the template does not exceed the maximum of 200 parameters."
}

func (r *E2010) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E2010"
}

func (r *E2010) Tags() []string {
	return []string{"parameters", "limits"}
}

// MaxParameters is the CloudFormation limit for parameters per template.
const MaxParameters = 200

func (r *E2010) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	count := len(tmpl.Parameters)
	if count > MaxParameters {
		matches = append(matches, rules.Match{
			Message: fmt.Sprintf("Template has %d parameters, exceeding the limit of %d", count, MaxParameters),
			Line:    1,
			Column:  1,
			Path:    []string{"Parameters"},
		})
	}

	return matches
}
