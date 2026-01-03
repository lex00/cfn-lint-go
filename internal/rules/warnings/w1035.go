package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W1035{})
}

// W1035 warns about Select function value issues.
type W1035 struct{}

func (r *W1035) ID() string { return "W1035" }

func (r *W1035) ShortDesc() string {
	return "Select function value validation"
}

func (r *W1035) Description() string {
	return "Warns about potential issues with Fn::Select values, such as static index that could be simplified."
}

func (r *W1035) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-select.html"
}

func (r *W1035) Tags() []string {
	return []string{"warnings", "functions", "select"}
}

func (r *W1035) Match(tmpl *template.Template) []rules.Match {
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

func (r *W1035) checkValue(v any, path []string, matches *[]rules.Match) {
	switch val := v.(type) {
	case map[string]any:
		if sel, ok := val["Fn::Select"]; ok {
			r.checkSelect(sel, path, matches)
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

func (r *W1035) checkSelect(sel any, path []string, matches *[]rules.Match) {
	arr, ok := sel.([]any)
	if !ok || len(arr) != 2 {
		return
	}

	// Get index (can be string or int)
	var index int
	switch idx := arr[0].(type) {
	case string:
		// Parse string index
		_, _ = fmt.Sscanf(idx, "%d", &index)
	case int:
		index = idx
	case float64:
		index = int(idx)
	default:
		return // Dynamic index, can't validate
	}

	// Check if selecting from a static list where the element could be used directly
	if list, ok := arr[1].([]any); ok {
		if index >= 0 && index < len(list) {
			// Check if the selected element is a simple literal
			switch list[index].(type) {
			case string, int, float64, bool:
				*matches = append(*matches, rules.Match{
					Message: fmt.Sprintf("Fn::Select with static index %d from a static list; consider using the value directly", index),
					Path:    path,
				})
			}
		}
	}

	// Warn about negative index (invalid but worth mentioning)
	if index < 0 {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("Fn::Select uses negative index %d which is invalid", index),
			Path:    path,
		})
	}
}
