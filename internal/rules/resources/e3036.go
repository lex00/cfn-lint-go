// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
	"gopkg.in/yaml.v3"
)

func init() {
	rules.Register(&E3036{})
}

// E3036 checks that UpdateReplacePolicy has a valid value.
type E3036 struct{}

func (r *E3036) ID() string { return "E3036" }

func (r *E3036) ShortDesc() string {
	return "Invalid UpdateReplacePolicy"
}

func (r *E3036) Description() string {
	return "Checks that UpdateReplacePolicy has a valid value (Delete, Retain, Snapshot)."
}

func (r *E3036) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-attribute-updatereplacepolicy.html"
}

func (r *E3036) Tags() []string {
	return []string{"resources", "updatereplacepolicy"}
}

var validUpdateReplacePolicies = map[string]bool{
	"Delete":   true,
	"Retain":   true,
	"Snapshot": true,
}

func (r *E3036) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Node == nil || res.Node.Kind != yaml.MappingNode {
			continue
		}

		for i := 0; i < len(res.Node.Content); i += 2 {
			key := res.Node.Content[i]
			value := res.Node.Content[i+1]

			if key.Value == "UpdateReplacePolicy" {
				policy := value.Value
				if !validUpdateReplacePolicies[policy] {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Invalid UpdateReplacePolicy '%s' in resource '%s'. Valid values: Delete, Retain, Snapshot", policy, resName),
						Line:    value.Line,
						Column:  value.Column,
						Path:    []string{"Resources", resName, "UpdateReplacePolicy"},
					})
				}
			}
		}
	}

	return matches
}
