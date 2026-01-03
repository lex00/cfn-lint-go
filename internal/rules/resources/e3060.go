// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"net"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3060{})
}

// E3060 validates that subnet CIDRs do not overlap.
type E3060 struct{}

func (r *E3060) ID() string { return "E3060" }

func (r *E3060) ShortDesc() string {
	return "Subnet CIDRs no overlap"
}

func (r *E3060) Description() string {
	return "Validates that subnet CIDR blocks within the same VPC do not overlap."
}

func (r *E3060) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3060"
}

func (r *E3060) Tags() []string {
	return []string{"resources", "properties", "vpc", "subnet", "cidr"}
}

func (r *E3060) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Group subnets by VPC
	subnetsByVPC := make(map[string][]subnetInfo)

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

		// Get VPC reference
		vpcRef := ""
		if vpcMap, ok := vpcID.(map[string]interface{}); ok {
			if ref, hasRef := vpcMap["Ref"]; hasRef {
				if refStr, ok := ref.(string); ok {
					vpcRef = refStr
				}
			}
		} else if vpcStr, ok := vpcID.(string); ok {
			vpcRef = vpcStr
		}

		if vpcRef == "" {
			continue
		}

		// Parse CIDR
		_, cidrNet, err := net.ParseCIDR(subnetCIDRStr)
		if err != nil {
			continue
		}

		subnetsByVPC[vpcRef] = append(subnetsByVPC[vpcRef], subnetInfo{
			name:   resName,
			cidr:   subnetCIDRStr,
			cidrNet: cidrNet,
			line:   res.Node.Line,
			column: res.Node.Column,
		})
	}

	// Check for overlaps within each VPC
	for vpcRef, subnets := range subnetsByVPC {
		for i := 0; i < len(subnets); i++ {
			for j := i + 1; j < len(subnets); j++ {
				if r.cidrsOverlap(subnets[i].cidrNet, subnets[j].cidrNet) {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf(
							"Subnet '%s' CIDR '%s' overlaps with subnet '%s' CIDR '%s' in VPC '%s'",
							subnets[i].name, subnets[i].cidr, subnets[j].name, subnets[j].cidr, vpcRef,
						),
						Line:   subnets[i].line,
						Column: subnets[i].column,
						Path:   []string{"Resources", subnets[i].name, "Properties", "CidrBlock"},
					})
				}
			}
		}
	}

	return matches
}

type subnetInfo struct {
	name    string
	cidr    string
	cidrNet *net.IPNet
	line    int
	column  int
}

// cidrsOverlap checks if two CIDR blocks overlap
func (r *E3060) cidrsOverlap(cidr1, cidr2 *net.IPNet) bool {
	return cidr1.Contains(cidr2.IP) || cidr2.Contains(cidr1.IP)
}
