// Package errors contains base error rules (E0xxx).
package errors

import (
	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E0200{})
}

// E0200 checks parameter file syntax validation.
// Note: This rule is registered but currently returns no matches as parameter file
// validation requires separate file parsing infrastructure not yet implemented.
type E0200 struct{}

func (r *E0200) ID() string { return "E0200" }

func (r *E0200) ShortDesc() string {
	return "Parameter file syntax validation"
}

func (r *E0200) Description() string {
	return "Validate if a parameter file has the correct syntax for one of the supported formats. " +
		"Parameter files are JSON files used to provide parameter values during CloudFormation stack deployment."
}

func (r *E0200) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md"
}

func (r *E0200) Tags() []string {
	return []string{"base", "parameters"}
}

func (r *E0200) Match(tmpl *template.Template) []rules.Match {
	// Parameter file validation requires separate file parsing infrastructure.
	// This rule is registered for compatibility but does not perform validation
	// on CloudFormation templates directly. It would need to be invoked with
	// parameter file-specific parsers.
	return nil
}
