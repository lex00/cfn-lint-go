// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3053{})
}

// E3053 validates ECS task HostPort values.
type E3053 struct{}

func (r *E3053) ID() string { return "E3053" }

func (r *E3053) ShortDesc() string {
	return "ECS task HostPort values"
}

func (r *E3053) Description() string {
	return "Validates that ECS task definition HostPort is either undefined or equal to ContainerPort."
}

func (r *E3053) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3053"
}

func (r *E3053) Tags() []string {
	return []string{"resources", "properties", "ecs", "task", "ports"}
}

func (r *E3053) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::ECS::TaskDefinition" {
			continue
		}

		// Check if using awsvpc network mode
		isAwsvpc := false
		if networkMode, hasNetworkMode := res.Properties["NetworkMode"]; hasNetworkMode {
			if networkModeStr, ok := networkMode.(string); ok && networkModeStr == "awsvpc" {
				isAwsvpc = true
			}
		}

		// Check if Fargate (which requires awsvpc)
		if reqCompat, hasReqCompat := res.Properties["RequiresCompatibilities"]; hasReqCompat {
			if compatList, ok := reqCompat.([]interface{}); ok {
				for _, compat := range compatList {
					if compatStr, ok := compat.(string); ok && compatStr == "FARGATE" {
						isAwsvpc = true
						break
					}
				}
			}
		}

		if !isAwsvpc {
			continue
		}

		// Check container definitions
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

			portMappings, hasPorts := containerMap["PortMappings"]
			if !hasPorts {
				continue
			}

			ports, ok := portMappings.([]interface{})
			if !ok {
				continue
			}

			for j, port := range ports {
				portMap, ok := port.(map[string]interface{})
				if !ok {
					continue
				}

				containerPort, hasContainerPort := portMap["ContainerPort"]
				hostPort, hasHostPort := portMap["HostPort"]

				if !hasContainerPort {
					continue
				}

				containerPortInt, ok1 := r.toInt(containerPort)
				if !ok1 {
					continue
				}

				// If HostPort is specified, it must equal ContainerPort for awsvpc mode
				if hasHostPort {
					hostPortInt, ok2 := r.toInt(hostPort)
					if !ok2 {
						continue
					}

					if hostPortInt != 0 && hostPortInt != containerPortInt {
						matches = append(matches, rules.Match{
							Message: fmt.Sprintf(
								"Resource '%s': Container %d PortMapping %d has HostPort %d which must be undefined or equal to ContainerPort %d in awsvpc network mode",
								resName, i, j, hostPortInt, containerPortInt,
							),
							Line:   res.Node.Line,
							Column: res.Node.Column,
							Path:   []string{"Resources", resName, "Properties", "ContainerDefinitions", fmt.Sprintf("[%d]", i), "PortMappings", fmt.Sprintf("[%d]", j), "HostPort"},
						})
					}
				}
			}
		}
	}

	return matches
}

func (r *E3053) toInt(value interface{}) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case float64:
		return int(v), true
	case string:
		// Try to parse string to int
		var i int
		if _, err := fmt.Sscanf(v, "%d", &i); err == nil {
			return i, true
		}
	}
	return 0, false
}
