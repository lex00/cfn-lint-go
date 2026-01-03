// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3639{})
}

// E3639 validates DynamoDB ProvisionedThroughput requirement.
type E3639 struct{}

func (r *E3639) ID() string { return "E3639" }

func (r *E3639) ShortDesc() string {
	return "DynamoDB PROVISIONED requires ProvisionedThroughput"
}

func (r *E3639) Description() string {
	return "Validates that DynamoDB tables with BillingMode 'PROVISIONED' or no BillingMode (default) specify ProvisionedThroughput property."
}

func (r *E3639) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3639"
}

func (r *E3639) Tags() []string {
	return []string{"resources", "properties", "dynamodb", "billing"}
}

func (r *E3639) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::DynamoDB::Table" {
			continue
		}

		billingMode, hasBillingMode := res.Properties["BillingMode"]
		requiresProvisionedThroughput := false

		if !hasBillingMode {
			// Default billing mode is PROVISIONED
			requiresProvisionedThroughput = true
		} else if !isIntrinsicFunction(billingMode) {
			if billingModeStr, ok := billingMode.(string); ok && billingModeStr == "PROVISIONED" {
				requiresProvisionedThroughput = true
			}
		}

		if requiresProvisionedThroughput {
			if _, hasProvisionedThroughput := res.Properties["ProvisionedThroughput"]; !hasProvisionedThroughput {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': DynamoDB Table with BillingMode 'PROVISIONED' must specify ProvisionedThroughput",
						resName,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	return matches
}
