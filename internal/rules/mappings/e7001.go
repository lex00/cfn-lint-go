// Package mappings contains mapping validation rules (E7xxx).
package mappings

import (
	"fmt"
	"regexp"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E7001{})
}

// E7001 validates that Mappings are properly configured.
type E7001 struct{}

func (r *E7001) ID() string { return "E7001" }

func (r *E7001) ShortDesc() string {
	return "Mappings are appropriately configured"
}

func (r *E7001) Description() string {
	return "Validates that Mappings have proper structure with valid keys and values."
}

func (r *E7001) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E7001"
}

func (r *E7001) Tags() []string {
	return []string{"mappings", "configuration"}
}

// alphanumericPattern matches valid mapping key names.
var alphanumericPattern = regexp.MustCompile(`^[a-zA-Z0-9.-]+$`)

func (r *E7001) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for name, mapping := range tmpl.Mappings {
		// Validate mapping name
		if len(name) > 255 {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Mapping '%s' name exceeds 255 character limit", name),
				Line:    mapping.Node.Line,
				Column:  mapping.Node.Column,
				Path:    []string{"Mappings", name},
			})
		}

		if !alphanumericPattern.MatchString(name) {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Mapping '%s' name must be alphanumeric (a-zA-Z0-9.-)", name),
				Line:    mapping.Node.Line,
				Column:  mapping.Node.Column,
				Path:    []string{"Mappings", name},
			})
		}

		// Validate mapping has at least one top-level key
		if len(mapping.Values) == 0 {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Mapping '%s' must have at least one top-level key", name),
				Line:    mapping.Node.Line,
				Column:  mapping.Node.Column,
				Path:    []string{"Mappings", name},
			})
		}

		// Validate top-level and second-level keys
		for topKey, secondLevel := range mapping.Values {
			if !alphanumericPattern.MatchString(topKey) {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Mapping '%s' key '%s' must be alphanumeric (a-zA-Z0-9.-)", name, topKey),
					Line:    mapping.Node.Line,
					Column:  mapping.Node.Column,
					Path:    []string{"Mappings", name, topKey},
				})
			}

			if len(secondLevel) == 0 {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Mapping '%s' top-level key '%s' must have at least one second-level key", name, topKey),
					Line:    mapping.Node.Line,
					Column:  mapping.Node.Column,
					Path:    []string{"Mappings", name, topKey},
				})
			}

			for secondKey := range secondLevel {
				if !alphanumericPattern.MatchString(secondKey) {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Mapping '%s' second-level key '%s' must be alphanumeric (a-zA-Z0-9.-)", name, secondKey),
						Line:    mapping.Node.Line,
						Column:  mapping.Node.Column,
						Path:    []string{"Mappings", name, topKey, secondKey},
					})
				}
			}
		}
	}

	// Check mapping count limit
	if len(tmpl.Mappings) > 200 {
		matches = append(matches, rules.Match{
			Message: fmt.Sprintf("Template has %d mappings, exceeding the 200 mapping limit", len(tmpl.Mappings)),
			Line:    1,
			Column:  1,
			Path:    []string{"Mappings"},
		})
	}

	return matches
}
