// Package warnings contains warning-level rules (Wxxx).
package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W4001{})
}

// W4001 warns when AWS::CloudFormation::Interface references a non-existent parameter.
type W4001 struct{}

func (r *W4001) ID() string { return "W4001" }

func (r *W4001) ShortDesc() string {
	return "Interface references undefined parameter"
}

func (r *W4001) Description() string {
	return "Warns when AWS::CloudFormation::Interface metadata references a parameter that doesn't exist."
}

func (r *W4001) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-cloudformation-interface.html"
}

func (r *W4001) Tags() []string {
	return []string{"warnings", "metadata", "interface"}
}

func (r *W4001) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	if len(tmpl.Metadata) == 0 {
		return matches
	}

	cfnInterface, ok := tmpl.Metadata["AWS::CloudFormation::Interface"].(map[string]any)
	if !ok {
		return matches
	}

	// Check ParameterGroups
	if paramGroups, ok := cfnInterface["ParameterGroups"].([]any); ok {
		for i, group := range paramGroups {
			if groupMap, ok := group.(map[string]any); ok {
				if params, ok := groupMap["Parameters"].([]any); ok {
					for j, param := range params {
						if paramName, ok := param.(string); ok {
							if _, exists := tmpl.Parameters[paramName]; !exists {
								matches = append(matches, rules.Match{
									Message: fmt.Sprintf("ParameterGroups references undefined parameter '%s'", paramName),
									Path:    []string{"Metadata", "AWS::CloudFormation::Interface", "ParameterGroups", fmt.Sprintf("[%d]", i), "Parameters", fmt.Sprintf("[%d]", j)},
								})
							}
						}
					}
				}
			}
		}
	}

	// Check ParameterLabels
	if paramLabels, ok := cfnInterface["ParameterLabels"].(map[string]any); ok {
		for paramName := range paramLabels {
			if _, exists := tmpl.Parameters[paramName]; !exists {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("ParameterLabels references undefined parameter '%s'", paramName),
					Path:    []string{"Metadata", "AWS::CloudFormation::Interface", "ParameterLabels", paramName},
				})
			}
		}
	}

	return matches
}
