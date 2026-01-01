// Package functions contains intrinsic function validation rules (E1xxx).
package functions

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1016{})
}

// E1016 checks that Fn::ImportValue is properly configured.
type E1016 struct{}

func (r *E1016) ID() string { return "E1016" }

func (r *E1016) ShortDesc() string {
	return "Fn::ImportValue function error"
}

func (r *E1016) Description() string {
	return "Checks that Fn::ImportValue is properly configured with a string or intrinsic function."
}

func (r *E1016) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-importvalue.html"
}

func (r *E1016) Tags() []string {
	return []string{"functions", "importvalue"}
}

func (r *E1016) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check all resources
	for resName, res := range tmpl.Resources {
		findImportValueErrors(res.Properties, []string{"Resources", resName, "Properties"}, &matches)
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		if out.Value != nil {
			findImportValueErrors(out.Value, []string{"Outputs", outName, "Value"}, &matches)
		}
	}

	return matches
}

func findImportValueErrors(v any, path []string, matches *[]rules.Match) {
	switch val := v.(type) {
	case map[string]any:
		if importVal, ok := val["Fn::ImportValue"]; ok {
			if !isValidImportValueArg(importVal) {
				*matches = append(*matches, rules.Match{
					Message: fmt.Sprintf("Fn::ImportValue argument must be a string or intrinsic function, got %T", importVal),
					Path:    path,
				})
			}
		}
		for key, child := range val {
			findImportValueErrors(child, append(path, key), matches)
		}
	case []any:
		for i, child := range val {
			findImportValueErrors(child, append(path, fmt.Sprintf("[%d]", i)), matches)
		}
	}
}

func isValidImportValueArg(v any) bool {
	switch val := v.(type) {
	case string:
		return true
	case map[string]any:
		// Allow intrinsic functions
		allowedFunctions := []string{"Ref", "Fn::Sub", "Fn::If", "Fn::Join"}
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
