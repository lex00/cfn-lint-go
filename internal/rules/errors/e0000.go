// Package errors contains template structure validation rules (E0xxx).
package errors

import (
	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E0000{})
}

// E0000 checks for template parse errors.
// This rule is implicitly triggered by the parser; this struct exists
// for documentation and rule listing purposes.
type E0000 struct{}

func (r *E0000) ID() string { return "E0000" }

func (r *E0000) ShortDesc() string {
	return "Template parse error"
}

func (r *E0000) Description() string {
	return "Checks that the template is valid YAML or JSON and can be parsed."
}

func (r *E0000) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E0000"
}

func (r *E0000) Tags() []string {
	return []string{"base", "template"}
}

func (r *E0000) Match(tmpl *template.Template) []rules.Match {
	// Parse errors are already caught during template parsing.
	// This rule exists for documentation purposes.
	return nil
}
