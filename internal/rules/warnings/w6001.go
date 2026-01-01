// Package warnings contains warning-level rules (Wxxx).
package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W6001{})
}

// W6001 warns when Fn::ImportValue is used in an Output value.
type W6001 struct{}

func (r *W6001) ID() string { return "W6001" }

func (r *W6001) ShortDesc() string {
	return "ImportValue in Output"
}

func (r *W6001) Description() string {
	return "Warns when Fn::ImportValue is used in an Output value, which creates a circular dependency risk between stacks."
}

func (r *W6001) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-importvalue.html"
}

func (r *W6001) Tags() []string {
	return []string{"warnings", "outputs", "importvalue"}
}

func (r *W6001) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for outName, out := range tmpl.Outputs {
		if hasImportValue(out.Value) {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Output '%s' uses Fn::ImportValue, which can create circular dependencies between stacks", outName),
				Path:    []string{"Outputs", outName, "Value"},
			})
		}
	}

	return matches
}

func hasImportValue(v any) bool {
	switch val := v.(type) {
	case map[string]any:
		if _, ok := val["Fn::ImportValue"]; ok {
			return true
		}
		for _, child := range val {
			if hasImportValue(child) {
				return true
			}
		}
	case []any:
		for _, child := range val {
			if hasImportValue(child) {
				return true
			}
		}
	}
	return false
}
