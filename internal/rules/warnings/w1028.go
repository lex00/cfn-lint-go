package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W1028{})
}

// W1028 warns when Fn::If has an unreachable path due to static condition evaluation.
type W1028 struct{}

func (r *W1028) ID() string { return "W1028" }

func (r *W1028) ShortDesc() string {
	return "Fn::If unreachable path"
}

func (r *W1028) Description() string {
	return "Warns when Fn::If has a path that can never be reached because the condition always evaluates to true or false."
}

func (r *W1028) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-conditions.html"
}

func (r *W1028) Tags() []string {
	return []string{"warnings", "functions", "conditions", "if"}
}

func (r *W1028) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Analyze conditions that are always true or false
	staticConditions := r.findStaticConditions(tmpl)

	// Check resources for Fn::If using static conditions
	for resName, res := range tmpl.Resources {
		r.checkValue(res.Properties, []string{"Resources", resName, "Properties"}, staticConditions, &matches)
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		r.checkValue(out.Value, []string{"Outputs", outName, "Value"}, staticConditions, &matches)
	}

	return matches
}

func (r *W1028) findStaticConditions(tmpl *template.Template) map[string]bool {
	static := make(map[string]bool)

	for condName, cond := range tmpl.Conditions {
		if result, isStatic := r.evaluateCondition(cond.Expression); isStatic {
			static[condName] = result
		}
	}

	return static
}

func (r *W1028) evaluateCondition(v any) (bool, bool) {
	condMap, ok := v.(map[string]any)
	if !ok {
		return false, false
	}

	// Check Fn::Equals with two identical literal values
	if equals, ok := condMap["Fn::Equals"]; ok {
		if arr, ok := equals.([]any); ok && len(arr) == 2 {
			// Both values are identical strings
			if str1, ok1 := arr[0].(string); ok1 {
				if str2, ok2 := arr[1].(string); ok2 {
					return str1 == str2, true
				}
			}
			// Both values are identical numbers
			if num1, ok1 := arr[0].(float64); ok1 {
				if num2, ok2 := arr[1].(float64); ok2 {
					return num1 == num2, true
				}
			}
			if num1, ok1 := arr[0].(int); ok1 {
				if num2, ok2 := arr[1].(int); ok2 {
					return num1 == num2, true
				}
			}
		}
	}

	return false, false
}

func (r *W1028) checkValue(v any, path []string, staticConditions map[string]bool, matches *[]rules.Match) {
	switch val := v.(type) {
	case map[string]any:
		if fnIf, ok := val["Fn::If"]; ok {
			if arr, ok := fnIf.([]any); ok && len(arr) >= 1 {
				if condName, ok := arr[0].(string); ok {
					if result, isStatic := staticConditions[condName]; isStatic {
						if result {
							*matches = append(*matches, rules.Match{
								Message: fmt.Sprintf("Fn::If condition '%s' always evaluates to true; the false branch is unreachable", condName),
								Path:    path,
							})
						} else {
							*matches = append(*matches, rules.Match{
								Message: fmt.Sprintf("Fn::If condition '%s' always evaluates to false; the true branch is unreachable", condName),
								Path:    path,
							})
						}
					}
				}
			}
		}
		for key, child := range val {
			r.checkValue(child, append(path, key), staticConditions, matches)
		}
	case []any:
		for i, child := range val {
			r.checkValue(child, append(path, fmt.Sprintf("[%d]", i)), staticConditions, matches)
		}
	}
}
