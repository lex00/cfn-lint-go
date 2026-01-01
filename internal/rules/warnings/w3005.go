// Package warnings contains warning-level rules (Wxxx).
package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W3005{})
}

// W3005 warns about redundant DependsOn declarations.
type W3005 struct{}

func (r *W3005) ID() string { return "W3005" }

func (r *W3005) ShortDesc() string {
	return "Redundant DependsOn"
}

func (r *W3005) Description() string {
	return "Warns when DependsOn specifies a resource that is already implicitly depended on through Ref or GetAtt."
}

func (r *W3005) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-attribute-dependson.html"
}

func (r *W3005) Tags() []string {
	return []string{"warnings", "resources", "dependson"}
}

func (r *W3005) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if len(res.DependsOn) == 0 {
			continue
		}

		// Find implicit dependencies through Ref and GetAtt
		implicitDeps := make(map[string]bool)
		findImplicitDependencies(res.Properties, implicitDeps, tmpl.Resources)

		// Check each DependsOn entry
		for _, dep := range res.DependsOn {
			if implicitDeps[dep] {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("DependsOn '%s' in resource '%s' is redundant (already an implicit dependency via Ref/GetAtt)", dep, resName),
					Path:    []string{"Resources", resName, "DependsOn"},
				})
			}
		}
	}

	return matches
}

func findImplicitDependencies(v any, deps map[string]bool, resources map[string]*template.Resource) {
	switch val := v.(type) {
	case map[string]any:
		// Check for Ref to a resource
		if ref, ok := val["Ref"].(string); ok {
			if _, isResource := resources[ref]; isResource {
				deps[ref] = true
			}
		}
		// Check for GetAtt
		if getAtt, ok := val["Fn::GetAtt"]; ok {
			switch ga := getAtt.(type) {
			case []any:
				if len(ga) >= 1 {
					if resName, ok := ga[0].(string); ok {
						if _, isResource := resources[resName]; isResource {
							deps[resName] = true
						}
					}
				}
			case string:
				// "Resource.Attribute" format
				for i, c := range ga {
					if c == '.' {
						resName := ga[:i]
						if _, isResource := resources[resName]; isResource {
							deps[resName] = true
						}
						break
					}
				}
			}
		}
		// Check Fn::Sub for resource references
		if sub, ok := val["Fn::Sub"]; ok {
			findSubResourceRefs(sub, deps, resources)
		}
		// Recurse
		for _, child := range val {
			findImplicitDependencies(child, deps, resources)
		}
	case []any:
		for _, child := range val {
			findImplicitDependencies(child, deps, resources)
		}
	}
}

func findSubResourceRefs(v any, deps map[string]bool, resources map[string]*template.Resource) {
	// Fn::Sub can reference resources via ${ResourceName} or ${ResourceName.Attribute}
	// This is a simplified check - a full implementation would parse the Sub string
	switch sub := v.(type) {
	case string:
		for resName := range resources {
			// Check for ${ResourceName} or ${ResourceName.
			if containsSubRef(sub, resName) {
				deps[resName] = true
			}
		}
	case []any:
		if len(sub) >= 1 {
			if str, ok := sub[0].(string); ok {
				for resName := range resources {
					if containsSubRef(str, resName) {
						deps[resName] = true
					}
				}
			}
		}
	}
}

func containsSubRef(str, resName string) bool {
	// Look for ${ResourceName} or ${ResourceName.Attr}
	pattern1 := "${" + resName + "}"
	pattern2 := "${" + resName + "."
	return len(str) > 0 && (indexOf(str, pattern1) >= 0 || indexOf(str, pattern2) >= 0)
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
