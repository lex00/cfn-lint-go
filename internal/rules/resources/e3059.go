// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"net"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3059{})
}

// E3059 validates that subnet CIDRs are within VPC CIDR blocks.
type E3059 struct{}

func (r *E3059) ID() string { return "E3059" }

func (r *E3059) ShortDesc() string {
	return "Subnet CIDRs within VPC"
}

func (r *E3059) Description() string {
	return "Validates that subnet CIDR blocks fall within the parent VPC CIDR blocks."
}

func (r *E3059) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3059"
}

func (r *E3059) Tags() []string {
	return []string{"resources", "properties", "vpc", "subnet", "cidr"}
}

func (r *E3059) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Collect VPC CIDRs
	vpcCIDRs := make(map[string][]string)

	for resName, res := range tmpl.Resources {
		if res.Type == "AWS::EC2::VPC" {
			if cidrBlock, hasCIDR := res.Properties["CidrBlock"]; hasCIDR {
				if cidrStr, ok := cidrBlock.(string); ok {
					vpcCIDRs[resName] = []string{cidrStr}
				}
			}
		}
	}

	// Check subnet CIDRs
	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::EC2::Subnet" {
			continue
		}

		subnetCIDR, hasCIDR := res.Properties["CidrBlock"]
		vpcID, hasVPC := res.Properties["VpcId"]

		if !hasCIDR || !hasVPC {
			continue
		}

		subnetCIDRStr, ok1 := subnetCIDR.(string)
		if !ok1 {
			continue
		}

		// Parse subnet CIDR
		_, subnetNet, err := net.ParseCIDR(subnetCIDRStr)
		if err != nil {
			continue
		}

		// Get VPC reference
		vpcRef := ""
		if vpcMap, ok := vpcID.(map[string]interface{}); ok {
			if ref, hasRef := vpcMap["Ref"]; hasRef {
				if refStr, ok := ref.(string); ok {
					vpcRef = refStr
				}
			}
		}

		if vpcRef == "" {
			continue
		}

		// Get VPC CIDRs
		vpcCIDRList, hasVPCCIDR := vpcCIDRs[vpcRef]
		if !hasVPCCIDR {
			continue
		}

		// Check if subnet CIDR is within any VPC CIDR
		isWithinVPC := false
		for _, vpcCIDRStr := range vpcCIDRList {
			_, vpcNet, err := net.ParseCIDR(vpcCIDRStr)
			if err != nil {
				continue
			}

			// Check if subnet is within VPC
			if r.isSubnet(subnetNet, vpcNet) {
				isWithinVPC = true
				break
			}
		}

		if !isWithinVPC {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Subnet CIDR '%s' is not within VPC '%s' CIDR blocks %v",
					resName, subnetCIDRStr, vpcRef, vpcCIDRList,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties", "CidrBlock"},
			})
		}
	}

	return matches
}

// isSubnet checks if subnet is within vpc
func (r *E3059) isSubnet(subnet, vpc *net.IPNet) bool {
	// Check if subnet's network address is within VPC
	if !vpc.Contains(subnet.IP) {
		return false
	}

	// Check if subnet's broadcast address is within VPC
	subnetBroadcast := r.getBroadcastAddr(subnet)
	if !vpc.Contains(subnetBroadcast) {
		return false
	}

	return true
}

// getBroadcastAddr calculates the broadcast address for a network
func (r *E3059) getBroadcastAddr(n *net.IPNet) net.IP {
	ip := n.IP.To4()
	if ip == nil {
		ip = n.IP.To16()
	}
	broadcast := make(net.IP, len(ip))
	for i := range ip {
		broadcast[i] = ip[i] | ^n.Mask[i]
	}
	return broadcast
}
