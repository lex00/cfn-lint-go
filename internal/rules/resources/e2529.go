package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E2529{})
}

// E2529 validates that LogGroups do not exceed the limit of 2 SubscriptionFilters.
type E2529 struct{}

func (r *E2529) ID() string { return "E2529" }

func (r *E2529) ShortDesc() string {
	return "Validate SubscriptionFilters limit per LogGroup"
}

func (r *E2529) Description() string {
	return "CloudWatch Logs LogGroups can have a maximum of 2 SubscriptionFilters. This rule validates that the SubscriptionFilters property does not exceed this limit."
}

func (r *E2529) Source() string {
	return "https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/cloudwatch_limits_cwl.html"
}

func (r *E2529) Tags() []string {
	return []string{"resources", "logs", "cloudwatch", "limits"}
}

const maxSubscriptionFiltersPerLogGroup = 2

func (r *E2529) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Logs::LogGroup" {
			continue
		}

		// Get the SubscriptionFilters property
		subscriptionFilters, ok := res.Properties["SubscriptionFilters"]
		if !ok {
			continue
		}

		// Check if it's an array
		filters, ok := subscriptionFilters.([]any)
		if !ok {
			// If it's an intrinsic function, we can't validate the count
			continue
		}

		filterCount := len(filters)
		if filterCount > maxSubscriptionFiltersPerLogGroup {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("LogGroup '%s' has %d SubscriptionFilters, exceeding the limit of %d", resName, filterCount, maxSubscriptionFiltersPerLogGroup),
				Line:    res.Node.Line,
				Column:  res.Node.Column,
				Path:    []string{"Resources", resName, "Properties", "SubscriptionFilters"},
			})
		}
	}

	return matches
}
