// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3689{})
}

// E3689 validates MonitoringInterval and MonitoringRoleArn dependency.
type E3689 struct{}

func (r *E3689) ID() string { return "E3689" }

func (r *E3689) ShortDesc() string {
	return "MonitoringInterval requires MonitoringRoleArn"
}

func (r *E3689) Description() string {
	return "Validates that RDS DB instances with MonitoringInterval > 0 specify MonitoringRoleArn."
}

func (r *E3689) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3689"
}

func (r *E3689) Tags() []string {
	return []string{"resources", "properties", "rds", "monitoring"}
}

func (r *E3689) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::RDS::DBInstance" {
			continue
		}

		monitoringInterval, hasMonitoringInterval := res.Properties["MonitoringInterval"]
		if !hasMonitoringInterval || isIntrinsicFunction(monitoringInterval) {
			continue
		}

		var intervalValue int
		switch v := monitoringInterval.(type) {
		case int:
			intervalValue = v
		case float64:
			intervalValue = int(v)
		default:
			continue
		}

		if intervalValue > 0 {
			_, hasMonitoringRoleArn := res.Properties["MonitoringRoleArn"]
			if !hasMonitoringRoleArn {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': RDS DB Instance with MonitoringInterval > 0 must specify MonitoringRoleArn",
						resName,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	return matches
}
