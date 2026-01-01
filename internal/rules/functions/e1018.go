// Package functions contains intrinsic function validation rules (E1xxx).
package functions

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1018{})
}

// E1018 checks that Fn::Split is properly configured.
type E1018 struct{}

func (r *E1018) ID() string { return "E1018" }

func (r *E1018) ShortDesc() string {
	return "Fn::Split function error"
}

func (r *E1018) Description() string {
	return "Checks that Fn::Split is properly configured with [delimiter, source_string]."
}

func (r *E1018) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-split.html"
}

func (r *E1018) Tags() []string {
	return []string{"functions", "split"}
}

func (r *E1018) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check all resources
	for resName, res := range tmpl.Resources {
		findSplitErrors(res.Properties, []string{"Resources", resName, "Properties"}, &matches)
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		if out.Value != nil {
			findSplitErrors(out.Value, []string{"Outputs", outName, "Value"}, &matches)
		}
	}

	return matches
}

func findSplitErrors(v any, path []string, matches *[]rules.Match) {
	switch val := v.(type) {
	case map[string]any:
		if splitVal, ok := val["Fn::Split"]; ok {
			validateSplitArgs(splitVal, path, matches)
		}
		for key, child := range val {
			findSplitErrors(child, append(path, key), matches)
		}
	case []any:
		for i, child := range val {
			findSplitErrors(child, append(path, fmt.Sprintf("[%d]", i)), matches)
		}
	}
}

func validateSplitArgs(v any, path []string, matches *[]rules.Match) {
	arr, ok := v.([]any)
	if !ok {
		*matches = append(*matches, rules.Match{
			Message: "Fn::Split requires an array of [delimiter, source_string]",
			Path:    path,
		})
		return
	}

	if len(arr) != 2 {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("Fn::Split requires exactly 2 elements [delimiter, source_string], got %d", len(arr)),
			Path:    path,
		})
		return
	}

	// First element must be a delimiter string
	if _, ok := arr[0].(string); !ok {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("Fn::Split delimiter must be a string, got %T", arr[0]),
			Path:    path,
		})
	}

	// Second element must be a string or intrinsic that returns a string
	if !isValidSplitSource(arr[1]) {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("Fn::Split source must be a string or intrinsic function, got %T", arr[1]),
			Path:    path,
		})
	}
}

func isValidSplitSource(v any) bool {
	switch val := v.(type) {
	case string:
		return true
	case map[string]any:
		// Allow intrinsic functions that return strings
		return len(val) > 0
	default:
		return false
	}
}
