// Package parameters contains parameter validation rules (E2xxx).
package parameters

import (
	"fmt"
	"regexp"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E2003{})
}

// E2003 checks that parameter names follow naming conventions.
type E2003 struct{}

func (r *E2003) ID() string { return "E2003" }

func (r *E2003) ShortDesc() string {
	return "Parameter naming convention error"
}

func (r *E2003) Description() string {
	return "Checks that parameter names are alphanumeric and follow CloudFormation naming conventions."
}

func (r *E2003) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/parameters-section-structure.html"
}

func (r *E2003) Tags() []string {
	return []string{"parameters", "naming"}
}

// Parameter names must be alphanumeric (A-Za-z0-9)
var validParamNamePattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]*$`)

func (r *E2003) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for paramName, param := range tmpl.Parameters {
		if !validParamNamePattern.MatchString(paramName) {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Parameter name '%s' must be alphanumeric and start with a letter", paramName),
				Line:    param.Node.Line,
				Column:  param.Node.Column,
				Path:    []string{"Parameters", paramName},
			})
		}
	}

	return matches
}
