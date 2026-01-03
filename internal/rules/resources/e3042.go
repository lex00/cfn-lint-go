// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3042{})
}

// E3042 checks that at least one container is marked as essential in ECS TaskDefinitions.
type E3042 struct{}

func (r *E3042) ID() string { return "E3042" }

func (r *E3042) ShortDesc() string {
	return "Essential container required"
}

func (r *E3042) Description() string {
	return "Validates that every AWS::ECS::TaskDefinition contains at least one ContainerDefinition that has Essential set to true."
}

func (r *E3042) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3042"
}

func (r *E3042) Tags() []string {
	return []string{"resources", "properties", "ecs", "task", "container"}
}

func (r *E3042) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::ECS::TaskDefinition" {
			continue
		}

		containerDefs, hasContainers := res.Properties["ContainerDefinitions"]
		if !hasContainers {
			continue
		}

		// Check if it's a slice
		containers, ok := containerDefs.([]interface{})
		if !ok {
			continue
		}

		// Check if at least one container has Essential = true
		hasEssential := false
		for _, container := range containers {
			containerMap, ok := container.(map[string]interface{})
			if !ok {
				continue
			}

			// Default is true if not specified
			essential, hasEssentialProp := containerMap["Essential"]
			if !hasEssentialProp {
				hasEssential = true
				break
			}

			// Check if Essential is explicitly true
			if essentialBool, ok := essential.(bool); ok && essentialBool {
				hasEssential = true
				break
			}
		}

		if !hasEssential {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': AWS::ECS::TaskDefinition must contain at least one ContainerDefinition with Essential set to true",
					resName,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties", "ContainerDefinitions"},
			})
		}
	}

	return matches
}
