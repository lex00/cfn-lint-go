package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W1033{})
}

// W1033 warns about Split function value issues.
type W1033 struct{}

func (r *W1033) ID() string { return "W1033" }

func (r *W1033) ShortDesc() string {
	return "Split function value validation"
}

func (r *W1033) Description() string {
	return "Warns about potential issues with Fn::Split values, such as empty delimiters."
}

func (r *W1033) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-split.html"
}

func (r *W1033) Tags() []string {
	return []string{"warnings", "functions", "split"}
}

func (r *W1033) Match(tmpl *template.Template) []rules.Match {
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

func (r *W1033) checkValue(v any, path []string, matches *[]rules.Match) {
	switch val := v.(type) {
	case map[string]any:
		if split, ok := val["Fn::Split"]; ok {
			r.checkSplit(split, path, matches)
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

func (r *W1033) checkSplit(split any, path []string, matches *[]rules.Match) {
	arr, ok := split.([]any)
	if !ok || len(arr) != 2 {
		return
	}

	// Check for empty delimiter
	if delimiter, ok := arr[0].(string); ok {
		if delimiter == "" {
			*matches = append(*matches, rules.Match{
				Message: "Fn::Split uses an empty delimiter which will split each character",
				Path:    path,
			})
		}
	}

	// Check if splitting a literal string that doesn't contain the delimiter
	if delimiter, ok := arr[0].(string); ok && delimiter != "" {
		if str, ok := arr[1].(string); ok {
			found := false
			for i := 0; i <= len(str)-len(delimiter); i++ {
				if str[i:i+len(delimiter)] == delimiter {
					found = true
					break
				}
			}
			if !found && str != "" {
				*matches = append(*matches, rules.Match{
					Message: fmt.Sprintf("Fn::Split delimiter '%s' is not found in the literal string; result will be a single-element list", delimiter),
					Path:    path,
				})
			}
		}
	}
}
