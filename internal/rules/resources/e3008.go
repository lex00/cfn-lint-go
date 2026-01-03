package resources

import (
	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3008{})
}

// E3008 validates arrays in order for schema validation (prefixItems).
type E3008 struct{}

func (r *E3008) ID() string {
	return "E3008"
}

func (r *E3008) ShortDesc() string {
	return "Validate an array in order"
}

func (r *E3008) Description() string {
	return "Will validate arrays in order for schema validation"
}

func (r *E3008) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/cfn-schema-specification.md#prefixitems"
}

func (r *E3008) Tags() []string {
	return []string{"resources", "properties", "array", "prefixItems"}
}

func (r *E3008) Match(tmpl *template.Template) []rules.Match {
	// This rule validates array ordering based on schema prefixItems.
	// The schema package doesn't currently expose prefixItems validation,
	// so this is a placeholder implementation.
	// TODO: Implement when schema package supports prefixItems.
	return nil
}
