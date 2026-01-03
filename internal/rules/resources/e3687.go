// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3687{})
}

// E3687 validates protocol port validation.
type E3687 struct{}

func (r *E3687) ID() string { return "E3687" }

func (r *E3687) ShortDesc() string {
	return "Validate protocol port validation"
}

func (r *E3687) Description() string {
	return "Validates that security group rules specify valid port ranges for the protocol."
}

func (r *E3687) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3687"
}

func (r *E3687) Tags() []string {
	return []string{"resources", "properties", "ec2", "securitygroup"}
}

func (r *E3687) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	checkRule := func(resName string, rule map[string]any, rulePath []string, res *template.Resource) {
		ipProtocol, hasProtocol := rule["IpProtocol"]
		if !hasProtocol || isIntrinsicFunction(ipProtocol) {
			return
		}

		protocolStr, ok := ipProtocol.(string)
		if !ok {
			return
		}

		fromPort, hasFromPort := rule["FromPort"]
		toPort, hasToPort := rule["ToPort"]

		// ICMP (protocol 1) uses FromPort/ToPort for type/code
		if protocolStr == "icmp" || protocolStr == "1" {
			// ICMP can have FromPort and ToPort for type and code
			return
		}

		// TCP (6) and UDP (17) require port ranges
		if protocolStr == "tcp" || protocolStr == "6" || protocolStr == "udp" || protocolStr == "17" {
			if !hasFromPort || !hasToPort {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': Security group rule with protocol '%s' must specify FromPort and ToPort",
						resName, protocolStr,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   rulePath,
				})
			}
		}

		// Validate port numbers if specified
		if hasFromPort && !isIntrinsicFunction(fromPort) {
			var fromPortNum int
			switch v := fromPort.(type) {
			case int:
				fromPortNum = v
			case float64:
				fromPortNum = int(v)
			default:
				return
			}

			if fromPortNum != -1 && (fromPortNum < 0 || fromPortNum > 65535) {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': FromPort must be between 0 and 65535, or -1. Got: %d",
						resName, fromPortNum,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   append(rulePath, "FromPort"),
				})
			}
		}

		if hasToPort && !isIntrinsicFunction(toPort) {
			var toPortNum int
			switch v := toPort.(type) {
			case int:
				toPortNum = v
			case float64:
				toPortNum = int(v)
			default:
				return
			}

			if toPortNum != -1 && (toPortNum < 0 || toPortNum > 65535) {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': ToPort must be between 0 and 65535, or -1. Got: %d",
						resName, toPortNum,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   append(rulePath, "ToPort"),
				})
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
