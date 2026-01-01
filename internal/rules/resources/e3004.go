// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3004{})
}

// E3004 checks for circular resource dependencies.
type E3004 struct{}

func (r *E3004) ID() string { return "E3004" }

func (r *E3004) ShortDesc() string {
	return "Circular resource dependency detected"
}

func (r *E3004) Description() string {
	return "Checks for circular dependencies between resources via DependsOn, Ref, and GetAtt."
}

func (r *E3004) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-attribute-dependson.html"
}

func (r *E3004) Tags() []string {
	return []string{"resources", "dependencies"}
}

func (r *E3004) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Build dependency graph
	deps := buildDependencyGraph(tmpl)

	// Find cycles using DFS
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for resName := range tmpl.Resources {
		if !visited[resName] {
			if cycle := findCycle(resName, deps, visited, recStack, []string{}); cycle != nil {
				res := tmpl.Resources[resName]
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Circular dependency detected: %s", strings.Join(cycle, " -> ")),
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", resName},
				})
				break // Report only one cycle to avoid duplicates
			}
		}
	}

	return matches
}

func buildDependencyGraph(tmpl *template.Template) map[string][]string {
	deps := make(map[string][]string)

	for resName, res := range tmpl.Resources {
		var resDeps []string

		// Add explicit DependsOn
		resDeps = append(resDeps, res.DependsOn...)

		// Add implicit dependencies from Ref
		refs := findResourceRefs(res.Properties, tmpl)
		resDeps = append(resDeps, refs...)

		// Add implicit dependencies from GetAtt
		getAtts := findResourceGetAtts(res.Properties, tmpl)
		resDeps = append(resDeps, getAtts...)

		// Deduplicate
		deps[resName] = uniqueStrings(resDeps)
	}

	return deps
}

func findResourceRefs(v any, tmpl *template.Template) []string {
	var refs []string
	findResourceRefsRecursive(v, tmpl, &refs)
	return refs
}

func findResourceRefsRecursive(v any, tmpl *template.Template, refs *[]string) {
	switch val := v.(type) {
	case map[string]any:
		if ref, ok := val["Ref"].(string); ok {
			// Only add if it's a resource (not a parameter or pseudo-parameter)
			if tmpl.HasResource(ref) {
				*refs = append(*refs, ref)
			}
		}
		for _, child := range val {
			findResourceRefsRecursive(child, tmpl, refs)
		}
	case []any:
		for _, child := range val {
			findResourceRefsRecursive(child, tmpl, refs)
		}
	}
}

func findResourceGetAtts(v any, tmpl *template.Template) []string {
	var refs []string
	findResourceGetAttsRecursive(v, tmpl, &refs)
	return refs
}

func findResourceGetAttsRecursive(v any, tmpl *template.Template, refs *[]string) {
	switch val := v.(type) {
	case map[string]any:
		if ga, ok := val["Fn::GetAtt"]; ok {
			resName := extractGetAttResource(ga)
			if resName != "" && tmpl.HasResource(resName) {
				*refs = append(*refs, resName)
			}
		}
		for _, child := range val {
			findResourceGetAttsRecursive(child, tmpl, refs)
		}
	case []any:
		for _, child := range val {
			findResourceGetAttsRecursive(child, tmpl, refs)
		}
	}
}

func extractGetAttResource(v any) string {
	switch val := v.(type) {
	case string:
		parts := strings.SplitN(val, ".", 2)
		if len(parts) >= 1 {
			return parts[0]
		}
	case []any:
		if len(val) >= 1 {
			if res, ok := val[0].(string); ok {
				return res
			}
		}
	}
	return ""
}

func findCycle(node string, deps map[string][]string, visited, recStack map[string]bool, path []string) []string {
	visited[node] = true
	recStack[node] = true
	path = append(path, node)

	for _, dep := range deps[node] {
		if !visited[dep] {
			if cycle := findCycle(dep, deps, visited, recStack, path); cycle != nil {
				return cycle
			}
		} else if recStack[dep] {
			// Found cycle - return path from dep back to dep
			cycle := append(path, dep)
			// Trim to start from the cycle start
			for i, n := range cycle {
				if n == dep {
					return cycle[i:]
				}
			}
			return cycle
		}
	}

	recStack[node] = false
	return nil
}

func uniqueStrings(s []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}
