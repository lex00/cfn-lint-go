package warnings

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W1031{})
}

// W1031 warns about Sub function value issues that may cause problems.
type W1031 struct{}

func (r *W1031) ID() string { return "W1031" }

func (r *W1031) ShortDesc() string {
	return "Sub function value validation"
}

func (r *W1031) Description() string {
	return "Warns about potential issues with Fn::Sub values, such as unclosed variable brackets."
}

func (r *W1031) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-sub.html"
}

func (r *W1031) Tags() []string {
	return []string{"warnings", "functions", "sub"}
}

func (r *W1031) Match(tmpl *template.Template) []rules.Match {
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

func (r *W1031) checkValue(v any, path []string, matches *[]rules.Match) {
	switch val := v.(type) {
	case map[string]any:
		if sub, ok := val["Fn::Sub"]; ok {
			r.checkSub(sub, path, matches)
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

func (r *W1031) checkSub(sub any, path []string, matches *[]rules.Match) {
	var subStr string

	switch s := sub.(type) {
	case string:
		subStr = s
	case []any:
		if len(s) >= 1 {
			if str, ok := s[0].(string); ok {
				subStr = str
			}
		}
	default:
		return
	}

	if subStr == "" {
		return
	}

	// Check for unbalanced ${} brackets
	openCount := strings.Count(subStr, "${")
	closeCount := strings.Count(subStr, "}")

	if openCount > closeCount {
		*matches = append(*matches, rules.Match{
			Message: "Fn::Sub string has unclosed variable bracket '${'",
			Path:    path,
		})
	}

	// Check for empty variable references ${}
	if strings.Contains(subStr, "${}") {
		*matches = append(*matches, rules.Match{
			Message: "Fn::Sub string contains empty variable reference '${}'",
			Path:    path,
		})
	}

	// Check for $ without { which might be unintended
	for i := 0; i < len(subStr)-1; i++ {
		if subStr[i] == '$' && subStr[i+1] != '{' && subStr[i+1] != '$' {
			*matches = append(*matches, rules.Match{
				Message: "Fn::Sub string contains '$' not followed by '{' which may be unintended",
				Path:    path,
			})
			break
		}
	}
}
