// Package warnings contains warning-level rules (Wxxx).
package warnings

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W2001{})
}

// W2001 warns about parameters that are defined but never used.
type W2001 struct{}

func (r *W2001) ID() string { return "W2001" }

func (r *W2001) ShortDesc() string {
	return "Unused parameter"
}

func (r *W2001) Description() string {
	return "Warns when a parameter is defined but never referenced in the template."
}

func (r *W2001) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/parameters-section-structure.html"
}

func (r *W2001) Tags() []string {
	return []string{"warnings", "parameters", "unused"}
}

func (r *W2001) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Find all parameter references in the template
	usedParams := make(map[string]bool)

	// Check Resources
	for _, res := range tmpl.Resources {
		findParamRefs(res.Properties, usedParams, tmpl.Parameters)
	}

	// Check Outputs
	for _, out := range tmpl.Outputs {
		findParamRefs(out.Value, usedParams, tmpl.Parameters)
	}

	// Check Conditions
	for _, cond := range tmpl.Conditions {
		findParamRefs(cond.Expression, usedParams, tmpl.Parameters)
	}

	// Check Metadata
	findParamRefs(tmpl.Metadata, usedParams, tmpl.Parameters)

	// Report unused parameters
	for paramName := range tmpl.Parameters {
		if !usedParams[paramName] {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Parameter '%s' is defined but never used", paramName),
				Path:    []string{"Parameters", paramName},
			})
		}
	}

	return matches
}

func findParamRefs(v any, usedParams map[string]bool, params map[string]*template.Parameter) {
	switch val := v.(type) {
	case map[string]any:
		// Check for Ref
		if ref, ok := val["Ref"].(string); ok {
			if _, isParam := params[ref]; isParam {
				usedParams[ref] = true
			}
		}
		// Check for Fn::Sub with parameter references
		if sub, ok := val["Fn::Sub"]; ok {
			findSubParamRefs(sub, usedParams, params)
		}
		// Recurse into children
		for _, child := range val {
			findParamRefs(child, usedParams, params)
		}
	case []any:
		for _, child := range val {
			findParamRefs(child, usedParams, params)
		}
	}
}

func findSubParamRefs(v any, usedParams map[string]bool, params map[string]*template.Parameter) {
	switch sub := v.(type) {
	case string:
		// Find ${ParamName} references
		for paramName := range params {
			if strings.Contains(sub, "${"+paramName+"}") {
				usedParams[paramName] = true
			}
		}
	case []any:
		// [string, {VarName: value}] format
		if len(sub) >= 1 {
			if str, ok := sub[0].(string); ok {
				for paramName := range params {
					if strings.Contains(str, "${"+paramName+"}") {
						usedParams[paramName] = true
					}
				}
			}
		}
	}
}
