package informational

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&I1022{})
}

// I1022 suggests using Fn::Sub instead of Fn::Join for simple string concatenation.
type I1022 struct{}

func (r *I1022) ID() string { return "I1022" }

func (r *I1022) ShortDesc() string {
	return "Prefer Fn::Sub over Fn::Join"
}

func (r *I1022) Description() string {
	return "Suggests using Fn::Sub instead of Fn::Join when concatenating strings, as it is more readable and maintainable."
}

func (r *I1022) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-sub.html"
}

func (r *I1022) Tags() []string {
	return []string{"functions", "join", "sub", "best-practice"}
}

func (r *I1022) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check all resources
	for resName, res := range tmpl.Resources {
		findSimpleJoins(res.Properties, []string{"Resources", resName, "Properties"}, &matches)
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		if out.Value != nil {
			findSimpleJoins(out.Value, []string{"Outputs", outName, "Value"}, &matches)
		}
	}

	// Check parameters (Default values)
	for paramName, param := range tmpl.Parameters {
		if param.Default != nil {
			findSimpleJoins(param.Default, []string{"Parameters", paramName, "Default"}, &matches)
		}
	}

	return matches
}

func findSimpleJoins(v any, path []string, matches *[]rules.Match) {
	switch val := v.(type) {
	case map[string]any:
		if joinVal, ok := val["Fn::Join"]; ok {
			checkJoinForSubReplacement(joinVal, path, matches)
		}
		for key, child := range val {
			findSimpleJoins(child, append(path, key), matches)
		}
	case []any:
		for i, child := range val {
			findSimpleJoins(child, append(path, fmt.Sprintf("[%d]", i)), matches)
		}
	}
}

func checkJoinForSubReplacement(v any, path []string, matches *[]rules.Match) {
	arr, ok := v.([]any)
	if !ok || len(arr) != 2 {
		return
	}

	delimiter, ok := arr[0].(string)
	if !ok {
		return
	}

	values, ok := arr[1].([]any)
	if !ok {
		return
	}

	// Check if this is a simple concatenation that could use Fn::Sub
	// Criteria:
	// 1. Empty delimiter "" (simple string concatenation)
	// 2. Multiple values with mix of strings and intrinsics
	// 3. Not too complex (e.g., not nested joins)
	if delimiter == "" && len(values) > 1 {
		hasString := false
		hasIntrinsic := false
		tooComplex := false

		for _, val := range values {
			switch v := val.(type) {
			case string:
				hasString = true
			case map[string]any:
				hasIntrinsic = true
				// Check for nested joins or complex structures
				if _, isJoin := v["Fn::Join"]; isJoin {
					tooComplex = true
				}
			}
		}

		// Suggest Fn::Sub if we have both strings and intrinsics
		if hasString && hasIntrinsic && !tooComplex {
			*matches = append(*matches, rules.Match{
				Message: "Consider using Fn::Sub instead of Fn::Join with empty delimiter for improved readability",
				Path:    append(path, "Fn::Join"),
			})
		}
	}

	// Also check for common patterns like joining with "-" or "/" with simple values
	if (delimiter == "-" || delimiter == "/" || delimiter == ":") && len(values) >= 2 {
		simplePattern := true
		for _, val := range values {
			switch v := val.(type) {
			case string:
				// OK
			case map[string]any:
				// Skip nested joins
				if _, isJoin := v["Fn::Join"]; isJoin {
					simplePattern = false
					break
				}
				// Allow simple Ref or AWS::pseudo parameters
				if len(v) > 1 {
					simplePattern = false
				}
			default:
				simplePattern = false
			}
		}

		if simplePattern {
			// Build example Sub format
			example := buildSubExample(delimiter, values)
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf("Consider using Fn::Sub instead of Fn::Join. Example: !Sub '%s'", example),
				Path:    append(path, "Fn::Join"),
			})
		}
	}
}

func buildSubExample(delimiter string, values []any) string {
	var parts []string
	for _, val := range values {
		switch v := val.(type) {
		case string:
			parts = append(parts, v)
		case map[string]any:
			if ref, ok := v["Ref"]; ok {
				if refStr, ok := ref.(string); ok {
					parts = append(parts, "${"+refStr+"}")
				}
			} else if getAtt, ok := v["Fn::GetAtt"]; ok {
				if getAttArr, ok := getAtt.([]any); ok && len(getAttArr) == 2 {
					if res, ok := getAttArr[0].(string); ok {
						if attr, ok := getAttArr[1].(string); ok {
							parts = append(parts, "${"+res+"."+attr+"}")
						}
					}
				}
			}
		}
	}
	return strings.Join(parts, delimiter)
}
