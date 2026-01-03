// Package formats contains format validation rules (E11xx).
package formats

import (
	"fmt"
	"regexp"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1151{})
}

// E1151 validates VPC ID format.
type E1151 struct{}

func (r *E1151) ID() string { return "E1151" }

func (r *E1151) ShortDesc() string {
	return "VPC id format validation"
}

func (r *E1151) Description() string {
	return "Validates that VPC IDs match the format vpc-[0-9a-f]{8,17}."
}

func (r *E1151) Source() string {
	return "https://docs.aws.amazon.com/vpc/latest/userguide/what-is-amazon-vpc.html"
}

func (r *E1151) Tags() []string {
	return []string{"format", "vpc"}
}

var vpcIDPattern = regexp.MustCompile(`^vpc-[0-9a-f]{8,17}$`)

func (r *E1151) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		vpcRefs := findFormatReferences(res.Properties, "vpc-", []string{"VpcId", "VPCId"})
		for _, ref := range vpcRefs {
			if ref.value != "" && !vpcIDPattern.MatchString(ref.value) && !isIntrinsicFunction(ref.rawValue) {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Invalid VPC ID format '%s' in resource '%s', expected format: vpc-[0-9a-f]{8,17}", ref.value, resName),
					Path:    append([]string{"Resources", resName, "Properties"}, ref.path...),
				})
			}
		}
	}

	return matches
}
