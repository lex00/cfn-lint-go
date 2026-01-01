// Package conditions contains condition validation rules (E8xxx).
package conditions

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E8007{})
}

// E8007 checks that the Condition intrinsic function references a valid condition.
type E8007 struct{}

func (r *E8007) ID() string { return "E8007" }

func (r *E8007) ShortDesc() string {
	return "Condition intrinsic function error"
}

func (r *E8007) Description() string {
	return "Checks that Condition intrinsic functions are used correctly within condition expressions."
}

func (r *E8007) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E8007"
}

func (r *E8007) Tags() []string {
	return []string{"conditions", "functions", "condition"}
}

func (r *E8007) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check that Condition intrinsic is only used within condition functions
	// and that it references a string (condition name)
	for name, cond := range tmpl.Conditions {
		found := findConditionIntrinsic(cond.Expression)
		for _, c := range found {
			if err := validateConditionIntrinsic(c); err != "" {
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

func validateConditionIntrinsic(value any) string {
	// The Condition intrinsic should have a string value (the condition name)
	_, ok := value.(string)
	if !ok {
		return "Condition intrinsic must reference a condition name (string)"
	}
	return ""
}

func findConditionIntrinsic(v any) []any {
	var results []any

	switch val := v.(type) {
	case map[string]any:
		if cond, ok := val["Condition"]; ok {
			results = append(results, cond)
		}
		for _, child := range val {
			results = append(results, findConditionIntrinsic(child)...)
		}
	case []any:
		for _, child := range val {
			results = append(results, findConditionIntrinsic(child)...)
		}
	}

	return results
}
