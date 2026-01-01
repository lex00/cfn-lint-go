// Package functions contains intrinsic function validation rules (E1xxx).
package functions

import (
	"regexp"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1027{})
}

// E1027 checks that dynamic references are used in valid locations.
type E1027 struct{}

func (r *E1027) ID() string { return "E1027" }

func (r *E1027) ShortDesc() string {
	return "Dynamic reference in invalid location"
}

func (r *E1027) Description() string {
	return "Checks that dynamic references ({{resolve:...}}) are only used in supported locations."
}

func (r *E1027) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/dynamic-references.html"
}

func (r *E1027) Tags() []string {
	return []string{"functions", "dynamic-references"}
}

var dynamicRefPattern = regexp.MustCompile(`\{\{resolve:[^}]+\}\}`)

func (r *E1027) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check parameters - dynamic references not allowed in Default values
	for paramName, param := range tmpl.Parameters {
		if param.Default != nil {
			if str, ok := param.Default.(string); ok {
				if dynamicRefPattern.MatchString(str) {
					matches = append(matches, rules.Match{
						Message: "Dynamic references are not allowed in parameter Default values",
						Line:    param.Node.Line,
						Column:  param.Node.Column,
						Path:    []string{"Parameters", paramName, "Default"},
					})
				}
			}
		}
	}

	// Check outputs - dynamic references not allowed in Export Name
	for outName, out := range tmpl.Outputs {
		if out.Export != nil {
			if name, ok := out.Export["Name"].(string); ok {
				if dynamicRefPattern.MatchString(name) {
					matches = append(matches, rules.Match{
						Message: "Dynamic references are not allowed in output Export Name",
						Line:    out.Node.Line,
						Column:  out.Node.Column,
						Path:    []string{"Outputs", outName, "Export", "Name"},
					})
				}
			}
		}
	}

	// Check resource DependsOn - dynamic references not allowed
	for resName, res := range tmpl.Resources {
		for _, dep := range res.DependsOn {
			if dynamicRefPattern.MatchString(dep) {
				matches = append(matches, rules.Match{
					Message: "Dynamic references are not allowed in DependsOn",
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", resName, "DependsOn"},
				})
			}
		}

		// Check resource Condition
		if res.Condition != "" && dynamicRefPattern.MatchString(res.Condition) {
			matches = append(matches, rules.Match{
				Message: "Dynamic references are not allowed in resource Condition",
				Line:    res.Node.Line,
				Column:  res.Node.Column,
				Path:    []string{"Resources", resName, "Condition"},
			})
		}
	}

	// Check conditions section - dynamic references not allowed
	for condName, cond := range tmpl.Conditions {
		if hasDynamicRef(cond.Expression) {
			matches = append(matches, rules.Match{
				Message: "Dynamic references are not allowed in Conditions",
				Line:    cond.Node.Line,
				Column:  cond.Node.Column,
				Path:    []string{"Conditions", condName},
			})
		}
	}

	return matches
}

func hasDynamicRef(v any) bool {
	switch val := v.(type) {
	case string:
		return dynamicRefPattern.MatchString(val)
	case map[string]any:
		for _, child := range val {
			if hasDynamicRef(child) {
				return true
			}
		}
	case []any:
		for _, child := range val {
			if hasDynamicRef(child) {
				return true
			}
		}
	}
	return false
}

// containsDynamicRef checks if a string contains a dynamic reference pattern
func containsDynamicRef(s string) bool {
	return strings.Contains(s, "{{resolve:")
}
