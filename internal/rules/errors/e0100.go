// Package errors contains base error rules (E0xxx).
package errors

import (
	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E0100{})
}

// E0100 checks deployment file syntax validation.
// Note: This rule is registered but currently returns no matches as deployment file
// validation requires separate file parsing infrastructure not yet implemented.
type E0100 struct{}

func (r *E0100) ID() string { return "E0100" }

func (r *E0100) ShortDesc() string {
	return "Deployment file syntax validation"
}

func (r *E0100) Description() string {
	return "Validate if a deployment file has the correct syntax for one of the supported formats. " +
		"Deployment files are used to configure template deployment parameters and are separate from CloudFormation templates."
}

func (r *E0100) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md"
}

func (r *E0100) Tags() []string {
	return []string{"base", "deployment"}
}

func (r *E0100) Match(tmpl *template.Template) []rules.Match {
	// Deployment file validation requires separate file parsing infrastructure.
	// This rule is registered for compatibility but does not perform validation
	// on CloudFormation templates directly. It would need to be invoked with
	// deployment-specific file parsers.
	return nil
}
