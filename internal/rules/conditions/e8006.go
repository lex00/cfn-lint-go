// Package conditions contains condition validation rules (E8xxx).
package conditions

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E8006{})
}

// E8006 checks that Fn::Or has valid structure.
type E8006 struct{}

func (r *E8006) ID() string { return "E8006" }

func (r *E8006) ShortDesc() string {
	return "Fn::Or structure error"
}

func (r *E8006) Description() string {
	return "Checks that Fn::Or has 2-10 condition elements."
}

func (r *E8006) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E8006"
}

func (r *E8006) Tags() []string {
	return []string{"conditions", "functions", "or"}
}

func (r *E8006) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check conditions
	for name, cond := range tmpl.Conditions {
		found := findFnOr(cond.Expression)
		for _, or := range found {
			if err := validateFnOr(or); err != "" {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Condition '%s': %s", name, err),
					Line:    cond.Node.Line,
					Column:  cond.Node.Column,
					Path:    []string{"Conditions", name},
				})
			}
		}
	}

	return matches
}

func validateFnOr(value any) string {
	arr, ok := value.([]any)
	if !ok {
		return "Fn::Or must be a list"
	}
	if len(arr) < 2 {
		return fmt.Sprintf("Fn::Or must have at least 2 conditions, got %d", len(arr))
	}
	if len(arr) > 10 {
		return fmt.Sprintf("Fn::Or must have at most 10 conditions, got %d", len(arr))
	}
	return ""
}

func findFnOr(v any) []any {
	var results []any

	switch val := v.(type) {
	case map[string]any:
		if or, ok := val["Fn::Or"]; ok {
			results = append(results, or)
		}
		for _, child := range val {
			results = append(results, findFnOr(child)...)
		}
	case []any:
		for _, child := range val {
			results = append(results, findFnOr(child)...)
		}
	}

	return results
}
