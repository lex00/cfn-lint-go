// Package metadata contains metadata validation rules (E4xxx).
package metadata

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E4002{})
}

// E4002 validates that the Metadata section is properly configured.
type E4002 struct{}

func (r *E4002) ID() string { return "E4002" }

func (r *E4002) ShortDesc() string {
	return "Metadata section is valid"
}

func (r *E4002) Description() string {
	return "Validates that the Metadata section is an object with no null values."
}

func (r *E4002) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E4002"
}

func (r *E4002) Tags() []string {
	return []string{"metadata", "configuration"}
}

func (r *E4002) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check for null values in Metadata
	if tmpl.MetadataNode != nil {
		nullPaths := findNullValues(tmpl.Metadata, []string{"Metadata"})
		for _, path := range nullPaths {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Metadata contains null value at path: %v", path),
				Line:    tmpl.MetadataNode.Line,
				Column:  tmpl.MetadataNode.Column,
				Path:    path,
			})
		}
	}

	return matches
}

// findNullValues recursively finds null values in a map structure.
func findNullValues(v any, path []string) [][]string {
	var nullPaths [][]string

	switch val := v.(type) {
	case nil:
		nullPaths = append(nullPaths, path)
	case map[string]any:
		for key, child := range val {
			childPath := append(append([]string{}, path...), key)
			nullPaths = append(nullPaths, findNullValues(child, childPath)...)
		}
	case []any:
		for i, child := range val {
			childPath := append(append([]string{}, path...), fmt.Sprintf("[%d]", i))
			nullPaths = append(nullPaths, findNullValues(child, childPath)...)
		}
	}

	return nullPaths
}
