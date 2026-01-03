// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3049{})
}

// E3049 validates ECS dynamic host port configuration with load balancers.
type E3049 struct{}

func (r *E3049) ID() string { return "E3049" }

func (r *E3049) ShortDesc() string {
	return "ECS dynamic host port configuration"
}

func (r *E3049) Description() string {
	return "Validates that ECS services using dynamic host ports (0) with load balancers specify 'traffic-port' as HealthCheckPort."
}

func (r *E3049) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3049"
}

func (r *E3049) Tags() []string {
	return []string{"resources", "properties", "ecs", "service", "loadbalancer"}
}

func (r *E3049) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::ECS::Service" {
			continue
		}

		// Check for LoadBalancers
		loadBalancers, hasLB := res.Properties["LoadBalancers"]
		if !hasLB {
			continue
		}

		lbList, ok := loadBalancers.([]interface{})
		if !ok {
			continue
		}

		// Get task definition to check for dynamic ports
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

		// Find the referenced task definition
		hasDynamicPort := false
		if taskDefRef != "" {
			if taskDefRes, exists := tmpl.Resources[taskDefRef]; exists && taskDefRes.Type == "AWS::ECS::TaskDefinition" {
				if containerDefs, hasContainers := taskDefRes.Properties["ContainerDefinitions"]; hasContainers {
					if containers, ok := containerDefs.([]interface{}); ok {
						for _, container := range containers {
							if containerMap, ok := container.(map[string]interface{}); ok {
								if portMappings, hasPorts := containerMap["PortMappings"]; hasPorts {
									if ports, ok := portMappings.([]interface{}); ok {
										for _, port := range ports {
											if portMap, ok := port.(map[string]interface{}); ok {
												if hostPort, hasHostPort := portMap["HostPort"]; hasHostPort {
													if hostPortInt, ok := hostPort.(int); ok && hostPortInt == 0 {
														hasDynamicPort = true
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}

		if !hasDynamicPort {
			continue
		}

		// Check load balancers for proper health check configuration
		for i, lb := range lbList {
			lbMap, ok := lb.(map[string]interface{})
			if !ok {
				continue
			}

			// Check for TargetGroupArn
			targetGroupArn, hasTargetGroup := lbMap["TargetGroupArn"]
			if !hasTargetGroup {
				continue
			}

			// Find the target group resource
			targetGroupRef := ""
			if tgMap, ok := targetGroupArn.(map[string]interface{}); ok {
				if ref, hasRef := tgMap["Ref"]; hasRef {
					if refStr, ok := ref.(string); ok {
						targetGroupRef = refStr
					}
				}
			}

			if targetGroupRef != "" {
				if tgRes, exists := tmpl.Resources[targetGroupRef]; exists {
					if tgRes.Type == "AWS::ElasticLoadBalancingV2::TargetGroup" {
						healthCheckPort, hasHealthCheckPort := tgRes.Properties["HealthCheckPort"]
						if !hasHealthCheckPort {
							matches = append(matches, rules.Match{
								Message: fmt.Sprintf(
									"Resource '%s': When using dynamic host ports (HostPort: 0), LoadBalancer target group '%s' must specify HealthCheckPort as 'traffic-port'",
									resName, targetGroupRef,
								),
								Line:   res.Node.Line,
								Column: res.Node.Column,
								Path:   []string{"Resources", resName, "Properties", "LoadBalancers", fmt.Sprintf("[%d]", i)},
							})
						} else if healthCheckPortStr, ok := healthCheckPort.(string); ok && healthCheckPortStr != "traffic-port" {
							matches = append(matches, rules.Match{
								Message: fmt.Sprintf(
									"Resource '%s': When using dynamic host ports, HealthCheckPort should be 'traffic-port' (got '%s')",
									resName, healthCheckPortStr,
								),
								Line:   res.Node.Line,
								Column: res.Node.Column,
								Path:   []string{"Resources", targetGroupRef, "Properties", "HealthCheckPort"},
							})
						}
					}
				}
			}
		}
	}

	return matches
}
