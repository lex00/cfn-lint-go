// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3052{})
}

// E3052 validates ECS service NetworkConfiguration.
type E3052 struct{}

func (r *E3052) ID() string { return "E3052" }

func (r *E3052) ShortDesc() string {
	return "ECS service NetworkConfiguration"
}

func (r *E3052) Description() string {
	return "Validates that ECS services with awsvpc network mode have NetworkConfiguration specified."
}

func (r *E3052) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3052"
}

func (r *E3052) Tags() []string {
	return []string{"resources", "properties", "ecs", "service", "networking"}
}

func (r *E3052) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::ECS::Service" {
			continue
		}

		// Get the task definition reference
		taskDef, hasTaskDef := res.Properties["TaskDefinition"]
		if !hasTaskDef {
			continue
		}

		// Check if TaskDefinition is a Ref
		taskDefRef := ""
		if taskDefStr, ok := taskDef.(string); ok {
			taskDefRef = taskDefStr
		} else if taskDefMap, ok := taskDef.(map[string]interface{}); ok {
			if ref, hasRef := taskDefMap["Ref"]; hasRef {
				if refStr, ok := ref.(string); ok {
					taskDefRef = refStr
				}
			}
		}

		// Find the referenced task definition and check its NetworkMode
		requiresNetworkConfig := false
		if taskDefRef != "" {
			if taskDefRes, exists := tmpl.Resources[taskDefRef]; exists && taskDefRes.Type == "AWS::ECS::TaskDefinition" {
				if networkMode, hasNetworkMode := taskDefRes.Properties["NetworkMode"]; hasNetworkMode {
					if networkModeStr, ok := networkMode.(string); ok && networkModeStr == "awsvpc" {
						requiresNetworkConfig = true
					}
				}
			}
		}

		// Also check LaunchType FARGATE which implies awsvpc
		if launchType, hasLaunchType := res.Properties["LaunchType"]; hasLaunchType {
			if launchTypeStr, ok := launchType.(string); ok && launchTypeStr == "FARGATE" {
				requiresNetworkConfig = true
			}
		}

		if requiresNetworkConfig {
			networkConfig, hasNetworkConfig := res.Properties["NetworkConfiguration"]
			if !hasNetworkConfig {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': ECS service with awsvpc network mode must specify NetworkConfiguration",
						resName,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties"},
				})
			} else {
				// Validate NetworkConfiguration structure
				if netConfigMap, ok := networkConfig.(map[string]interface{}); ok {
					if awsvpcConfig, hasAwsvpc := netConfigMap["AwsvpcConfiguration"]; !hasAwsvpc {
						matches = append(matches, rules.Match{
							Message: fmt.Sprintf(
								"Resource '%s': NetworkConfiguration must include AwsvpcConfiguration",
								resName,
							),
							Line:   res.Node.Line,
							Column: res.Node.Column,
							Path:   []string{"Resources", resName, "Properties", "NetworkConfiguration"},
						})
					} else {
						// Validate AwsvpcConfiguration has Subnets
						if awsvpcMap, ok := awsvpcConfig.(map[string]interface{}); ok {
							if _, hasSubnets := awsvpcMap["Subnets"]; !hasSubnets {
								matches = append(matches, rules.Match{
									Message: fmt.Sprintf(
										"Resource '%s': AwsvpcConfiguration must specify Subnets",
										resName,
									),
									Line:   res.Node.Line,
									Column: res.Node.Column,
									Path:   []string{"Resources", resName, "Properties", "NetworkConfiguration", "AwsvpcConfiguration"},
								})
							}
						}
					}
				}
			}
		}
	}

	return matches
}
