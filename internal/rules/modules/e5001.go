// Package modules contains module validation rules (E5xxx).
package modules

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E5001{})
}

// E5001 checks that CloudFormation Modules resources are valid.
// CloudFormation Modules are reusable resource configurations that follow
// the naming pattern Organization::Service::Resource::MODULE.
type E5001 struct{}

func (r *E5001) ID() string { return "E5001" }

func (r *E5001) ShortDesc() string {
	return "Modules resource validation"
}

func (r *E5001) Description() string {
	return "Checks that CloudFormation Modules resources are valid. " +
		"Modules must follow the format 'Organization::Service::Resource::MODULE' and have valid Properties."
}

func (r *E5001) Source() string {
	return "https://docs.aws.amazon.com/cloudformation-cli/latest/userguide/modules.html"
}

func (r *E5001) Tags() []string {
	return []string{"resources", "modules"}
}

func (r *E5001) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for name, res := range tmpl.Resources {
		// Check if this is a module resource (ends with ::MODULE)
		if !strings.HasSuffix(res.Type, "::MODULE") {
			continue
		}

		// Validate module type format: Organization::Service::Resource::MODULE
		parts := strings.Split(res.Type, "::")
		if len(parts) != 4 {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Module resource '%s' has invalid type format '%s'. Expected format: 'Organization::Service::Resource::MODULE'", name, res.Type),
				Line:    res.Node.Line,
				Column:  res.Node.Column,
				Path:    []string{"Resources", name, "Type"},
			})
			continue
		}

		// Validate that each part is not empty
		for i, part := range parts {
			if part == "" {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Module resource '%s' has empty component in type '%s'", name, res.Type),
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", name, "Type"},
				})
				break
			}
			// Last part must be MODULE
			if i == 3 && part != "MODULE" {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Module resource '%s' type must end with '::MODULE', got '%s'", name, res.Type),
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", name, "Type"},
				})
			}
		}

		// Modules must have Properties
		if res.Properties == nil || len(res.Properties) == 0 {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Module resource '%s' must have Properties defined", name),
				Line:    res.Node.Line,
				Column:  res.Node.Column,
				Path:    []string{"Resources", name},
			})
		}
	}

	return matches
}
