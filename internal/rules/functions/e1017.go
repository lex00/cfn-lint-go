// Package functions contains intrinsic function validation rules (E1xxx).
package functions

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1017{})
}

// E1017 checks that Fn::Select is properly configured.
type E1017 struct{}

func (r *E1017) ID() string { return "E1017" }

func (r *E1017) ShortDesc() string {
	return "Fn::Select function error"
}

func (r *E1017) Description() string {
	return "Checks that Fn::Select is properly configured with [index, list]."
}

func (r *E1017) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-select.html"
}

func (r *E1017) Tags() []string {
	return []string{"functions", "select"}
}

func (r *E1017) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check all resources
	for resName, res := range tmpl.Resources {
		findSelectErrors(res.Properties, []string{"Resources", resName, "Properties"}, &matches)
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		if out.Value != nil {
			findSelectErrors(out.Value, []string{"Outputs", outName, "Value"}, &matches)
		}
	}

	return matches
}

func findSelectErrors(v any, path []string, matches *[]rules.Match) {
	switch val := v.(type) {
	case map[string]any:
		if selectVal, ok := val["Fn::Select"]; ok {
			validateSelectArgs(selectVal, path, matches)
		}
		for key, child := range val {
			findSelectErrors(child, append(path, key), matches)
		}
	case []any:
		for i, child := range val {
			findSelectErrors(child, append(path, fmt.Sprintf("[%d]", i)), matches)
		}
	}
}

func validateSelectArgs(v any, path []string, matches *[]rules.Match) {
	arr, ok := v.([]any)
	if !ok {
		*matches = append(*matches, rules.Match{
			Message: "Fn::Select requires an array of [index, list]",
			Path:    path,
		})
		return
	}

	if len(arr) != 2 {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("Fn::Select requires exactly 2 elements [index, list], got %d", len(arr)),
			Path:    path,
		})
		return
	}

	// First element must be an index (integer or intrinsic)
	if !isValidSelectIndex(arr[0]) {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("Fn::Select index must be an integer or intrinsic function, got %T", arr[0]),
			Path:    path,
		})
	}

	// Second element must be an array or intrinsic that returns an array
	if !isValidSelectList(arr[1]) {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("Fn::Select list must be an array or intrinsic function, got %T", arr[1]),
			Path:    path,
		})
	}
}

func isValidSelectIndex(v any) bool {
	switch val := v.(type) {
	case int, int64, float64:
		return true
	case string:
		// String representation of number
		return true
	case map[string]any:
		// Allow intrinsic functions
		return len(val) > 0
	default:
		return false
	}
}

func isValidSelectList(v any) bool {
	switch val := v.(type) {
	case []any:
		return true
	case map[string]any:
		// Allow intrinsic functions that return arrays
		allowedFunctions := []string{
			"Ref", "Fn::GetAZs", "Fn::Split", "Fn::If",
			"Fn::FindInMap", "Fn::GetAtt",
		}
		for _, fn := range allowedFunctions {
			if _, ok := val[fn]; ok {
				return true
			}
		}
		return false
	default:
		return false
	}
}
