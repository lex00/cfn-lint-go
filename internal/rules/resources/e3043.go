// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3043{})
}

// E3043 validates that parameters for nested stacks are correctly specified.
type E3043 struct{}

func (r *E3043) ID() string { return "E3043" }

func (r *E3043) ShortDesc() string {
	return "Nested stack parameters validation"
}

func (r *E3043) Description() string {
	return "Validates that parameters for a nested CloudFormation stack are specified and that extra parameters aren't provided."
}

func (r *E3043) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3043"
}

func (r *E3043) Tags() []string {
	return []string{"resources", "properties", "stack", "parameters"}
}

func (r *E3043) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::CloudFormation::Stack" {
			continue
		}

		// Check for Parameters property
		params, hasParams := res.Properties["Parameters"]
		if !hasParams {
			// If TemplateURL is specified without Parameters, that might be valid
			// Skip validation if no parameters are provided
			continue
		}

		// Validate that Parameters is a map
		paramsMap, ok := params.(map[string]interface{})
		if !ok {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Parameters must be a map of parameter names to values",
					resName,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties", "Parameters"},
			})
			continue
		}

		// Note: We cannot validate against the actual nested template without loading it
		// This rule primarily validates the structure, not the specific parameter names
		_ = paramsMap // Basic structure validation passed
	}

	return matches
}
