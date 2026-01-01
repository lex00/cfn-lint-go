// Package functions contains intrinsic function validation rules (E1xxx).
package functions

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1022{})
}

// E1022 checks that Fn::Join is properly configured.
type E1022 struct{}

func (r *E1022) ID() string { return "E1022" }

func (r *E1022) ShortDesc() string {
	return "Fn::Join function error"
}

func (r *E1022) Description() string {
	return "Checks that Fn::Join is properly configured with [delimiter, [values]]."
}

func (r *E1022) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-join.html"
}

func (r *E1022) Tags() []string {
	return []string{"functions", "join"}
}

func (r *E1022) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check all resources
	for resName, res := range tmpl.Resources {
		findJoinErrors(res.Properties, []string{"Resources", resName, "Properties"}, &matches)
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		if out.Value != nil {
			findJoinErrors(out.Value, []string{"Outputs", outName, "Value"}, &matches)
		}
	}

	return matches
}

func findJoinErrors(v any, path []string, matches *[]rules.Match) {
	switch val := v.(type) {
	case map[string]any:
		if joinVal, ok := val["Fn::Join"]; ok {
			validateJoinArgs(joinVal, path, matches)
		}
		for key, child := range val {
			findJoinErrors(child, append(path, key), matches)
		}
	case []any:
		for i, child := range val {
			findJoinErrors(child, append(path, fmt.Sprintf("[%d]", i)), matches)
		}
	}
}

func validateJoinArgs(v any, path []string, matches *[]rules.Match) {
	arr, ok := v.([]any)
	if !ok {
		*matches = append(*matches, rules.Match{
			Message: "Fn::Join requires an array of [delimiter, [values]]",
			Path:    path,
		})
		return
	}

	if len(arr) != 2 {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("Fn::Join requires exactly 2 elements [delimiter, [values]], got %d", len(arr)),
			Path:    path,
		})
		return
	}

	// First element must be a delimiter string
	if _, ok := arr[0].(string); !ok {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("Fn::Join delimiter must be a string, got %T", arr[0]),
			Path:    path,
		})
	}

	// Second element must be an array or intrinsic that returns an array
	if !isValidJoinValues(arr[1]) {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("Fn::Join values must be an array or intrinsic function, got %T", arr[1]),
			Path:    path,
		})
	}
}

func isValidJoinValues(v any) bool {
	switch val := v.(type) {
	case []any:
		return true
	case map[string]any:
		// Allow intrinsic functions that return arrays
		return len(val) > 0
	default:
		return false
	}
}
