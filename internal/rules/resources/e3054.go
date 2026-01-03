// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3054{})
}

// E3054 validates ECS Fargate task definition compatibility.
type E3054 struct{}

func (r *E3054) ID() string { return "E3054" }

func (r *E3054) ShortDesc() string {
	return "ECS Fargate TaskDefinition compatibility"
}

func (r *E3054) Description() string {
	return "Validates that ECS services with Fargate LaunchType reference task definitions with FARGATE in RequiresCompatibilities."
}

func (r *E3054) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3054"
}

func (r *E3054) Tags() []string {
	return []string{"resources", "properties", "ecs", "fargate", "service"}
}

func (r *E3054) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::ECS::Service" {
			continue
		}

		// Check if LaunchType is FARGATE
		launchType, hasLaunchType := res.Properties["LaunchType"]
		if !hasLaunchType {
			continue
		}

		launchTypeStr, ok := launchType.(string)
		if !ok || launchTypeStr != "FARGATE" {
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

		// Find the referenced task definition
		if taskDefRef != "" {
			if taskDefRes, exists := tmpl.Resources[taskDefRef]; exists && taskDefRes.Type == "AWS::ECS::TaskDefinition" {
				// Check RequiresCompatibilities
				reqCompat, hasReqCompat := taskDefRes.Properties["RequiresCompatibilities"]
				if !hasReqCompat {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': ECS service uses Fargate LaunchType but TaskDefinition '%s' does not specify RequiresCompatibilities",
							resName, taskDefRef,
						),
						Line:   taskDefRes.Node.Line,
						Column: taskDefRes.Node.Column,
						Path:   []string{"Resources", taskDefRef, "Properties"},
					})
					continue
				}

				// Check if FARGATE is in the list
				hasFargate := false
				if compatList, ok := reqCompat.([]interface{}); ok {
					for _, compat := range compatList {
						if compatStr, ok := compat.(string); ok && compatStr == "FARGATE" {
							hasFargate = true
							break
						}
					}
				}

				if !hasFargate {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': ECS service uses Fargate LaunchType but TaskDefinition '%s' RequiresCompatibilities does not include FARGATE",
							resName, taskDefRef,
						),
						Line:   taskDefRes.Node.Line,
						Column: taskDefRes.Node.Column,
						Path:   []string{"Resources", taskDefRef, "Properties", "RequiresCompatibilities"},
					})
				}
			}
		}
	}

	return matches
}
