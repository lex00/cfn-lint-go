// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"regexp"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3011{})
}

// E3011 checks that resource property names are valid.
type E3011 struct{}

func (r *E3011) ID() string { return "E3011" }

func (r *E3011) ShortDesc() string {
	return "Invalid property name"
}

func (r *E3011) Description() string {
	return "Checks that resource property names are alphanumeric."
}

func (r *E3011) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/resources-section-structure.html"
}

func (r *E3011) Tags() []string {
	return []string{"resources", "properties"}
}

// Property names must be alphanumeric (PascalCase typically)
var validPropNamePattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]*$`)

func (r *E3011) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		for propName := range res.Properties {
			if !validPropNamePattern.MatchString(propName) {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Property name '%s' in resource '%s' must be alphanumeric", propName, resName),
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", resName, "Properties", propName},
				})
			}
		}
	}

	return matches
}
