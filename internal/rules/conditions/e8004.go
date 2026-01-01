// Package conditions contains condition validation rules (E8xxx).
package conditions

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E8004{})
}

// E8004 checks that Fn::And has valid structure.
type E8004 struct{}

func (r *E8004) ID() string { return "E8004" }

func (r *E8004) ShortDesc() string {
	return "Fn::And structure error"
}

func (r *E8004) Description() string {
	return "Checks that Fn::And has 2-10 condition elements."
}

func (r *E8004) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E8004"
}

func (r *E8004) Tags() []string {
	return []string{"conditions", "functions", "and"}
}

func (r *E8004) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check conditions
	for name, cond := range tmpl.Conditions {
		found := findFnAnd(cond.Expression)
		for _, and := range found {
			if err := validateFnAnd(and); err != "" {
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

func validateFnAnd(value any) string {
	arr, ok := value.([]any)
	if !ok {
		return "Fn::And must be a list"
	}
	if len(arr) < 2 {
		return fmt.Sprintf("Fn::And must have at least 2 conditions, got %d", len(arr))
	}
	if len(arr) > 10 {
		return fmt.Sprintf("Fn::And must have at most 10 conditions, got %d", len(arr))
	}
	return ""
}

func findFnAnd(v any) []any {
	var results []any

	switch val := v.(type) {
	case map[string]any:
		if and, ok := val["Fn::And"]; ok {
			results = append(results, and)
		}
		for _, child := range val {
			results = append(results, findFnAnd(child)...)
		}
	case []any:
		for _, child := range val {
			results = append(results, findFnAnd(child)...)
		}
	}

	return results
}
