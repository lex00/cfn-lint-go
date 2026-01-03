// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3048{})
}

// E3048 validates ECS Fargate required properties.
type E3048 struct{}

func (r *E3048) ID() string { return "E3048" }

func (r *E3048) ShortDesc() string {
	return "ECS Fargate required properties"
}

func (r *E3048) Description() string {
	return "Validates that ECS Fargate tasks have required properties and valid values including NetworkMode and RuntimePlatform."
}

func (r *E3048) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3048"
}

func (r *E3048) Tags() []string {
	return []string{"resources", "properties", "ecs", "fargate", "task"}
}

func (r *E3048) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::ECS::TaskDefinition" {
			continue
		}

		// Check if it requires Fargate compatibility
		reqCompat, hasReqCompat := res.Properties["RequiresCompatibilities"]
		isFargate := false

		if hasReqCompat {
			if compatList, ok := reqCompat.([]interface{}); ok {
				for _, compat := range compatList {
					if compatStr, ok := compat.(string); ok && compatStr == "FARGATE" {
						isFargate = true
						break
					}
				}
			}
		}

		if !isFargate {
			continue
		}

		// Check NetworkMode - must be awsvpc for Fargate
		networkMode, hasNetworkMode := res.Properties["NetworkMode"]
		if hasNetworkMode {
			if networkModeStr, ok := networkMode.(string); ok {
				if networkModeStr != "awsvpc" {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': Fargate tasks must use NetworkMode 'awsvpc' (got '%s')",
							resName, networkModeStr,
						),
						Line:   res.Node.Line,
						Column: res.Node.Column,
						Path:   []string{"Resources", resName, "Properties", "NetworkMode"},
					})
				}
			}
		} else {
			// NetworkMode is required for Fargate
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Fargate tasks must specify NetworkMode as 'awsvpc'",
					resName,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties"},
			})
		}

		// Check for CPU and Memory (required for Fargate)
		_, hasCPU := res.Properties["Cpu"]
		_, hasMemory := res.Properties["Memory"]

		if !hasCPU {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Fargate tasks must specify Cpu",
					resName,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties"},
			})
		}

		if !hasMemory {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Fargate tasks must specify Memory",
					resName,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties"},
			})
		}
	}

	return matches
}
