// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3661{})
}

// E3661 validates Route53 health check AlarmIdentifier.
type E3661 struct{}

func (r *E3661) ID() string { return "E3661" }

func (r *E3661) ShortDesc() string {
	return "Validate Route53 health check AlarmIdentifier"
}

func (r *E3661) Description() string {
	return "Validates that AWS::Route53::HealthCheck resources with Type CLOUDWATCH_METRIC have AlarmIdentifier specified."
}

func (r *E3661) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3661"
}

func (r *E3661) Tags() []string {
	return []string{"resources", "properties", "route53", "healthcheck"}
}

func (r *E3661) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Route53::HealthCheck" {
			continue
		}

		healthCheckConfig, hasConfig := res.Properties["HealthCheckConfig"]
		if !hasConfig || isIntrinsicFunction(healthCheckConfig) {
			continue
		}

		configMap, ok := healthCheckConfig.(map[string]any)
		if !ok {
			continue
		}

		healthCheckType, hasType := configMap["Type"]
		if !hasType || isIntrinsicFunction(healthCheckType) {
			continue
		}

		typeStr, ok := healthCheckType.(string)
		if !ok {
			continue
		}

		if typeStr == "CLOUDWATCH_METRIC" {
			_, hasAlarmIdentifier := configMap["AlarmIdentifier"]
			if !hasAlarmIdentifier {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': Route53 HealthCheck with Type 'CLOUDWATCH_METRIC' must specify AlarmIdentifier",
						resName,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", "HealthCheckConfig"},
				})
			}
		}
	}

	return matches
}
