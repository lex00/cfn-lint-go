// Package rulessection contains Rules section validation rules (E17xx).
package rulessection

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
	"gopkg.in/yaml.v3"
)

func init() {
	rules.Register(&E1700{})
}

// E1700 checks that Rules have appropriate configuration.
type E1700 struct{}

func (r *E1700) ID() string { return "E1700" }

func (r *E1700) ShortDesc() string {
	return "Rules have appropriate configuration"
}

func (r *E1700) Description() string {
	return "Validates that each Rule has proper configuration with required Assertions property and optional RuleCondition."
}

func (r *E1700) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E1700"
}

func (r *E1700) Tags() []string {
	return []string{"rules", "configuration"}
}

// Valid rule properties per CloudFormation spec
var validRuleProperties = map[string]bool{
	"RuleCondition": true,
	"Assertions":    true,
}

func (r *E1700) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for name, rule := range tmpl.Rules {
		// Check if rule is not a mapping
		if rule.Node != nil && rule.Node.Kind != yaml.MappingNode {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Rule '%s' must be a mapping with Assertions property", name),
				Line:    rule.Node.Line,
				Column:  rule.Node.Column,
				Path:    []string{"Rules", name},
			})
			continue
		}

		// Check for required Assertions property
		if len(rule.Assertions) == 0 {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Rule '%s' is missing required property 'Assertions'", name),
				Line:    rule.Node.Line,
				Column:  rule.Node.Column,
				Path:    []string{"Rules", name},
			})
		}

		// Check for invalid properties
		if rule.Node != nil && rule.Node.Kind == yaml.MappingNode {
			for i := 0; i < len(rule.Node.Content); i += 2 {
				propKey := rule.Node.Content[i]
				if !validRuleProperties[propKey.Value] {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Rule '%s' has invalid property '%s'. Valid properties are: Assertions, RuleCondition", name, propKey.Value),
						Line:    propKey.Line,
						Column:  propKey.Column,
						Path:    []string{"Rules", name, propKey.Value},
					})
				}
			}
		}
	}

	return matches
}
