package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W1032{})
}

// W1032 warns about Join function value issues.
type W1032 struct{}

func (r *W1032) ID() string { return "W1032" }

func (r *W1032) ShortDesc() string {
	return "Join function value validation"
}

func (r *W1032) Description() string {
	return "Warns about potential issues with Fn::Join values, such as joining a single element."
}

func (r *W1032) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-join.html"
}

func (r *W1032) Tags() []string {
	return []string{"warnings", "functions", "join"}
}

func (r *W1032) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check resources
	for resName, res := range tmpl.Resources {
		r.checkValue(res.Properties, []string{"Resources", resName, "Properties"}, &matches)
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		r.checkValue(out.Value, []string{"Outputs", outName, "Value"}, &matches)
	}

	return matches
}

func (r *W1032) checkValue(v any, path []string, matches *[]rules.Match) {
	switch val := v.(type) {
	case map[string]any:
		if join, ok := val["Fn::Join"]; ok {
			r.checkJoin(join, path, matches)
		}
		for key, child := range val {
			r.checkValue(child, append(path, key), matches)
		}
	case []any:
		for i, child := range val {
			r.checkValue(child, append(path, fmt.Sprintf("[%d]", i)), matches)
		}
	}
}

func (r *W1032) checkJoin(join any, path []string, matches *[]rules.Match) {
	arr, ok := join.([]any)
	if !ok || len(arr) != 2 {
		return
	}

	// Check if the list to join has only one element
	if list, ok := arr[1].([]any); ok {
		if len(list) == 1 {
			*matches = append(*matches, rules.Match{
				Message: "Fn::Join is used with only one element; consider using the element directly",
				Path:    path,
			})
		}
		if len(list) == 0 {
			*matches = append(*matches, rules.Match{
				Message: "Fn::Join is used with an empty list; this will result in an empty string",
				Path:    path,
			})
		}
	}

	// Check if delimiter is non-empty but list has only literal empty strings
	if delimiter, ok := arr[0].(string); ok && delimiter != "" {
		if list, ok := arr[1].([]any); ok {
			allEmpty := true
			for _, item := range list {
				if str, ok := item.(string); !ok || str != "" {
					allEmpty = false
					break
				}
			}
			if allEmpty && len(list) > 0 {
				*matches = append(*matches, rules.Match{
					Message: "Fn::Join has a delimiter but all list elements are empty strings",
					Path:    path,
				})
			}
		}
	}
}
