// Package errors contains base error rules (E0xxx).
package errors

import (
	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E0003{})
}

// E0003 checks for configuration errors in the template.
type E0003 struct{}

func (r *E0003) ID() string { return "E0003" }

func (r *E0003) ShortDesc() string {
	return "Configuration error"
}

func (r *E0003) Description() string {
	return "Checks for configuration errors in the CloudFormation template structure."
}

func (r *E0003) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-anatomy.html"
}

func (r *E0003) Tags() []string {
	return []string{"base", "configuration"}
}

func (r *E0003) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check for valid AWSTemplateFormatVersion
	if tmpl.AWSTemplateFormatVersion != "" && tmpl.AWSTemplateFormatVersion != "2010-09-09" {
		matches = append(matches, rules.Match{
			Message: "AWSTemplateFormatVersion must be '2010-09-09'",
			Path:    []string{"AWSTemplateFormatVersion"},
		})
	}

	// Check that template has at least Resources section
	if len(tmpl.Resources) == 0 {
		matches = append(matches, rules.Match{
			Message: "Template must have at least one resource in the Resources section",
			Path:    []string{"Resources"},
		})
	}

	return matches
}
