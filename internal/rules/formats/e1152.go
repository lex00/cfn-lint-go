// Package formats contains format validation rules (E11xx).
package formats

import (
	"fmt"
	"regexp"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1152{})
}

// E1152 validates AMI ID format.
type E1152 struct{}

func (r *E1152) ID() string { return "E1152" }

func (r *E1152) ShortDesc() string {
	return "AMI id format validation"
}

func (r *E1152) Description() string {
	return "Validates that AMI IDs match the format ami-[0-9a-f]{8,17}."
}

func (r *E1152) Source() string {
	return "https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html"
}

func (r *E1152) Tags() []string {
	return []string{"format", "ami"}
}

var amiIDPattern = regexp.MustCompile(`^ami-[0-9a-f]{8,17}$`)

func (r *E1152) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		amiRefs := findFormatReferences(res.Properties, "ami-", []string{"ImageId"})
		for _, ref := range amiRefs {
			if ref.value != "" && !amiIDPattern.MatchString(ref.value) && !isIntrinsicFunction(ref.rawValue) {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Invalid AMI ID format '%s' in resource '%s', expected format: ami-[0-9a-f]{8,17}", ref.value, resName),
					Path:    append([]string{"Resources", resName, "Properties"}, ref.path...),
				})
			}
		}
	}

	return matches
}
