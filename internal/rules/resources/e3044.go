// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3044{})
}

// E3044 validates that ECS services with FARGATE or EXTERNAL launch type use REPLICA scheduling strategy.
type E3044 struct{}

func (r *E3044) ID() string { return "E3044" }

func (r *E3044) ShortDesc() string {
	return "FARGATE/EXTERNAL scheduling strategy"
}

func (r *E3044) Description() string {
	return "Validates that ECS services using FARGATE or EXTERNAL launch types only use REPLICA scheduling strategy."
}

func (r *E3044) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3044"
}

func (r *E3044) Tags() []string {
	return []string{"resources", "properties", "ecs", "service", "fargate"}
}

func (r *E3044) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::ECS::Service" {
			continue
		}

		launchType, hasLaunchType := res.Properties["LaunchType"]
		schedulingStrategy, hasSchedulingStrategy := res.Properties["SchedulingStrategy"]

		if !hasLaunchType {
			continue
		}

		launchTypeStr, ok := launchType.(string)
		if !ok {
			continue
		}

		// Check if FARGATE or EXTERNAL
		if launchTypeStr != "FARGATE" && launchTypeStr != "EXTERNAL" {
			continue
		}

		// If SchedulingStrategy is specified, it must be REPLICA
		if hasSchedulingStrategy {
			strategyStr, ok := schedulingStrategy.(string)
			if ok && strategyStr != "REPLICA" {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': ECS services with LaunchType '%s' must use SchedulingStrategy 'REPLICA' (got '%s')",
						resName, launchTypeStr, strategyStr,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", "SchedulingStrategy"},
				})
			}
		}
		// Note: If SchedulingStrategy is not specified, REPLICA is the default, which is valid
	}

	return matches
}
