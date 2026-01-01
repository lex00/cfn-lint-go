// Package functions contains intrinsic function validation rules (E1xxx).
package functions

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1015{})
}

// E1015 checks that Fn::GetAZs is properly configured.
type E1015 struct{}

func (r *E1015) ID() string { return "E1015" }

func (r *E1015) ShortDesc() string {
	return "Fn::GetAZs function error"
}

func (r *E1015) Description() string {
	return "Checks that Fn::GetAZs is properly configured with a string region or empty string."
}

func (r *E1015) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getavailabilityzones.html"
}

func (r *E1015) Tags() []string {
	return []string{"functions", "getazs"}
}

func (r *E1015) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check all resources
	for resName, res := range tmpl.Resources {
		findGetAZsErrors(res.Properties, []string{"Resources", resName, "Properties"}, &matches)
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		if out.Value != nil {
			findGetAZsErrors(out.Value, []string{"Outputs", outName, "Value"}, &matches)
		}
	}

	return matches
}

func findGetAZsErrors(v any, path []string, matches *[]rules.Match) {
	switch val := v.(type) {
	case map[string]any:
		if getAZs, ok := val["Fn::GetAZs"]; ok {
			if !isValidGetAZsArg(getAZs) {
				*matches = append(*matches, rules.Match{
					Message: fmt.Sprintf("Fn::GetAZs argument must be a string (region name) or empty string, got %T", getAZs),
					Path:    path,
				})
			}
		}
		for key, child := range val {
			findGetAZsErrors(child, append(path, key), matches)
		}
	case []any:
		for i, child := range val {
			findGetAZsErrors(child, append(path, fmt.Sprintf("[%d]", i)), matches)
		}
	}
}

func isValidGetAZsArg(v any) bool {
	switch val := v.(type) {
	case string:
		return true // Empty string or region name
	case map[string]any:
		// Allow intrinsic functions like Ref
		_, hasRef := val["Ref"]
		_, hasSub := val["Fn::Sub"]
		_, hasIf := val["Fn::If"]
		return hasRef || hasSub || hasIf
	default:
		return false
	}
}
