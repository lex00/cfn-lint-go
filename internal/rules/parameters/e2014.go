package parameters

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E2014{})
}

// E2014 checks that ConstraintDescription is only used with constraints.
type E2014 struct{}

func (r *E2014) ID() string { return "E2014" }

func (r *E2014) ShortDesc() string {
	return "Parameter ConstraintDescription usage"
}

func (r *E2014) Description() string {
	return "Checks that ConstraintDescription is only specified when constraints (AllowedPattern, AllowedValues, MinLength, MaxLength, MinValue, MaxValue) are defined."
}

func (r *E2014) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E2014"
}

func (r *E2014) Tags() []string {
	return []string{"parameters", "constraints"}
}

func (r *E2014) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for name, param := range tmpl.Parameters {
		// Check if ConstraintDescription is specified
		if param.ConstraintDescription != "" {
			// Check if any constraints are defined
			hasConstraints := param.AllowedPattern != "" ||
				len(param.AllowedValues) > 0 ||
				param.MinLength != nil ||
				param.MaxLength != nil ||
				param.MinValue != nil ||
				param.MaxValue != nil

			if !hasConstraints {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Parameter '%s' has ConstraintDescription but no constraints (AllowedPattern, AllowedValues, MinLength, MaxLength, MinValue, MaxValue) defined", name),
					Line:    param.Node.Line,
					Column:  param.Node.Column,
					Path:    []string{"Parameters", name, "ConstraintDescription"},
				})
			}
		}
	}

	return matches
}
