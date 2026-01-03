// Package formats contains format validation rules (E11xx).
package formats

import (
	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1103{})
}

// E1103 is a parent rule for validating format keyword in schemas.
type E1103 struct{}

func (r *E1103) ID() string { return "E1103" }

func (r *E1103) ShortDesc() string {
	return "Value format validation"
}

func (r *E1103) Description() string {
	return "Parent rule for validating the format keyword in schemas."
}

func (r *E1103) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-resource-specification.html"
}

func (r *E1103) Tags() []string {
	return []string{"format", "validation"}
}

func (r *E1103) Match(tmpl *template.Template) []rules.Match {
	// E1103 is a parent rule - actual validation is done by child rules
	// E1150-E1156, etc.
	return nil
}
