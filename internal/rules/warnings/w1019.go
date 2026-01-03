package warnings

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W1019{})
}

// W1019 warns when Fn::Sub has parameters that are never used in the string.
type W1019 struct{}

func (r *W1019) ID() string { return "W1019" }

func (r *W1019) ShortDesc() string {
	return "Unused Sub parameters"
}

func (r *W1019) Description() string {
	return "Warns when Fn::Sub defines parameters that are not used in the substitution string."
}

func (r *W1019) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-sub.html"
}

func (r *W1019) Tags() []string {
	return []string{"warnings", "functions", "sub"}
}

// subVarPattern matches ${VarName} in Sub strings
var subVarPattern = regexp.MustCompile(`\$\{([^}!][^}]*)\}`)

func (r *W1019) Match(tmpl *template.Template) []rules.Match {
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

func (r *W1019) checkValue(v any, path []string, matches *[]rules.Match) {
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

func (r *W1019) checkSub(sub any, path []string, matches *[]rules.Match) {
	// Fn::Sub can be a string or [string, {params}]
	subArray, ok := sub.([]any)
	if !ok || len(subArray) != 2 {
		return // Only check when there are explicit parameters
	}

	subStr, ok := subArray[0].(string)
	if !ok {
		return
	}

	params, ok := subArray[1].(map[string]any)
	if !ok {
		return
	}

	// Find all variables used in the string
	usedVars := make(map[string]bool)
	varMatches := subVarPattern.FindAllStringSubmatch(subStr, -1)
	for _, m := range varMatches {
		if len(m) >= 2 {
			// Handle nested references like ${AWS::StackName} or ${MyVar}
			varName := m[1]
			// Strip any .Attribute suffix for GetAtt-style references
			if idx := strings.Index(varName, "."); idx > 0 {
				varName = varName[:idx]
			}
			usedVars[varName] = true
		}
	}

	// Check for unused parameters
	for paramName := range params {
		if !usedVars[paramName] {
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf("Fn::Sub parameter '%s' is defined but never used in the substitution string", paramName),
				Path:    path,
			})
		}
	}
}
