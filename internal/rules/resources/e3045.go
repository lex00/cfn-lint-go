// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3045{})
}

// E3045 validates that S3 buckets using AccessControl also configure OwnershipControls.
type E3045 struct{}

func (r *E3045) ID() string { return "E3045" }

func (r *E3045) ShortDesc() string {
	return "S3 AccessControl with OwnershipControls"
}

func (r *E3045) Description() string {
	return "Validates that S3 buckets using AccessControl (ACLs) have explicit OwnershipControls configuration."
}

func (r *E3045) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3045"
}

func (r *E3045) Tags() []string {
	return []string{"resources", "properties", "s3", "bucket"}
}

func (r *E3045) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::S3::Bucket" {
			continue
		}

		_, hasAccessControl := res.Properties["AccessControl"]
		_, hasOwnershipControls := res.Properties["OwnershipControls"]

		// If AccessControl is specified but OwnershipControls is not, warn
		if hasAccessControl && !hasOwnershipControls {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': S3 buckets using AccessControl should explicitly configure OwnershipControls to avoid ACL-related issues",
					resName,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties"},
			})
		}
	}

	return matches
}
