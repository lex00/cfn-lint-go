// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3046{})
}

// E3046 validates awslogs configuration in ECS task definitions.
type E3046 struct{}

func (r *E3046) ID() string { return "E3046" }

func (r *E3046) ShortDesc() string {
	return "ECS task awslogs configuration"
}

func (r *E3046) Description() string {
	return "Validates that ECS task definitions using the awslogs log driver have required 'awslogs-group' and 'awslogs-region' options."
}

func (r *E3046) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3046"
}

func (r *E3046) Tags() []string {
	return []string{"resources", "properties", "ecs", "task", "logging"}
}

func (r *E3046) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::ECS::TaskDefinition" {
			continue
		}

		containerDefs, hasContainers := res.Properties["ContainerDefinitions"]
		if !hasContainers {
			continue
		}

		containers, ok := containerDefs.([]interface{})
		if !ok {
			continue
		}

		for i, container := range containers {
			containerMap, ok := container.(map[string]interface{})
			if !ok {
				continue
			}

			logConfig, hasLogConfig := containerMap["LogConfiguration"]
			if !hasLogConfig {
				continue
			}

			logConfigMap, ok := logConfig.(map[string]interface{})
			if !ok {
				continue
			}

			logDriver, hasDriver := logConfigMap["LogDriver"]
			if !hasDriver {
				continue
			}

			logDriverStr, ok := logDriver.(string)
			if !ok || logDriverStr != "awslogs" {
				continue
			}

			// Check for required options
			options, hasOptions := logConfigMap["Options"]
			if !hasOptions {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': Container %d using awslogs driver must specify Options with 'awslogs-group' and 'awslogs-region'",
						resName, i,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", "ContainerDefinitions", fmt.Sprintf("[%d]", i), "LogConfiguration"},
				})
				continue
			}

			optionsMap, ok := options.(map[string]interface{})
			if !ok {
				continue
			}

			_, hasGroup := optionsMap["awslogs-group"]
			_, hasRegion := optionsMap["awslogs-region"]

			if !hasGroup || !hasRegion {
				missing := []string{}
				if !hasGroup {
					missing = append(missing, "awslogs-group")
				}
				if !hasRegion {
					missing = append(missing, "awslogs-region")
				}

				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': Container %d using awslogs driver is missing required options: %v",
						resName, i, missing,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", "ContainerDefinitions", fmt.Sprintf("[%d]", i), "LogConfiguration", "Options"},
				})
			}
		}
	}

	return matches
}
