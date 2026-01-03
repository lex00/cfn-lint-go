package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W3687{})
}

// W3687 warns when security group rules specify ports for protocols that don't use them.
type W3687 struct{}

func (r *W3687) ID() string { return "W3687" }

func (r *W3687) ShortDesc() string {
	return "Ports not for certain protocols"
}

func (r *W3687) Description() string {
	return "Warns when security group rules specify FromPort/ToPort for protocols like ICMP or ALL (-1) that don't use traditional ports."
}

func (r *W3687) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-security-group-ingress.html"
}

func (r *W3687) Tags() []string {
	return []string{"warnings", "ec2", "security-group", "ports"}
}

// Protocols that don't use traditional ports
var protocolsWithoutPorts = map[string]bool{
	"-1":   true, // All traffic
	"all":  true,
	"icmp": true,
	"58":   true, // ICMPv6
}

func (r *W3687) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		switch res.Type {
		case "AWS::EC2::SecurityGroup":
			r.checkSecurityGroup(resName, res, &matches)
		case "AWS::EC2::SecurityGroupIngress":
			r.checkSecurityGroupRule(resName, res.Properties, "Ingress", &matches)
		case "AWS::EC2::SecurityGroupEgress":
			r.checkSecurityGroupRule(resName, res.Properties, "Egress", &matches)
		}
	}

	return matches
}

func (r *W3687) checkSecurityGroup(resName string, res *template.Resource, matches *[]rules.Match) {
	// Check SecurityGroupIngress
	if ingress, ok := res.Properties["SecurityGroupIngress"].([]any); ok {
		for i, rule := range ingress {
			if ruleMap, ok := rule.(map[string]any); ok {
				r.checkRule(resName, ruleMap, fmt.Sprintf("SecurityGroupIngress[%d]", i), matches)
			}
		}
	}

	// Check SecurityGroupEgress
	if egress, ok := res.Properties["SecurityGroupEgress"].([]any); ok {
		for i, rule := range egress {
			if ruleMap, ok := rule.(map[string]any); ok {
				r.checkRule(resName, ruleMap, fmt.Sprintf("SecurityGroupEgress[%d]", i), matches)
			}
		}
	}
}

func (r *W3687) checkSecurityGroupRule(resName string, props map[string]any, ruleType string, matches *[]rules.Match) {
	r.checkRule(resName, props, ruleType, matches)
}

func (r *W3687) checkRule(resName string, rule map[string]any, context string, matches *[]rules.Match) {
	protocol := r.getProtocol(rule["IpProtocol"])
	if protocol == "" {
		return
	}

	// Check if this protocol doesn't use ports
	if !protocolsWithoutPorts[protocol] {
		return
	}

	fromPort := r.getPort(rule["FromPort"])
	toPort := r.getPort(rule["ToPort"])

	// For protocol -1 (all), ports should not be specified or should be -1
	if protocol == "-1" || protocol == "all" {
		if fromPort != 0 && fromPort != -1 {
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf("Security group '%s' %s specifies FromPort %d for protocol -1 (all); ports are ignored for this protocol", resName, context, fromPort),
				Path:    []string{"Resources", resName, "Properties", "FromPort"},
			})
		}
		if toPort != 0 && toPort != -1 {
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf("Security group '%s' %s specifies ToPort %d for protocol -1 (all); ports are ignored for this protocol", resName, context, toPort),
				Path:    []string{"Resources", resName, "Properties", "ToPort"},
			})
		}
	}

	// For ICMP, FromPort is ICMP type and ToPort is ICMP code
	if protocol == "icmp" {
		// Valid ICMP types are -1 (all) or 0-255
		if fromPort < -1 || fromPort > 255 {
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf("Security group '%s' %s has invalid ICMP type %d (FromPort); valid range is -1 to 255", resName, context, fromPort),
				Path:    []string{"Resources", resName, "Properties", "FromPort"},
			})
		}
		if toPort < -1 || toPort > 255 {
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf("Security group '%s' %s has invalid ICMP code %d (ToPort); valid range is -1 to 255", resName, context, toPort),
				Path:    []string{"Resources", resName, "Properties", "ToPort"},
			})
		}
	}
}

func (r *W3687) getProtocol(v any) string {
	switch p := v.(type) {
	case string:
		return p
	case int:
		return fmt.Sprintf("%d", p)
	case float64:
		return fmt.Sprintf("%d", int(p))
	}
	return ""
}

func (r *W3687) getPort(v any) int {
	switch p := v.(type) {
	case int:
		return p
	case float64:
		return int(p)
	case string:
		var port int
		_, _ = fmt.Sscanf(p, "%d", &port)
		return port
	}
	return 0
}
