// Package conditions contains condition validation rules (E8xxx).
package conditions

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E8003{})
}

// E8003 checks that Fn::Equals has valid structure.
type E8003 struct{}

func (r *E8003) ID() string { return "E8003" }

func (r *E8003) ShortDesc() string {
	return "Fn::Equals structure error"
}

func (r *E8003) Description() string {
	return "Checks that Fn::Equals has exactly two elements to compare."
}

func (r *E8003) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E8003"
}

func (r *E8003) Tags() []string {
	return []string{"conditions", "functions", "equals"}
}

func (r *E8003) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check conditions
	for name, cond := range tmpl.Conditions {
		found := findFnEquals(cond.Expression)
		for _, eq := range found {
			if err := validateFnEquals(eq); err != "" {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Condition '%s': %s", name, err),
					Line:    cond.Node.Line,
					Column:  cond.Node.Column,
					Path:    []string{"Conditions", name},
				})
			}
		}
	}

	// Check resources for Fn::Equals in Fn::If
	for resName, res := range tmpl.Resources {
		found := findFnEquals(res.Properties)
		for _, eq := range found {
			if err := validateFnEquals(eq); err != "" {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Resource '%s': %s", resName, err),
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	return matches
}

func validateFnEquals(value any) string {
	arr, ok := value.([]any)
	if !ok {
		return "Fn::Equals must be a list"
	}
	if len(arr) != 2 {
		return fmt.Sprintf("Fn::Equals must have exactly 2 elements, got %d", len(arr))
	}
	return ""
}

func findFnEquals(v any) []any {
	var results []any

	switch val := v.(type) {
	case map[string]any:
		if eq, ok := val["Fn::Equals"]; ok {
			results = append(results, eq)
		}
		for _, child := range val {
			results = append(results, findFnEquals(child)...)
		}
	case []any:
		for _, child := range val {
			results = append(results, findFnEquals(child)...)
		}
	}

	return results
}
