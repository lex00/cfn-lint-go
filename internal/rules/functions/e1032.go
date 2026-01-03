// Package functions contains intrinsic function validation rules (E1xxx).
package functions

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1032{})
}

// E1032 checks that Fn::ForEach is configured correctly.
type E1032 struct{}

func (r *E1032) ID() string { return "E1032" }

func (r *E1032) ShortDesc() string {
	return "Fn::ForEach function validation"
}

func (r *E1032) Description() string {
	return "Validates that Fn::ForEach parameters have a valid configuration and AWS::LanguageExtensions transform is present."
}

func (r *E1032) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-foreach.html"
}

func (r *E1032) Tags() []string {
	return []string{"functions", "foreach"}
}

func (r *E1032) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check if AWS::LanguageExtensions transform is present
	hasLanguageExtensions := hasLanguageExtensionsTransform(tmpl)

	// Fn::ForEach appears as resource names like "Fn::ForEach::Buckets"
	// Check for resource names starting with "Fn::ForEach::"
	for resName := range tmpl.Resources {
		if len(resName) >= 13 && resName[:13] == "Fn::ForEach::" {
			if !hasLanguageExtensions {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Missing Transform: Declare the AWS::LanguageExtensions Transform globally to enable use of the intrinsic function Fn::ForEach in resource '%s'", resName),
					Path:    []string{"Resources", resName},
				})
			}
		}
	}

	// Also check for Fn::ForEach in properties (alternative usage)
	for resName, res := range tmpl.Resources {
		forEachs := findAllFnForEach(res.Properties)
		if len(forEachs) > 0 && !hasLanguageExtensions {
			for _, fe := range forEachs {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Missing Transform: Declare the AWS::LanguageExtensions Transform globally to enable use of the intrinsic function Fn::ForEach in resource '%s'", resName),
					Line:    fe.line,
					Column:  fe.column,
					Path:    []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		forEachs := findAllFnForEach(out.Value)
		if len(forEachs) > 0 && !hasLanguageExtensions {
			for _, fe := range forEachs {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Missing Transform: Declare the AWS::LanguageExtensions Transform globally to enable use of the intrinsic function Fn::ForEach in output '%s'", outName),
					Line:    fe.line,
					Column:  fe.column,
					Path:    []string{"Outputs", outName, "Value"},
				})
			}
		}
	}

	return matches
}

type fnForEachInfo struct {
	value  any
	line   int
	column int
}

func findAllFnForEach(v any) []fnForEachInfo {
	var results []fnForEachInfo
	findFnForEachRecursive(v, &results)
	return results
}

func findFnForEachRecursive(v any, results *[]fnForEachInfo) {
	switch val := v.(type) {
	case map[string]any:
		if forEachVal, ok := val["Fn::ForEach"]; ok {
			*results = append(*results, fnForEachInfo{
				value: forEachVal,
			})
		}
		for _, child := range val {
			findFnForEachRecursive(child, results)
		}
	case []any:
		for _, child := range val {
			findFnForEachRecursive(child, results)
		}
	}
}
