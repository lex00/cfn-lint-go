// Package metadata contains metadata validation rules (E4xxx).
package metadata

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E4001{})
}

// E4001 checks that AWS::CloudFormation::Interface metadata is properly configured.
type E4001 struct{}

func (r *E4001) ID() string { return "E4001" }

func (r *E4001) ShortDesc() string {
	return "Interface metadata error"
}

func (r *E4001) Description() string {
	return "Checks that AWS::CloudFormation::Interface metadata references valid parameters."
}

func (r *E4001) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-cloudformation-interface.html"
}

func (r *E4001) Tags() []string {
	return []string{"metadata", "interface"}
}

func (r *E4001) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check if Metadata has AWS::CloudFormation::Interface
	if tmpl.Metadata == nil {
		return matches
	}

	interfaceMeta, ok := tmpl.Metadata["AWS::CloudFormation::Interface"]
	if !ok {
		return matches
	}

	interfaceMap, ok := interfaceMeta.(map[string]any)
	if !ok {
		return matches
	}

	// Check ParameterGroups
	if groups, ok := interfaceMap["ParameterGroups"]; ok {
		groupsList, ok := groups.([]any)
		if ok {
			for i, group := range groupsList {
				groupMap, ok := group.(map[string]any)
				if !ok {
					continue
				}

				// Check Parameters in each group
				if params, ok := groupMap["Parameters"]; ok {
					paramsList, ok := params.([]any)
					if ok {
						for _, param := range paramsList {
							paramName, ok := param.(string)
							if ok {
								if !tmpl.HasParameter(paramName) {
									matches = append(matches, rules.Match{
										Message: fmt.Sprintf("ParameterGroups[%d] references undefined parameter '%s'", i, paramName),
										Path:    []string{"Metadata", "AWS::CloudFormation::Interface", "ParameterGroups"},
									})
								}
							}
						}
					}
				}
			}
		}
	}

	// Check ParameterLabels
	if labels, ok := interfaceMap["ParameterLabels"]; ok {
		labelsMap, ok := labels.(map[string]any)
		if ok {
			for paramName := range labelsMap {
				if !tmpl.HasParameter(paramName) {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("ParameterLabels references undefined parameter '%s'", paramName),
						Path:    []string{"Metadata", "AWS::CloudFormation::Interface", "ParameterLabels"},
					})
				}
			}
		}
	}

	return matches
}
