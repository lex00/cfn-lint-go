// Package functions contains intrinsic function validation rules (E1xxx).
package functions

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1031{})
}

// E1031 checks that Fn::ToJsonString is configured correctly.
type E1031 struct{}

func (r *E1031) ID() string { return "E1031" }

func (r *E1031) ShortDesc() string {
	return "Fn::ToJsonString function validation"
}

func (r *E1031) Description() string {
	return "Validates that Fn::ToJsonString is only used when AWS::LanguageExtensions transform is present."
}

func (r *E1031) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-tojsonstring.html"
}

func (r *E1031) Tags() []string {
	return []string{"functions", "tojsonstring"}
}

func (r *E1031) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check if AWS::LanguageExtensions transform is present
	hasLanguageExtensions := hasLanguageExtensionsTransform(tmpl)

	// Find all Fn::ToJsonString uses
	for resName, res := range tmpl.Resources {
		toJsonStrings := findAllFnToJsonString(res.Properties)
		if len(toJsonStrings) > 0 && !hasLanguageExtensions {
			for _, tjs := range toJsonStrings {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Fn::ToJsonString is not supported without the AWS::LanguageExtensions transform in resource '%s'", resName),
					Line:    tjs.line,
					Column:  tjs.column,
					Path:    []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		toJsonStrings := findAllFnToJsonString(out.Value)
		if len(toJsonStrings) > 0 && !hasLanguageExtensions {
			for _, tjs := range toJsonStrings {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Fn::ToJsonString is not supported without the AWS::LanguageExtensions transform in output '%s'", outName),
					Line:    tjs.line,
					Column:  tjs.column,
					Path:    []string{"Outputs", outName, "Value"},
				})
			}
		}
	}

	return matches
}

type fnToJsonStringInfo struct {
	value  any
	line   int
	column int
}

func findAllFnToJsonString(v any) []fnToJsonStringInfo {
	var results []fnToJsonStringInfo
	findFnToJsonStringRecursive(v, &results)
	return results
}

func findFnToJsonStringRecursive(v any, results *[]fnToJsonStringInfo) {
	switch val := v.(type) {
	case map[string]any:
		if toJsonStringVal, ok := val["Fn::ToJsonString"]; ok {
			*results = append(*results, fnToJsonStringInfo{
				value: toJsonStringVal,
			})
		}
		for _, child := range val {
			findFnToJsonStringRecursive(child, results)
		}
	case []any:
		for _, child := range val {
			findFnToJsonStringRecursive(child, results)
		}
	}
}
