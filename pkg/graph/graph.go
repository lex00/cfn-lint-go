// Package graph generates DOT format dependency graphs from CloudFormation templates.
package graph

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

// Generator creates DOT graphs from CloudFormation templates.
type Generator struct {
	// IncludeParameters includes parameter references in the graph.
	IncludeParameters bool
}

// Edge represents a dependency between resources.
type Edge struct {
	From string
	To   string
	Type string // "Ref", "GetAtt", "DependsOn"
}

// Generate creates a DOT graph from a template.
func (g *Generator) Generate(tmpl *template.Template, w io.Writer) error {
	edges := g.extractEdges(tmpl)

	fmt.Fprintln(w, "digraph G {")
	fmt.Fprintln(w, "  rankdir=TB;")
	fmt.Fprintln(w, "  node [shape=box];")
	fmt.Fprintln(w)

	// Write resource nodes
	for name, res := range tmpl.Resources {
		label := fmt.Sprintf("%s\\n[%s]", name, res.Type)
		fmt.Fprintf(w, "  %q [label=%q];\n", name, label)
	}

	if g.IncludeParameters {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "  // Parameters")
		for name := range tmpl.Parameters {
			fmt.Fprintf(w, "  %q [shape=ellipse, style=dashed];\n", name)
		}
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "  // Dependencies")
	for _, edge := range edges {
		style := ""
		if edge.Type == "DependsOn" {
			style = " [style=dashed]"
		}
		fmt.Fprintf(w, "  %q -> %q%s;\n", edge.From, edge.To, style)
	}

	fmt.Fprintln(w, "}")
	return nil
}

func (g *Generator) extractEdges(tmpl *template.Template) []Edge {
	var edges []Edge
	seen := make(map[string]bool)

	addEdge := func(from, to, typ string) {
		key := from + "->" + to
		if !seen[key] {
			seen[key] = true
			edges = append(edges, Edge{From: from, To: to, Type: typ})
		}
	}

	// Extract from DependsOn
	for name, res := range tmpl.Resources {
		for _, dep := range res.DependsOn {
			addEdge(name, dep, "DependsOn")
		}
	}

	// Extract Ref and GetAtt from properties
	refPattern := regexp.MustCompile(`\{["']?Ref["']?\s*:\s*["']?(\w+)["']?\}`)
	getAttPattern := regexp.MustCompile(`\{["']?Fn::GetAtt["']?\s*:\s*\[["']?(\w+)["']?`)

	for name, res := range tmpl.Resources {
		propsStr := fmt.Sprintf("%v", res.Properties)

		// Find Ref references
		for _, match := range refPattern.FindAllStringSubmatch(propsStr, -1) {
			target := match[1]
			if tmpl.HasResource(target) {
				addEdge(name, target, "Ref")
			} else if g.IncludeParameters && tmpl.HasParameter(target) {
				addEdge(name, target, "Ref")
			}
		}

		// Find GetAtt references
		for _, match := range getAttPattern.FindAllStringSubmatch(propsStr, -1) {
			target := match[1]
			if tmpl.HasResource(target) {
				addEdge(name, target, "GetAtt")
			}
		}

		// Also check for simple string patterns in properties
		for _, ref := range findRefs(res.Properties) {
			if tmpl.HasResource(ref) {
				addEdge(name, ref, "Ref")
			}
		}
	}

	return edges
}

// findRefs recursively searches for Ref intrinsics in a value.
func findRefs(v any) []string {
	var refs []string

	switch val := v.(type) {
	case map[string]any:
		if ref, ok := val["Ref"].(string); ok {
			refs = append(refs, ref)
		}
		if getAtt, ok := val["Fn::GetAtt"]; ok {
			switch ga := getAtt.(type) {
			case []any:
				if len(ga) > 0 {
					if s, ok := ga[0].(string); ok {
						refs = append(refs, s)
					}
				}
			case string:
				parts := strings.Split(ga, ".")
				if len(parts) > 0 {
					refs = append(refs, parts[0])
				}
			}
		}
		for _, child := range val {
			refs = append(refs, findRefs(child)...)
		}
	case []any:
		for _, child := range val {
			refs = append(refs, findRefs(child)...)
		}
	}

	return refs
}

// GenerateString is a convenience method that returns the DOT graph as a string.
func (g *Generator) GenerateString(tmpl *template.Template) (string, error) {
	var sb strings.Builder
	if err := g.Generate(tmpl, &sb); err != nil {
		return "", err
	}
	return sb.String(), nil
}
