// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3683{})
}

// E3683 validates target group protocol restrictions.
type E3683 struct{}

func (r *E3683) ID() string { return "E3683" }

func (r *E3683) ShortDesc() string {
	return "Validate target group protocol restrictions"
}

func (r *E3683) Description() string {
	return "Validates that AWS::ElasticLoadBalancingV2::TargetGroup Protocol matches allowed values for target type."
}

func (r *E3683) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3683"
}

func (r *E3683) Tags() []string {
	return []string{"resources", "properties", "elasticloadbalancingv2", "targetgroup"}
}

func (r *E3683) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::ElasticLoadBalancingV2::TargetGroup" {
			continue
		}

		protocol, hasProtocol := res.Properties["Protocol"]
		if !hasProtocol || isIntrinsicFunction(protocol) {
			continue
		}

		protocolStr, ok := protocol.(string)
		if !ok {
			continue
		}

		targetType, hasTargetType := res.Properties["TargetType"]
		if hasTargetType && !isIntrinsicFunction(targetType) {
			targetTypeStr, ok := targetType.(string)
			if !ok {
				continue
			}

			// Validate protocol based on target type
			if targetTypeStr == "lambda" {
				if protocolStr != "HTTP" && protocolStr != "HTTPS" {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': TargetGroup with TargetType 'lambda' must use HTTP or HTTPS protocol. Got: %s",
							resName, protocolStr,
						),
						Line:   res.Node.Line,
						Column: res.Node.Column,
						Path:   []string{"Resources", resName, "Properties", "Protocol"},
					})
				}
			}
		}
	}

	return matches
}
