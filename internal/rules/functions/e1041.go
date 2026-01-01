// Package functions contains intrinsic function validation rules (E1xxx).
package functions

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1041{})
}

// E1041 checks that Ref has valid format.
type E1041 struct{}

func (r *E1041) ID() string { return "E1041" }

func (r *E1041) ShortDesc() string {
	return "Ref format error"
}

func (r *E1041) Description() string {
	return "Checks that Ref has valid format: must be a non-empty string."
}

func (r *E1041) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-ref.html"
}

func (r *E1041) Tags() []string {
	return []string{"functions", "ref"}
}

func (r *E1041) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check all resources
	for resName, res := range tmpl.Resources {
		refs := findAllRefFormats(res.Properties)
		for _, ref := range refs {
			if err := r.validateRefFormat(ref); err != "" {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("%s in resource '%s'", err, resName),
					Line:    ref.line,
					Column:  ref.column,
					Path:    []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		refs := findAllRefFormats(out.Value)
		for _, ref := range refs {
			if err := r.validateRefFormat(ref); err != "" {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("%s in output '%s'", err, outName),
					Line:    ref.line,
					Column:  ref.column,
					Path:    []string{"Outputs", outName, "Value"},
				})
			}
		}
	}

	// Check conditions
	for condName, cond := range tmpl.Conditions {
		refs := findAllRefFormats(cond.Expression)
		for _, ref := range refs {
			if err := r.validateRefFormat(ref); err != "" {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("%s in condition '%s'", err, condName),
					Line:    ref.line,
					Column:  ref.column,
					Path:    []string{"Conditions", condName},
				})
			}
		}
	}

	return matches
}

func (r *E1041) validateRefFormat(ref refFormatInfo) string {
	if !ref.isString {
		return fmt.Sprintf("Ref must be a string, got %s", ref.valueType)
	}
	if ref.value == "" {
		return "Ref value cannot be empty"
	}
	return ""
}

type refFormatInfo struct {
	isString  bool
	valueType string
	value     string
	line      int
	column    int
}

func findAllRefFormats(v any) []refFormatInfo {
	var results []refFormatInfo
	findRefFormatsRecursive(v, &results)
	return results
}

func findRefFormatsRecursive(v any, results *[]refFormatInfo) {
	switch val := v.(type) {
	case map[string]any:
		if refVal, ok := val["Ref"]; ok {
			info := refFormatInfo{}
			switch rv := refVal.(type) {
			case string:
				info.isString = true
				info.valueType = "string"
				info.value = rv
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
			findRefFormatsRecursive(child, results)
		}
	case []any:
		for _, child := range val {
			findRefFormatsRecursive(child, results)
		}
	}
}
