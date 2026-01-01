// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3005{})
}

// E3005 checks that DependsOn references existing resources.
type E3005 struct{}

func (r *E3005) ID() string { return "E3005" }

func (r *E3005) ShortDesc() string {
	return "DependsOn references undefined resource"
}

func (r *E3005) Description() string {
	return "Checks that all DependsOn values reference valid resource logical IDs."
}

func (r *E3005) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-attribute-dependson.html"
}

func (r *E3005) Tags() []string {
	return []string{"resources", "dependencies"}
}

func (r *E3005) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		for _, dep := range res.DependsOn {
			if !tmpl.HasResource(dep) {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("DependsOn references undefined resource '%s' in resource '%s'", dep, resName),
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", resName, "DependsOn"},
				})
			}
			// Also check that resource doesn't depend on itself
			if dep == resName {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Resource '%s' cannot depend on itself", resName),
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", resName, "DependsOn"},
				})
			}
		}
	}

	return matches
}
