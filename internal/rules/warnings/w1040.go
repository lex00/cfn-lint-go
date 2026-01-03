package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W1040{})
}

// W1040 warns about ToJsonString function value issues.
type W1040 struct{}

func (r *W1040) ID() string { return "W1040" }

func (r *W1040) ShortDesc() string {
	return "ToJsonString function value validation"
}

func (r *W1040) Description() string {
	return "Warns about potential issues with Fn::ToJsonString values, such as serializing simple values."
}

func (r *W1040) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-ToJsonString.html"
}

func (r *W1040) Tags() []string {
	return []string{"warnings", "functions", "tojsonstring"}
}

func (r *W1040) Match(tmpl *template.Template) []rules.Match {
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

func (r *W1040) checkValue(v any, path []string, matches *[]rules.Match) {
	switch val := v.(type) {
	case map[string]any:
		if toJson, ok := val["Fn::ToJsonString"]; ok {
			r.checkToJsonString(toJson, path, matches)
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

func (r *W1040) checkToJsonString(toJson any, path []string, matches *[]rules.Match) {
	// Check if serializing a simple value that doesn't need ToJsonString
	switch v := toJson.(type) {
	case string:
		*matches = append(*matches, rules.Match{
			Message: "Fn::ToJsonString is used on a string; consider using the string directly",
			Path:    path,
		})
	case int, float64:
		*matches = append(*matches, rules.Match{
			Message: "Fn::ToJsonString is used on a number; consider using the number directly",
			Path:    path,
		})
	case bool:
		*matches = append(*matches, rules.Match{
			Message: "Fn::ToJsonString is used on a boolean; consider using the boolean directly",
			Path:    path,
		})
	case map[string]any:
		// Check if it's an empty object
		if len(v) == 0 {
			*matches = append(*matches, rules.Match{
				Message: "Fn::ToJsonString is used on an empty object; this will result in '{}'",
				Path:    path,
			})
		}
	case []any:
		// Check if it's an empty array
		if len(v) == 0 {
			*matches = append(*matches, rules.Match{
				Message: "Fn::ToJsonString is used on an empty array; this will result in '[]'",
				Path:    path,
			})
		}
	}
}
