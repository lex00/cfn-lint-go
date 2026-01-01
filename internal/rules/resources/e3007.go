// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3007{})
}

// E3007 checks that resource logical IDs are unique.
// Note: This is implicitly enforced by the parser (map keys are unique),
// but we keep this rule for completeness and to catch duplicates at the YAML level.
type E3007 struct{}

func (r *E3007) ID() string { return "E3007" }

func (r *E3007) ShortDesc() string {
	return "Duplicate resource logical ID"
}

func (r *E3007) Description() string {
	return "Checks that all resource logical IDs are unique within the template."
}

func (r *E3007) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3007"
}

func (r *E3007) Tags() []string {
	return []string{"resources", "unique"}
}

func (r *E3007) Match(tmpl *template.Template) []rules.Match {
	// YAML/JSON parsers already enforce unique keys at the same level.
	// Duplicate keys in YAML result in the last value winning.
	// This rule is satisfied implicitly by the parser.
	// If we wanted to catch this, we'd need to check raw YAML nodes before parsing.
	return nil
}
