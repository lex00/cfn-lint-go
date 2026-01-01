// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
	"gopkg.in/yaml.v3"
)

func init() {
	rules.Register(&E3035{})
}

// E3035 checks that DeletionPolicy has a valid value.
type E3035 struct{}

func (r *E3035) ID() string { return "E3035" }

func (r *E3035) ShortDesc() string {
	return "Invalid DeletionPolicy"
}

func (r *E3035) Description() string {
	return "Checks that DeletionPolicy has a valid value (Delete, Retain, Snapshot, RetainExceptOnCreate)."
}

func (r *E3035) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-attribute-deletionpolicy.html"
}

func (r *E3035) Tags() []string {
	return []string{"resources", "deletionpolicy"}
}

var validDeletionPolicies = map[string]bool{
	"Delete":               true,
	"Retain":               true,
	"Snapshot":             true,
	"RetainExceptOnCreate": true,
}

func (r *E3035) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Node == nil || res.Node.Kind != yaml.MappingNode {
			continue
		}

		for i := 0; i < len(res.Node.Content); i += 2 {
			key := res.Node.Content[i]
			value := res.Node.Content[i+1]

			if key.Value == "DeletionPolicy" {
				policy := value.Value
				if !validDeletionPolicies[policy] {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Invalid DeletionPolicy '%s' in resource '%s'. Valid values: Delete, Retain, Snapshot, RetainExceptOnCreate", policy, resName),
						Line:    value.Line,
						Column:  value.Column,
						Path:    []string{"Resources", resName, "DeletionPolicy"},
					})
				}
			}
		}
	}

	return matches
}
