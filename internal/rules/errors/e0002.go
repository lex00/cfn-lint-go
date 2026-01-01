// Package errors contains template structure validation rules (E0xxx).
package errors

import (
	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E0002{})
}

// E0002 checks for rule processing errors.
// This rule is triggered when a linting rule encounters an internal error.
type E0002 struct{}

func (r *E0002) ID() string { return "E0002" }

func (r *E0002) ShortDesc() string {
	return "Rule processing error"
}

func (r *E0002) Description() string {
	return "Checks for errors during rule execution. This indicates a bug in cfn-lint-go."
}

func (r *E0002) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E0002"
}

func (r *E0002) Tags() []string {
	return []string{"base", "internal"}
}

func (r *E0002) Match(tmpl *template.Template) []rules.Match {
	// Rule processing errors are caught at runtime by the linter.
	// This rule exists for documentation purposes.
	return nil
}
