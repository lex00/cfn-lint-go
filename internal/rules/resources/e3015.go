// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3015{})
}

// E3015 checks that resource Condition references a defined condition.
type E3015 struct{}

func (r *E3015) ID() string { return "E3015" }

func (r *E3015) ShortDesc() string {
	return "Resource condition references undefined condition"
}

func (r *E3015) Description() string {
	return "Checks that resource Condition attribute references a defined condition."
}

func (r *E3015) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/conditions-section-structure.html"
}

func (r *E3015) Tags() []string {
	return []string{"resources", "conditions"}
}

func (r *E3015) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Condition != "" {
			if _, ok := tmpl.Conditions[res.Condition]; !ok {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Resource '%s' references undefined condition '%s'", resName, res.Condition),
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", resName, "Condition"},
				})
			}
		}
	}

	return matches
}
