// Package functions contains intrinsic function validation rules (E1xxx).
package functions

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1030{})
}

// E1030 checks that Fn::Length is configured correctly.
type E1030 struct{}

func (r *E1030) ID() string { return "E1030" }

func (r *E1030) ShortDesc() string {
	return "Fn::Length function validation"
}

func (r *E1030) Description() string {
	return "Validates that Fn::Length is only used when AWS::LanguageExtensions transform is present."
}

func (r *E1030) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-length.html"
}

func (r *E1030) Tags() []string {
	return []string{"functions", "length"}
}

func (r *E1030) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check if AWS::LanguageExtensions transform is present
	hasLanguageExtensions := hasLanguageExtensionsTransform(tmpl)

	// Find all Fn::Length uses
	for resName, res := range tmpl.Resources {
		lengths := findAllFnLength(res.Properties)
		if len(lengths) > 0 && !hasLanguageExtensions {
			for _, length := range lengths {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Fn::Length is not supported without the AWS::LanguageExtensions transform in resource '%s'", resName),
					Line:    length.line,
					Column:  length.column,
					Path:    []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		lengths := findAllFnLength(out.Value)
		if len(lengths) > 0 && !hasLanguageExtensions {
			for _, length := range lengths {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Fn::Length is not supported without the AWS::LanguageExtensions transform in output '%s'", outName),
					Line:    length.line,
					Column:  length.column,
					Path:    []string{"Outputs", outName, "Value"},
				})
			}
		}
	}

	return matches
}

// hasLanguageExtensionsTransform checks if AWS::LanguageExtensions transform is present
func hasLanguageExtensionsTransform(tmpl *template.Template) bool {
	if tmpl.Transform == nil {
		return false
	}

	// Transform can be a string or array of strings
	switch t := tmpl.Transform.(type) {
	case string:
		return t == "AWS::LanguageExtensions"
	case []any:
		for _, transform := range t {
			if str, ok := transform.(string); ok && str == "AWS::LanguageExtensions" {
				return true
			}
		}
	}

	return false
}

type fnLengthInfo struct {
	value  any
	line   int
	column int
}

func findAllFnLength(v any) []fnLengthInfo {
	var results []fnLengthInfo
	findFnLengthRecursive(v, &results)
	return results
}

func findFnLengthRecursive(v any, results *[]fnLengthInfo) {
	switch val := v.(type) {
	case map[string]any:
		if lengthVal, ok := val["Fn::Length"]; ok {
			*results = append(*results, fnLengthInfo{
				value: lengthVal,
			})
		}
		for _, child := range val {
			findFnLengthRecursive(child, results)
		}
	case []any:
		for _, child := range val {
			findFnLengthRecursive(child, results)
		}
	}
}
