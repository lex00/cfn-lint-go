package warnings

import (
	"fmt"
	"regexp"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W1020{})
}

// W1020 warns when Fn::Sub is used but contains no variable substitutions.
type W1020 struct{}

func (r *W1020) ID() string { return "W1020" }

func (r *W1020) ShortDesc() string {
	return "Sub not needed without variables"
}

func (r *W1020) Description() string {
	return "Warns when Fn::Sub is used but the string contains no variable substitutions, making Sub unnecessary."
}

func (r *W1020) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-sub.html"
}

func (r *W1020) Tags() []string {
	return []string{"warnings", "functions", "sub"}
}

// subVarPatternW1020 matches ${VarName} in Sub strings (excludes ${! literal escapes)
var subVarPatternW1020 = regexp.MustCompile(`\$\{[^}!][^}]*\}`)

func (r *W1020) Match(tmpl *template.Template) []rules.Match {
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

func (r *W1020) checkValue(v any, path []string, matches *[]rules.Match) {
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

func (r *W1020) checkSub(sub any, path []string, matches *[]rules.Match) {
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

	// Check if the string contains any variable substitutions
	if !subVarPatternW1020.MatchString(subStr) {
		*matches = append(*matches, rules.Match{
			Message: "Fn::Sub is used but contains no variable substitutions; consider using a plain string instead",
			Path:    path,
		})
	}
}
