// Package conditions contains condition validation rules (E8xxx).
package conditions

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E8002{})
}

// E8002 checks that referenced conditions are defined.
type E8002 struct{}

func (r *E8002) ID() string { return "E8002" }

func (r *E8002) ShortDesc() string {
	return "Referenced Conditions are defined"
}

func (r *E8002) Description() string {
	return "Checks that all Condition references point to defined conditions in the Conditions section."
}

func (r *E8002) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E8002"
}

func (r *E8002) Tags() []string {
	return []string{"conditions", "reference"}
}

func (r *E8002) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Build set of defined conditions
	definedConditions := make(map[string]bool)
	for name := range tmpl.Conditions {
		definedConditions[name] = true
	}

	// Check resource Condition attributes
	for resName, res := range tmpl.Resources {
		if res.Condition != "" && !definedConditions[res.Condition] {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Resource '%s' references undefined condition '%s'", resName, res.Condition),
				Line:    res.Node.Line,
				Column:  res.Node.Column,
				Path:    []string{"Resources", resName, "Condition"},
			})
		}

		// Check Fn::If in properties
		condRefs := findConditionRefs(res.Properties)
		for _, condRef := range condRefs {
			if !definedConditions[condRef] {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Resource '%s' references undefined condition '%s' in Fn::If", resName, condRef),
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	// Check output Condition attributes
	for outName, out := range tmpl.Outputs {
		if out.Condition != "" && !definedConditions[out.Condition] {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Output '%s' references undefined condition '%s'", outName, out.Condition),
				Line:    out.Node.Line,
				Column:  out.Node.Column,
				Path:    []string{"Outputs", outName, "Condition"},
			})
		}
	}

	// Check Condition references within conditions (Fn::And, Fn::Or, Fn::Not using Condition)
	for condName, cond := range tmpl.Conditions {
		condRefs := findConditionRefs(cond.Expression)
		for _, ref := range condRefs {
			if !definedConditions[ref] {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Condition '%s' references undefined condition '%s'", condName, ref),
					Line:    cond.Node.Line,
					Column:  cond.Node.Column,
					Path:    []string{"Conditions", condName},
				})
			}
		}
	}

	return matches
}

// findConditionRefs recursively finds all Condition and Fn::If references in a value.
func findConditionRefs(v any) []string {
	var refs []string

	switch val := v.(type) {
	case map[string]any:
		// Check for Condition intrinsic
		if condName, ok := val["Condition"].(string); ok {
			refs = append(refs, condName)
		}
		// Check for Fn::If (first element is condition name)
		if fnIf, ok := val["Fn::If"].([]any); ok && len(fnIf) > 0 {
			if condName, ok := fnIf[0].(string); ok {
				refs = append(refs, condName)
			}
		}
		// Recurse into all values
		for _, child := range val {
			refs = append(refs, findConditionRefs(child)...)
		}
	case []any:
		for _, child := range val {
			refs = append(refs, findConditionRefs(child)...)
		}
	}

	return refs
}
