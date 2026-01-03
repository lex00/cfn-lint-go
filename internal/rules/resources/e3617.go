// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3617{})
}

// E3617 validates ManagedBlockchain instance types.
type E3617 struct{}

func (r *E3617) ID() string { return "E3617" }

func (r *E3617) ShortDesc() string {
	return "Validate ManagedBlockchain instance type"
}

func (r *E3617) Description() string {
	return "Validates that AWS::ManagedBlockchain::Node resources specify valid instance types."
}

func (r *E3617) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3617"
}

func (r *E3617) Tags() []string {
	return []string{"resources", "properties", "managedblockchain", "instancetype"}
}

func (r *E3617) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	validInstanceTypes := map[string]bool{
		"bc.t3.small":  true,
		"bc.t3.medium": true,
		"bc.t3.large":  true,
		"bc.t3.xlarge": true,
		"bc.m5.large":  true,
		"bc.m5.xlarge": true,
		"bc.m5.2xlarge": true,
		"bc.m5.4xlarge": true,
		"bc.c5.large":  true,
		"bc.c5.xlarge": true,
		"bc.c5.2xlarge": true,
		"bc.c5.4xlarge": true,
	}

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::ManagedBlockchain::Node" {
			continue
		}

		nodeConfig, hasNodeConfig := res.Properties["NodeConfiguration"]
		if !hasNodeConfig || isIntrinsicFunction(nodeConfig) {
			continue
		}

		nodeConfigMap, ok := nodeConfig.(map[string]any)
		if !ok {
			continue
		}

		instanceType, hasInstanceType := nodeConfigMap["InstanceType"]
		if !hasInstanceType || isIntrinsicFunction(instanceType) {
			continue
		}

		instanceTypeStr, ok := instanceType.(string)
		if !ok {
			continue
		}

		if !validInstanceTypes[instanceTypeStr] {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Invalid ManagedBlockchain instance type '%s'. Must be a valid bc.* instance type",
					resName, instanceTypeStr,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties", "NodeConfiguration", "InstanceType"},
			})
		}
	}

	return matches
}
