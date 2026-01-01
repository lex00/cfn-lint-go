// Package functions contains intrinsic function validation rules (E1xxx).
package functions

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1020{})
}

// E1020 checks that Ref values are strings.
type E1020 struct{}

func (r *E1020) ID() string { return "E1020" }

func (r *E1020) ShortDesc() string {
	return "Ref value must be a string"
}

func (r *E1020) Description() string {
	return "Checks that Ref intrinsic function values are strings, not arrays or objects."
}

func (r *E1020) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-ref.html"
}

func (r *E1020) Tags() []string {
	return []string{"functions", "ref"}
}

func (r *E1020) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check all resources
	for resName, res := range tmpl.Resources {
		refs := findAllRefValues(res.Properties)
		for _, ref := range refs {
			if !ref.isString {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Ref value must be a string, got %s in resource '%s'", ref.valueType, resName),
					Line:    ref.line,
					Column:  ref.column,
					Path:    []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		refs := findAllRefValues(out.Value)
		for _, ref := range refs {
			if !ref.isString {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Ref value must be a string, got %s in output '%s'", ref.valueType, outName),
					Line:    ref.line,
					Column:  ref.column,
					Path:    []string{"Outputs", outName, "Value"},
				})
			}
		}
	}

	// Check conditions
	for condName, cond := range tmpl.Conditions {
		refs := findAllRefValues(cond.Expression)
		for _, ref := range refs {
			if !ref.isString {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Ref value must be a string, got %s in condition '%s'", ref.valueType, condName),
					Line:    ref.line,
					Column:  ref.column,
					Path:    []string{"Conditions", condName},
				})
			}
		}
	}

	return matches
}

type refValueInfo struct {
	isString  bool
	valueType string
	line      int
	column    int
}

func findAllRefValues(v any) []refValueInfo {
	var results []refValueInfo
	findRefValuesRecursive(v, &results)
	return results
}

func findRefValuesRecursive(v any, results *[]refValueInfo) {
	switch val := v.(type) {
	case map[string]any:
		if refVal, ok := val["Ref"]; ok {
			info := refValueInfo{}
			switch refVal.(type) {
			case string:
				info.isString = true
				info.valueType = "string"
			case []any:
				info.isString = false
				info.valueType = "array"
			case map[string]any:
				info.isString = false
				info.valueType = "object"
			default:
				info.isString = false
				info.valueType = fmt.Sprintf("%T", refVal)
			}
			*results = append(*results, info)
		}
		for _, child := range val {
			findRefValuesRecursive(child, results)
		}
	case []any:
		for _, child := range val {
			findRefValuesRecursive(child, results)
		}
	}
}
