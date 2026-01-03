// Package formats contains format validation rules (E11xx).
package formats

import (
	"fmt"
	"regexp"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1154{})
}

// E1154 validates VPC subnet ID format.
type E1154 struct{}

func (r *E1154) ID() string { return "E1154" }

func (r *E1154) ShortDesc() string {
	return "VPC subnet id format validation"
}

func (r *E1154) Description() string {
	return "Validates that subnet IDs match the format subnet-[0-9a-f]{8,17}."
}

func (r *E1154) Source() string {
	return "https://docs.aws.amazon.com/vpc/latest/userguide/configure-subnets.html"
}

func (r *E1154) Tags() []string {
	return []string{"format", "subnet", "vpc"}
}

var subnetIDPattern = regexp.MustCompile(`^subnet-[0-9a-f]{8,17}$`)

func (r *E1154) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		subnetRefs := findFormatReferences(res.Properties, "subnet-", []string{"SubnetId", "SubnetIds"})
		for _, ref := range subnetRefs {
			if ref.value != "" && !subnetIDPattern.MatchString(ref.value) && !isIntrinsicFunction(ref.rawValue) {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Invalid subnet ID format '%s' in resource '%s', expected format: subnet-[0-9a-f]{8,17}", ref.value, resName),
					Path:    append([]string{"Resources", resName, "Properties"}, ref.path...),
				})
			}
		}
	}

	return matches
}
