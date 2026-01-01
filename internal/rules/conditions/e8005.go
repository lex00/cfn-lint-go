// Package conditions contains condition validation rules (E8xxx).
package conditions

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E8005{})
}

// E8005 checks that Fn::Not has valid structure.
type E8005 struct{}

func (r *E8005) ID() string { return "E8005" }

func (r *E8005) ShortDesc() string {
	return "Fn::Not structure error"
}

func (r *E8005) Description() string {
	return "Checks that Fn::Not has exactly one condition element."
}

func (r *E8005) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E8005"
}

func (r *E8005) Tags() []string {
	return []string{"conditions", "functions", "not"}
}

func (r *E8005) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check conditions
	for name, cond := range tmpl.Conditions {
		found := findFnNot(cond.Expression)
		for _, not := range found {
			if err := validateFnNot(not); err != "" {
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

func validateFnNot(value any) string {
	arr, ok := value.([]any)
	if !ok {
		return "Fn::Not must be a list"
	}
	if len(arr) != 1 {
		return fmt.Sprintf("Fn::Not must have exactly 1 condition, got %d", len(arr))
	}
	return ""
}

func findFnNot(v any) []any {
	var results []any

	switch val := v.(type) {
	case map[string]any:
		if not, ok := val["Fn::Not"]; ok {
			results = append(results, not)
		}
		for _, child := range val {
			results = append(results, findFnNot(child)...)
		}
	case []any:
		for _, child := range val {
			results = append(results, findFnNot(child)...)
		}
	}

	return results
}
