// Package warnings contains warning-level rules (Wxxx).
package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W1001{})
}

// W1001 warns when Ref/GetAtt references a resource that has a Condition.
// The referenced resource might not exist if its condition evaluates to false.
type W1001 struct{}

func (r *W1001) ID() string { return "W1001" }

func (r *W1001) ShortDesc() string {
	return "Ref/GetAtt to conditional resource"
}

func (r *W1001) Description() string {
	return "Warns when Ref or GetAtt references a resource with a Condition, as the resource may not exist."
}

func (r *W1001) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-conditions.html"
}

func (r *W1001) Tags() []string {
	return []string{"warnings", "functions", "ref", "getatt", "conditions"}
}

func (r *W1001) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Build set of conditional resources
	conditionalResources := make(map[string]string)
	for resName, res := range tmpl.Resources {
		if res.Condition != "" {
			conditionalResources[resName] = res.Condition
		}
	}

	// If no conditional resources, nothing to check
	if len(conditionalResources) == 0 {
		return matches
	}

	// Check all resources for Refs/GetAtts to conditional resources
	for resName, res := range tmpl.Resources {
		// Skip checking within the same conditional resource
		refs := findAllRefTargets(res.Properties)
		for _, ref := range refs {
			if condition, ok := conditionalResources[ref]; ok && res.Condition != condition {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Ref to '%s' in resource '%s' references a conditional resource (Condition: %s)", ref, resName, condition),
					Path:    []string{"Resources", resName, "Properties"},
				})
			}
		}

		getAtts := findAllGetAttTargets(res.Properties)
		for _, target := range getAtts {
			if condition, ok := conditionalResources[target]; ok && res.Condition != condition {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("GetAtt to '%s' in resource '%s' references a conditional resource (Condition: %s)", target, resName, condition),
					Path:    []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		refs := findAllRefTargets(out.Value)
		for _, ref := range refs {
			if condition, ok := conditionalResources[ref]; ok && out.Condition != condition {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Ref to '%s' in output '%s' references a conditional resource (Condition: %s)", ref, outName, condition),
					Path:    []string{"Outputs", outName, "Value"},
				})
			}
		}

		getAtts := findAllGetAttTargets(out.Value)
		for _, target := range getAtts {
			if condition, ok := conditionalResources[target]; ok && out.Condition != condition {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("GetAtt to '%s' in output '%s' references a conditional resource (Condition: %s)", target, outName, condition),
					Path:    []string{"Outputs", outName, "Value"},
				})
			}
		}
	}

	return matches
}

func findAllRefTargets(v any) []string {
	var refs []string
	findRefTargetsRecursive(v, &refs)
	return refs
}

func findRefTargetsRecursive(v any, refs *[]string) {
	switch val := v.(type) {
	case map[string]any:
		if ref, ok := val["Ref"].(string); ok {
			*refs = append(*refs, ref)
		}
		for _, child := range val {
			findRefTargetsRecursive(child, refs)
		}
	case []any:
		for _, child := range val {
			findRefTargetsRecursive(child, refs)
		}
	}
}

func findAllGetAttTargets(v any) []string {
	var targets []string
	findGetAttTargetsRecursive(v, &targets)
	return targets
}

func findGetAttTargetsRecursive(v any, targets *[]string) {
	switch val := v.(type) {
	case map[string]any:
		if getAtt, ok := val["Fn::GetAtt"]; ok {
			switch ga := getAtt.(type) {
			case []any:
				if len(ga) >= 1 {
					if resName, ok := ga[0].(string); ok {
						*targets = append(*targets, resName)
					}
				}
			case string:
				// "Resource.Attribute" format
				for i, c := range ga {
					if c == '.' {
						*targets = append(*targets, ga[:i])
						break
					}
				}
			}
		}
		for _, child := range val {
			findGetAttTargetsRecursive(child, targets)
		}
	case []any:
		for _, child := range val {
			findGetAttTargetsRecursive(child, targets)
		}
	}
}
