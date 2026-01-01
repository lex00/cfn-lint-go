// Package warnings contains warning-level rules (Wxxx).
package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W2010{})
}

// W2010 warns when a NoEcho parameter is used in a way that might expose its value.
type W2010 struct{}

func (r *W2010) ID() string { return "W2010" }

func (r *W2010) ShortDesc() string {
	return "NoEcho parameter may be exposed"
}

func (r *W2010) Description() string {
	return "Warns when a NoEcho parameter is referenced in an Output value, which could expose the sensitive data."
}

func (r *W2010) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/parameters-section-structure.html#parameters-section-structure-properties-noecho"
}

func (r *W2010) Tags() []string {
	return []string{"warnings", "parameters", "security", "noecho"}
}

func (r *W2010) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Find all NoEcho parameters
	noEchoParams := make(map[string]bool)
	for paramName, param := range tmpl.Parameters {
		if param.NoEcho {
			noEchoParams[paramName] = true
		}
	}

	if len(noEchoParams) == 0 {
		return matches
	}

	// Check if NoEcho parameters are referenced in Outputs
	for outName, out := range tmpl.Outputs {
		refs := findParamRefsInValue(out.Value, noEchoParams)
		for _, paramName := range refs {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("NoEcho parameter '%s' is referenced in output '%s', which may expose its value", paramName, outName),
				Path:    []string{"Outputs", outName, "Value"},
			})
		}
	}

	return matches
}

func findParamRefsInValue(v any, targetParams map[string]bool) []string {
	var refs []string
	findParamRefsInValueRecursive(v, targetParams, &refs)
	return refs
}

func findParamRefsInValueRecursive(v any, targetParams map[string]bool, refs *[]string) {
	switch val := v.(type) {
	case map[string]any:
		if ref, ok := val["Ref"].(string); ok {
			if targetParams[ref] {
				*refs = append(*refs, ref)
			}
		}
		for _, child := range val {
			findParamRefsInValueRecursive(child, targetParams, refs)
		}
	case []any:
		for _, child := range val {
			findParamRefsInValueRecursive(child, targetParams, refs)
		}
	}
}
