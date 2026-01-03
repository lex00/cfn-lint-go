// Package errors contains base error rules (E0xxx and E2900).
package errors

import (
	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E2900{})
}

// E2900 validates deployment file parameters.
// Note: This rule is registered but currently returns no matches as deployment file
// validation requires separate file parsing infrastructure not yet implemented.
type E2900 struct{}

func (r *E2900) ID() string { return "E2900" }

func (r *E2900) ShortDesc() string {
	return "Deployment file parameters validation"
}

func (r *E2900) Description() string {
	return "Validate that parameters defined in deployment files (e.g., AWS SAM deploy configuration) " +
		"match the parameters defined in the CloudFormation template. This ensures that all required " +
		"parameters are provided and no extra parameters are specified."
}

func (r *E2900) Source() string {
	return "https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-deploying.html"
}

func (r *E2900) Tags() []string {
	return []string{"parameters", "deployment"}
}

func (r *E2900) Match(tmpl *template.Template) []rules.Match {
	// Deployment file parameter validation requires separate file parsing infrastructure.
	// This rule is registered for compatibility but does not perform validation
	// on CloudFormation templates directly. It would need to be invoked with
	// deployment configuration file parsers (e.g., samconfig.toml, parameter overrides).
	//
	// Future implementation would:
	// 1. Parse deployment configuration files (samconfig.toml, parameter JSON files)
	// 2. Extract parameter overrides from deployment configuration
	// 3. Validate against template.Parameters map
	// 4. Check for missing required parameters (those without defaults)
	// 5. Check for extra parameters not defined in template
	//
	// Example validation logic:
	// - Ensure all required template parameters are provided
	// - Ensure no undefined parameters are specified
	// - Validate parameter value types match template constraints
	return nil
}
