package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W3663{})
}

// W3663 warns when Lambda Permission is missing SourceAccount for S3 or SNS.
type W3663 struct{}

func (r *W3663) ID() string { return "W3663" }

func (r *W3663) ShortDesc() string {
	return "SourceAccount required"
}

func (r *W3663) Description() string {
	return "Warns when Lambda Permission resources for S3 or SNS sources don't include SourceAccount, which could allow cross-account invocations."
}

func (r *W3663) Source() string {
	return "https://docs.aws.amazon.com/lambda/latest/dg/access-control-resource-based.html"
}

func (r *W3663) Tags() []string {
	return []string{"warnings", "lambda", "permissions", "security"}
}

// Services that should have SourceAccount
var servicesNeedingSourceAccount = map[string]bool{
	"s3.amazonaws.com":  true,
	"sns.amazonaws.com": true,
}

func (r *W3663) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Lambda::Permission" {
			continue
		}

		principal, hasPrincipal := res.Properties["Principal"].(string)
		if !hasPrincipal {
			continue
		}

		// Check if this is a service that needs SourceAccount
		if !servicesNeedingSourceAccount[principal] {
			continue
		}

		// Check for SourceAccount
		_, hasSourceAccount := res.Properties["SourceAccount"]
		_, hasSourceArn := res.Properties["SourceArn"]

		if !hasSourceAccount && hasSourceArn {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Lambda Permission '%s' for '%s' has SourceArn but no SourceAccount; consider adding SourceAccount to prevent cross-account invocations", resName, principal),
				Path:    []string{"Resources", resName, "Properties"},
			})
		}

		if !hasSourceAccount && !hasSourceArn {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Lambda Permission '%s' for '%s' has neither SourceAccount nor SourceArn; this allows any %s resource to invoke the function", resName, principal, principal),
				Path:    []string{"Resources", resName, "Properties"},
			})
		}
	}

	return matches
}
