// Package functions contains intrinsic function validation rules (E1xxx).
package functions

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1028{})
}

// E1028 checks that Fn::If has correct structure and references valid conditions.
type E1028 struct{}

func (r *E1028) ID() string { return "E1028" }

func (r *E1028) ShortDesc() string {
	return "Fn::If structure error"
}

func (r *E1028) Description() string {
	return "Checks that Fn::If has exactly 3 elements and references a defined condition."
}

func (r *E1028) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-conditions.html#intrinsic-function-reference-conditions-if"
}

func (r *E1028) Tags() []string {
	return []string{"functions", "conditions", "if"}
}

func (r *E1028) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check all resources
	for resName, res := range tmpl.Resources {
		ifs := findAllFnIfs(res.Properties)
		for _, fnIf := range ifs {
			errs := r.validateFnIf(tmpl, fnIf)
			for _, err := range errs {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("%s in resource '%s'", err, resName),
					Line:    fnIf.line,
					Column:  fnIf.column,
					Path:    []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		ifs := findAllFnIfs(out.Value)
		for _, fnIf := range ifs {
			errs := r.validateFnIf(tmpl, fnIf)
			for _, err := range errs {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("%s in output '%s'", err, outName),
					Line:    fnIf.line,
					Column:  fnIf.column,
					Path:    []string{"Outputs", outName, "Value"},
				})
			}
		}
	}

	return matches
}

func (r *E1028) validateFnIf(tmpl *template.Template, fnIf fnIfInfo) []string {
	var errors []string

	// Check structure - must have exactly 3 elements
	if fnIf.argCount != 3 {
		errors = append(errors, fmt.Sprintf("Fn::If must have exactly 3 elements [condition_name, value_if_true, value_if_false], got %d", fnIf.argCount))
		return errors
	}

	// Check condition name is a string
	if fnIf.conditionName == "" {
		errors = append(errors, "Fn::If first element must be a condition name string")
		return errors
	}

	// Check condition exists
	if _, ok := tmpl.Conditions[fnIf.conditionName]; !ok {
		errors = append(errors, fmt.Sprintf("Fn::If references undefined condition '%s'", fnIf.conditionName))
	}

	return errors
}

type fnIfInfo struct {
	conditionName string
	argCount      int
	line          int
	column        int
}

func findAllFnIfs(v any) []fnIfInfo {
	var results []fnIfInfo
	findFnIfsRecursive(v, &results)
	return results
}

func findFnIfsRecursive(v any, results *[]fnIfInfo) {
	switch val := v.(type) {
	case map[string]any:
		if fnIf, ok := val["Fn::If"]; ok {
			info := parseFnIfValue(fnIf)
			*results = append(*results, info)
		}
		for _, child := range val {
			findFnIfsRecursive(child, results)
		}
	case []any:
		for _, child := range val {
			findFnIfsRecursive(child, results)
		}
	}
}

func parseFnIfValue(v any) fnIfInfo {
	arr, ok := v.([]any)
	if !ok {
		return fnIfInfo{}
	}

	info := fnIfInfo{argCount: len(arr)}

	if len(arr) >= 1 {
		if name, ok := arr[0].(string); ok {
			info.conditionName = name
		}
	}

	return info
}
