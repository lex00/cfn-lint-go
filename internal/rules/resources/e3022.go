package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3022{})
}

// E3022 validates SubnetRouteTableAssociation has only one association per subnet.
type E3022 struct{}

func (r *E3022) ID() string {
	return "E3022"
}

func (r *E3022) ShortDesc() string {
	return "Resource SubnetRouteTableAssociation Properties"
}

func (r *E3022) Description() string {
	return "Confirms only one association per subnet exists"
}

func (r *E3022) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-subnetroutetableassociation.html"
}

func (r *E3022) Tags() []string {
	return []string{"resources", "ec2", "subnet"}
}

func (r *E3022) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Track SubnetId to resource name mappings
	subnetAssociations := make(map[string]string) // SubnetId -> ResourceName

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::EC2::SubnetRouteTableAssociation" {
			continue
		}

		subnetID, ok := res.Properties["SubnetId"]
		if !ok {
			continue
		}

		// Skip intrinsic functions
		if isIntrinsicFunction(subnetID) {
			continue
		}

		subnetStr, ok := subnetID.(string)
		if !ok {
			continue
		}

		if existingRes, exists := subnetAssociations[subnetStr]; exists {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Subnet '%s' has multiple route table associations ('%s' and '%s')", subnetStr, existingRes, resName),
				Line:    res.Node.Line,
				Column:  res.Node.Column,
				Path:    []string{"Resources", resName, "Properties", "SubnetId"},
			})
		} else {
			subnetAssociations[subnetStr] = resName
		}
	}

	return matches
}
