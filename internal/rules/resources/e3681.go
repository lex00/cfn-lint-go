// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3681{})
}

// E3681 validates target group target type restrictions.
type E3681 struct{}

func (r *E3681) ID() string { return "E3681" }

func (r *E3681) ShortDesc() string {
	return "Validate target group target type restrictions"
}

func (r *E3681) Description() string {
	return "Validates that AWS::ElasticLoadBalancingV2::TargetGroup resources with TargetType 'lambda' do not specify certain properties."
}

func (r *E3681) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3681"
}

func (r *E3681) Tags() []string {
	return []string{"resources", "properties", "elasticloadbalancingv2", "targetgroup"}
}

func (r *E3681) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::ElasticLoadBalancingV2::TargetGroup" {
			continue
		}

		targetType, hasTargetType := res.Properties["TargetType"]
		if !hasTargetType || isIntrinsicFunction(targetType) {
			continue
		}

		targetTypeStr, ok := targetType.(string)
		if !ok {
			continue
		}

		if targetTypeStr == "lambda" {
			// Lambda target groups cannot have Port, Protocol, or VpcId
			invalidProps := []string{"Port", "Protocol", "VpcId"}
			for _, prop := range invalidProps {
				if _, hasProp := res.Properties[prop]; hasProp {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': TargetGroup with TargetType 'lambda' must not specify %s",
							resName, prop,
						),
						Line:   res.Node.Line,
						Column: res.Node.Column,
						Path:   []string{"Resources", resName, "Properties", prop},
					})
				}
			}
		}
	}

	return matches
}
