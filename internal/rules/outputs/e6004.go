// Package outputs contains output validation rules (E6xxx).
package outputs

import (
	"fmt"
	"regexp"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E6004{})
}

// E6004 checks that output names follow naming conventions.
type E6004 struct{}

func (r *E6004) ID() string { return "E6004" }

func (r *E6004) ShortDesc() string {
	return "Output naming convention error"
}

func (r *E6004) Description() string {
	return "Checks that output names are alphanumeric and follow CloudFormation naming conventions."
}

func (r *E6004) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/outputs-section-structure.html"
}

func (r *E6004) Tags() []string {
	return []string{"outputs", "naming"}
}

// Output names must be alphanumeric
var validOutputNamePattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]*$`)

func (r *E6004) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for outName, out := range tmpl.Outputs {
		if !validOutputNamePattern.MatchString(outName) {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Output name '%s' must be alphanumeric and start with a letter", outName),
				Line:    out.Node.Line,
				Column:  out.Node.Column,
				Path:    []string{"Outputs", outName},
			})
		}
	}

	return matches
}
