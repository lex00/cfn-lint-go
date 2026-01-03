// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3638{})
}

// E3638 validates DynamoDB BillingMode pay-per-request configuration.
type E3638 struct{}

func (r *E3638) ID() string { return "E3638" }

func (r *E3638) ShortDesc() string {
	return "DynamoDB PAY_PER_REQUEST excludes ProvisionedThroughput"
}

func (r *E3638) Description() string {
	return "Validates that DynamoDB tables with BillingMode 'PAY_PER_REQUEST' do not specify ProvisionedThroughput property."
}

func (r *E3638) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3638"
}

func (r *E3638) Tags() []string {
	return []string{"resources", "properties", "dynamodb", "billing"}
}

func (r *E3638) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::DynamoDB::Table" {
			continue
		}

		billingMode, hasBillingMode := res.Properties["BillingMode"]
		if !hasBillingMode || isIntrinsicFunction(billingMode) {
			continue
		}

		billingModeStr, ok := billingMode.(string)
		if !ok {
			continue
		}

		if billingModeStr == "PAY_PER_REQUEST" {
			if _, hasProvisionedThroughput := res.Properties["ProvisionedThroughput"]; hasProvisionedThroughput {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': DynamoDB Table with BillingMode 'PAY_PER_REQUEST' must not specify ProvisionedThroughput",
						resName,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", "ProvisionedThroughput"},
				})
			}
		}
	}

	return matches
}
