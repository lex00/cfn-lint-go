// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3680{})
}

// E3680 validates ALB minimum subnets.
type E3680 struct{}

func (r *E3680) ID() string { return "E3680" }

func (r *E3680) ShortDesc() string {
	return "ALB requires minimum 2 subnets"
}

func (r *E3680) Description() string {
	return "Validates that AWS::ElasticLoadBalancingV2::LoadBalancer Application Load Balancers specify at least 2 subnets."
}

func (r *E3680) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3680"
}

func (r *E3680) Tags() []string {
	return []string{"resources", "properties", "elasticloadbalancingv2", "loadbalancer"}
}

func (r *E3680) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::ElasticLoadBalancingV2::LoadBalancer" {
			continue
		}

		// Check if it's an Application Load Balancer (default is application)
		lbType, hasType := res.Properties["Type"]
		if hasType && !isIntrinsicFunction(lbType) {
			typeStr, ok := lbType.(string)
			if ok && typeStr != "application" {
				continue
			}
		}

		subnets, hasSubnets := res.Properties["Subnets"]
		if !hasSubnets || isIntrinsicFunction(subnets) {
			continue
		}

		subnetList, ok := subnets.([]any)
		if !ok {
			continue
		}

		if len(subnetList) < 2 {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Application Load Balancer must specify at least 2 subnets. Found: %d",
					resName, len(subnetList),
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties", "Subnets"},
			})
		}
	}

	return matches
}
