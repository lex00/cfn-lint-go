package rulessection

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
	"gopkg.in/yaml.v3"
)

func init() {
	rules.Register(&E1701{})
}

// E1701 validates the configuration of Assertions.
type E1701 struct{}

func (r *E1701) ID() string { return "E1701" }

func (r *E1701) ShortDesc() string {
	return "Validate Assertions configuration"
}

func (r *E1701) Description() string {
	return "Validates that Assertions in Rules are properly configured with required Assert property and optional AssertDescription."
}

func (r *E1701) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E1701"
}

func (r *E1701) Tags() []string {
	return []string{"rules", "assertions"}
}

// Valid assertion properties per CloudFormation spec
var validAssertionProperties = map[string]bool{
	"Assert":            true,
	"AssertDescription": true,
}

func (r *E1701) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for ruleName, rule := range tmpl.Rules {
		// Validate each assertion
		for idx, assertion := range rule.Assertions {
			// Check if assertion is not a mapping
			if assertion.Node != nil && assertion.Node.Kind != yaml.MappingNode {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Rule '%s' Assertion[%d] must be a mapping with Assert property", ruleName, idx),
					Line:    assertion.Node.Line,
					Column:  assertion.Node.Column,
					Path:    []string{"Rules", ruleName, "Assertions", fmt.Sprintf("[%d]", idx)},
				})
				continue
			}

			// Check for required Assert property
			if assertion.Assert == nil {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Rule '%s' Assertion[%d] is missing required property 'Assert'", ruleName, idx),
					Line:    assertion.Node.Line,
					Column:  assertion.Node.Column,
					Path:    []string{"Rules", ruleName, "Assertions", fmt.Sprintf("[%d]", idx)},
				})
			}

			// Check for invalid properties
			if assertion.Node != nil && assertion.Node.Kind == yaml.MappingNode {
				for i := 0; i < len(assertion.Node.Content); i += 2 {
					propKey := assertion.Node.Content[i]
					if !validAssertionProperties[propKey.Value] {
						matches = append(matches, rules.Match{
							Message: fmt.Sprintf("Rule '%s' Assertion[%d] has invalid property '%s'. Valid properties are: Assert, AssertDescription", ruleName, idx, propKey.Value),
							Line:    propKey.Line,
							Column:  propKey.Column,
							Path:    []string{"Rules", ruleName, "Assertions", fmt.Sprintf("[%d]", idx), propKey.Value},
						})
					}
				}
			}
		}
	}

	return matches
}
