// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3684{})
}

// E3684 validates target group health check protocol.
type E3684 struct{}

func (r *E3684) ID() string { return "E3684" }

func (r *E3684) ShortDesc() string {
	return "Validate target group health check protocol"
}

func (r *E3684) Description() string {
	return "Validates that AWS::ElasticLoadBalancingV2::TargetGroup HealthCheckProtocol is valid for the target type."
}

func (r *E3684) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3684"
}

func (r *E3684) Tags() []string {
	return []string{"resources", "properties", "elasticloadbalancingv2", "targetgroup"}
}

func (r *E3684) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::ElasticLoadBalancingV2::TargetGroup" {
			continue
		}

		healthCheckProtocol, hasHealthCheckProtocol := res.Properties["HealthCheckProtocol"]
		if !hasHealthCheckProtocol || isIntrinsicFunction(healthCheckProtocol) {
			continue
		}

		healthCheckProtocolStr, ok := healthCheckProtocol.(string)
		if !ok {
			continue
		}

		targetType, hasTargetType := res.Properties["TargetType"]
		if hasTargetType && !isIntrinsicFunction(targetType) {
			targetTypeStr, ok := targetType.(string)
			if !ok {
				continue
			}

			// Lambda target groups cannot have health checks
			if targetTypeStr == "lambda" {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': TargetGroup with TargetType 'lambda' must not specify HealthCheckProtocol",
						resName,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", "HealthCheckProtocol"},
				})
			}
		}

		// Validate protocol values
		validProtocols := map[string]bool{
			"HTTP":  true,
			"HTTPS": true,
			"TCP":   true,
			"TLS":   true,
			"UDP":   true,
			"TCP_UDP": true,
			"GENEVE": true,
		}

		if !validProtocols[healthCheckProtocolStr] {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Invalid HealthCheckProtocol '%s'",
					resName, healthCheckProtocolStr,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties", "HealthCheckProtocol"},
			})
		}
	}

	return matches
}
