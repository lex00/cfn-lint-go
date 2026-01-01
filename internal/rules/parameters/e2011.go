// Package parameters contains parameter validation rules (E2xxx).
package parameters

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E2011{})
}

// E2011 checks that parameter names don't exceed the length limit.
type E2011 struct{}

func (r *E2011) ID() string { return "E2011" }

func (r *E2011) ShortDesc() string {
	return "Parameter name length error"
}

func (r *E2011) Description() string {
	return "Checks that parameter names don't exceed 255 characters."
}

func (r *E2011) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cloudformation-limits.html"
}

func (r *E2011) Tags() []string {
	return []string{"parameters", "limits"}
}

const maxParamNameLength = 255

func (r *E2011) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for paramName, param := range tmpl.Parameters {
		if len(paramName) > maxParamNameLength {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Parameter name '%s' exceeds maximum length of %d characters (got %d)", paramName, maxParamNameLength, len(paramName)),
				Line:    param.Node.Line,
				Column:  param.Node.Column,
				Path:    []string{"Parameters", paramName},
			})
		}
	}

	return matches
}
