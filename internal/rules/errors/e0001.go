// Package errors contains template structure validation rules (E0xxx).
package errors

import (
	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E0001{})
}

// E0001 checks for template transformation errors.
// This rule is triggered when AWS::Include or other transforms fail.
type E0001 struct{}

func (r *E0001) ID() string { return "E0001" }

func (r *E0001) ShortDesc() string {
	return "Template transformation error"
}

func (r *E0001) Description() string {
	return "Checks for errors during template transformation (e.g., AWS::Include, AWS::Serverless)."
}

func (r *E0001) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E0001"
}

func (r *E0001) Tags() []string {
	return []string{"base", "transform"}
}

func (r *E0001) Match(tmpl *template.Template) []rules.Match {
	// Transform errors are handled at parse/transform time.
	// This rule exists for documentation and future transform validation.
	return nil
}
