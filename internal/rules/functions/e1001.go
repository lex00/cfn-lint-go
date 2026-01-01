// Package functions contains intrinsic function validation rules (E1xxx).
package functions

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1001{})
}

// E1001 checks that Ref references point to existing resources or parameters.
type E1001 struct{}

func (r *E1001) ID() string { return "E1001" }

func (r *E1001) ShortDesc() string {
	return "Ref to undefined resource or parameter"
}

func (r *E1001) Description() string {
	return "Checks that all Ref intrinsic functions reference valid resources or parameters."
}

func (r *E1001) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E1001"
}

func (r *E1001) Tags() []string {
	return []string{"functions", "ref"}
}

func (r *E1001) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Build set of valid references
	validRefs := make(map[string]bool)
	for name := range tmpl.Resources {
		validRefs[name] = true
	}
	for name := range tmpl.Parameters {
		validRefs[name] = true
	}
	// Add pseudo-parameters
	for _, pseudo := range []string{
		"AWS::AccountId",
		"AWS::NotificationARNs",
		"AWS::NoValue",
		"AWS::Partition",
		"AWS::Region",
		"AWS::StackId",
		"AWS::StackName",
		"AWS::URLSuffix",
	} {
		validRefs[pseudo] = true
	}

	// Check all resources for invalid Refs
	for resName, res := range tmpl.Resources {
		refs := findAllRefs(res.Properties)
		for _, ref := range refs {
			if !validRefs[ref.target] {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Ref '%s' in resource '%s' references undefined resource or parameter", ref.target, resName),
					Line:    ref.line,
					Column:  ref.column,
					Path:    []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	return matches
}

type refInfo struct {
	target string
	line   int
	column int
}

func findAllRefs(v any) []refInfo {
	var refs []refInfo
	findRefsRecursive(v, &refs)
	return refs
}

func findRefsRecursive(v any, refs *[]refInfo) {
	switch val := v.(type) {
	case map[string]any:
		if ref, ok := val["Ref"].(string); ok {
			*refs = append(*refs, refInfo{target: ref})
		}
		for _, child := range val {
			findRefsRecursive(child, refs)
		}
	case []any:
		for _, child := range val {
			findRefsRecursive(child, refs)
		}
	}
}
