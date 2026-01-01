// Package mappings contains mapping validation rules (E7xxx).
package mappings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E7010{})
}

// E7010 checks that the template does not exceed the mapping limit.
type E7010 struct{}

func (r *E7010) ID() string { return "E7010" }

func (r *E7010) ShortDesc() string {
	return "Mapping limit exceeded"
}

func (r *E7010) Description() string {
	return "Checks that the template does not exceed the maximum of 200 mappings."
}

func (r *E7010) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E7010"
}

func (r *E7010) Tags() []string {
	return []string{"mappings", "limits"}
}

// MaxMappings is the CloudFormation limit for mappings per template.
const MaxMappings = 200

func (r *E7010) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	count := len(tmpl.Mappings)
	if count > MaxMappings {
		matches = append(matches, rules.Match{
			Message: fmt.Sprintf("Template has %d mappings, exceeding the limit of %d", count, MaxMappings),
			Line:    1,
			Column:  1,
			Path:    []string{"Mappings"},
		})
	}

	return matches
}
