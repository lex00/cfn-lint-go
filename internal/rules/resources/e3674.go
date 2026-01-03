// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3674{})
}

// E3674 validates Primary and PrivateIpAddress conflict.
type E3674 struct{}

func (r *E3674) ID() string { return "E3674" }

func (r *E3674) ShortDesc() string {
	return "Validate Primary and PrivateIpAddress conflict"
}

func (r *E3674) Description() string {
	return "Validates that AWS::EC2::NetworkInterface PrivateIpAddresses entries do not specify both Primary and PrivateIpAddress."
}

func (r *E3674) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3674"
}

func (r *E3674) Tags() []string {
	return []string{"resources", "properties", "ec2", "networkinterface"}
}

func (r *E3674) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::EC2::NetworkInterface" {
			continue
		}

		privateIpAddresses, hasPrivateIpAddresses := res.Properties["PrivateIpAddresses"]
		if !hasPrivateIpAddresses || isIntrinsicFunction(privateIpAddresses) {
			continue
		}

		ipList, ok := privateIpAddresses.([]any)
		if !ok {
			continue
		}

		for _, ip := range ipList {
			ipMap, ok := ip.(map[string]any)
			if !ok {
				continue
			}

			_, hasPrimary := ipMap["Primary"]
			_, hasPrivateIpAddress := ipMap["PrivateIpAddress"]

			// Both Primary and PrivateIpAddress cannot be specified
			if hasPrimary && hasPrivateIpAddress {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': NetworkInterface PrivateIpAddresses entry cannot specify both Primary and PrivateIpAddress",
						resName,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", "PrivateIpAddresses"},
				})
			}
		}
	}

	return matches
}
