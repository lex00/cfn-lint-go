// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3670{})
}

// E3670 validates AmazonMQ Broker instance types.
type E3670 struct{}

func (r *E3670) ID() string { return "E3670" }

func (r *E3670) ShortDesc() string {
	return "Validate AmazonMQ Broker instance types"
}

func (r *E3670) Description() string {
	return "Validates that AWS::AmazonMQ::Broker resources specify valid instance types."
}

func (r *E3670) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3670"
}

func (r *E3670) Tags() []string {
	return []string{"resources", "properties", "amazonmq", "instancetype"}
}

func (r *E3670) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// AmazonMQ instance type families
	validFamilies := []string{
		"mq.t2.", "mq.t3.",
		"mq.m4.", "mq.m5.",
	}

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::AmazonMQ::Broker" {
			continue
		}

		hostInstanceType, hasHostInstanceType := res.Properties["HostInstanceType"]
		if !hasHostInstanceType || isIntrinsicFunction(hostInstanceType) {
			continue
		}

		hostInstanceTypeStr, ok := hostInstanceType.(string)
		if !ok {
			continue
		}

		// Check if instance type starts with a valid family
		isValid := false
		for _, family := range validFamilies {
			if strings.HasPrefix(hostInstanceTypeStr, family) {
				isValid = true
				break
			}
		}

		if !isValid {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Invalid AmazonMQ instance type '%s'. Must start with 'mq.'",
					resName, hostInstanceTypeStr,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties", "HostInstanceType"},
			})
		}
	}

	return matches
}
