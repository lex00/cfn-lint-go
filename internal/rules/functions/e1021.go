// Package functions contains intrinsic function validation rules (E1xxx).
package functions

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1021{})
}

// E1021 checks that Fn::Base64 is properly configured.
type E1021 struct{}

func (r *E1021) ID() string { return "E1021" }

func (r *E1021) ShortDesc() string {
	return "Fn::Base64 function error"
}

func (r *E1021) Description() string {
	return "Checks that Fn::Base64 is properly configured with a string value."
}

func (r *E1021) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-base64.html"
}

func (r *E1021) Tags() []string {
	return []string{"functions", "base64"}
}

func (r *E1021) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check all resources
	for resName, res := range tmpl.Resources {
		findBase64Errors(res.Properties, []string{"Resources", resName, "Properties"}, &matches)
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		if out.Value != nil {
			findBase64Errors(out.Value, []string{"Outputs", outName, "Value"}, &matches)
		}
	}

	return matches
}

func findBase64Errors(v any, path []string, matches *[]rules.Match) {
	switch val := v.(type) {
	case map[string]any:
		if base64Val, ok := val["Fn::Base64"]; ok {
			if !isValidBase64Arg(base64Val) {
				*matches = append(*matches, rules.Match{
					Message: fmt.Sprintf("Fn::Base64 argument must be a string or intrinsic function, got %T", base64Val),
					Path:    path,
				})
			}
		}
		for key, child := range val {
			findBase64Errors(child, append(path, key), matches)
		}
	case []any:
		for i, child := range val {
			findBase64Errors(child, append(path, fmt.Sprintf("[%d]", i)), matches)
		}
	}
}

func isValidBase64Arg(v any) bool {
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
