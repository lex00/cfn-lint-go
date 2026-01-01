// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3010{})
}

// E3010 checks that the template does not exceed the resource limit.
type E3010 struct{}

func (r *E3010) ID() string { return "E3010" }

func (r *E3010) ShortDesc() string {
	return "Resource limit exceeded"
}

func (r *E3010) Description() string {
	return "Checks that the template does not exceed the maximum of 500 resources."
}

func (r *E3010) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3010"
}

func (r *E3010) Tags() []string {
	return []string{"resources", "limits"}
}

// MaxResources is the CloudFormation limit for resources per template.
const MaxResources = 500

func (r *E3010) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	count := len(tmpl.Resources)
	if count > MaxResources {
		matches = append(matches, rules.Match{
			Message: fmt.Sprintf("Template has %d resources, exceeding the limit of %d", count, MaxResources),
			Line:    1,
			Column:  1,
			Path:    []string{"Resources"},
		})
	}

	return matches
}
