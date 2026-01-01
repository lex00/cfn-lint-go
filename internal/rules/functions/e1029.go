// Package functions contains intrinsic function validation rules (E1xxx).
package functions

import (
	"fmt"
	"regexp"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1029{})
}

// E1029 checks that variable substitution syntax is only used within Fn::Sub.
type E1029 struct{}

func (r *E1029) ID() string { return "E1029" }

func (r *E1029) ShortDesc() string {
	return "Sub required for variable substitution"
}

func (r *E1029) Description() string {
	return "Checks that ${variable} syntax is only used within Fn::Sub functions."
}

func (r *E1029) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-sub.html"
}

func (r *E1029) Tags() []string {
	return []string{"functions", "sub"}
}

// Matches ${VarName} but not $${literal} or ${!literal}
var subRequiredVarPattern = regexp.MustCompile(`\$\{([^!$][^}]*)\}`)

func (r *E1029) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check all resources
	for resName, res := range tmpl.Resources {
		findSubRequiredErrors(res.Properties, []string{"Resources", resName, "Properties"}, false, &matches)
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		if out.Value != nil {
			findSubRequiredErrors(out.Value, []string{"Outputs", outName, "Value"}, false, &matches)
		}
		if out.Description != "" {
			if subRequiredVarPattern.MatchString(out.Description) {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Variable substitution syntax ${...} found outside Fn::Sub in output '%s' Description", outName),
					Path:    []string{"Outputs", outName, "Description"},
				})
			}
		}
	}

	return matches
}

func findSubRequiredErrors(v any, path []string, inSub bool, matches *[]rules.Match) {
	switch val := v.(type) {
	case string:
		if !inSub && subRequiredVarPattern.MatchString(val) {
			*matches = append(*matches, rules.Match{
				Message: "Variable substitution syntax ${...} found outside Fn::Sub",
				Path:    path,
			})
		}
	case map[string]any:
		// Check if we're inside a Fn::Sub
		if _, ok := val["Fn::Sub"]; ok {
			// Process Fn::Sub contents as being inside Sub
			for key, child := range val {
				findSubRequiredErrors(child, append(path, key), true, matches)
			}
		} else {
			for key, child := range val {
				findSubRequiredErrors(child, append(path, key), inSub, matches)
			}
		}
	case []any:
		for i, child := range val {
			findSubRequiredErrors(child, append(path, fmt.Sprintf("[%d]", i)), inSub, matches)
		}
	}
}
