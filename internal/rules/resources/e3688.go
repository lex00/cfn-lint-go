// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3688{})
}

// E3688 validates port -1 validation.
type E3688 struct{}

func (r *E3688) ID() string { return "E3688" }

func (r *E3688) ShortDesc() string {
	return "Validate port -1 validation"
}

func (r *E3688) Description() string {
	return "Validates that security group rules with port -1 use protocol -1 (all protocols)."
}

func (r *E3688) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3688"
}

func (r *E3688) Tags() []string {
	return []string{"resources", "properties", "ec2", "securitygroup"}
}

func (r *E3688) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	checkRule := func(resName string, rule map[string]any, rulePath []string, res *template.Resource) {
		fromPort, hasFromPort := rule["FromPort"]
		toPort, hasToPort := rule["ToPort"]
		ipProtocol, hasProtocol := rule["IpProtocol"]

		var fromPortNum, toPortNum int
		var hasFromPortNum, hasToPortNum bool

		if hasFromPort && !isIntrinsicFunction(fromPort) {
			switch v := fromPort.(type) {
			case int:
				fromPortNum = v
				hasFromPortNum = true
			case float64:
				fromPortNum = int(v)
				hasFromPortNum = true
			}
		}

		if hasToPort && !isIntrinsicFunction(toPort) {
			switch v := toPort.(type) {
			case int:
				toPortNum = v
				hasToPortNum = true
			case float64:
				toPortNum = int(v)
				hasToPortNum = true
			}
		}

		// If FromPort or ToPort is -1, IpProtocol must be -1
		if hasFromPortNum && fromPortNum == -1 || hasToPortNum && toPortNum == -1 {
			if hasProtocol && !isIntrinsicFunction(ipProtocol) {
				protocolStr, ok := ipProtocol.(string)
				if ok && protocolStr != "-1" {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': Security group rule with FromPort or ToPort set to -1 must use IpProtocol -1",
							resName,
						),
						Line:   res.Node.Line,
						Column: res.Node.Column,
						Path:   rulePath,
					})
				}
			}
		}
	}

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::EC2::SecurityGroup" && res.Type != "AWS::EC2::SecurityGroupIngress" && res.Type != "AWS::EC2::SecurityGroupEgress" {
			continue
		}

		if res.Type == "AWS::EC2::SecurityGroup" {
			// Check ingress rules
			if ingress, hasIngress := res.Properties["SecurityGroupIngress"]; hasIngress && !isIntrinsicFunction(ingress) {
				if ingressList, ok := ingress.([]any); ok {
					for _, rule := range ingressList {
						if ruleMap, ok := rule.(map[string]any); ok {
							checkRule(resName, ruleMap, []string{"Resources", resName, "Properties", "SecurityGroupIngress"}, res)
						}
					}
				}
			}

			// Check egress rules
			if egress, hasEgress := res.Properties["SecurityGroupEgress"]; hasEgress && !isIntrinsicFunction(egress) {
				if egressList, ok := egress.([]any); ok {
					for _, rule := range egressList {
						if ruleMap, ok := rule.(map[string]any); ok {
							checkRule(resName, ruleMap, []string{"Resources", resName, "Properties", "SecurityGroupEgress"}, res)
						}
					}
				}
			}
		} else {
			// Individual ingress/egress resource
			checkRule(resName, res.Properties, []string{"Resources", resName, "Properties"}, res)
		}
	}

	return matches
}
