// Package graph generates DOT and Mermaid format dependency graphs from CloudFormation templates.
package graph

import (
	"io"
	"strings"

	"github.com/emicklei/dot"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

// Format specifies the output format for the graph.
type Format string

const (
	// FormatDOT outputs Graphviz DOT format.
	FormatDOT Format = "dot"
	// FormatMermaid outputs Mermaid format for GitHub/markdown rendering.
	FormatMermaid Format = "mermaid"
)

// Generator creates dependency graphs from CloudFormation templates.
type Generator struct {
	// IncludeParameters includes parameter references in the graph.
	IncludeParameters bool

	// Format specifies the output format (dot or mermaid). Defaults to dot.
	Format Format

	// ClusterByType groups resources by AWS service type.
	ClusterByType bool
}

// Edge represents a dependency between resources.
type Edge struct {
	From string
	To   string
	Type string // "Ref", "GetAtt", "DependsOn"
}

// Generate creates a dependency graph from a template and writes it to w.
func (g *Generator) Generate(tmpl *template.Template, w io.Writer) error {
	graph := g.buildGraph(tmpl)

	format := g.Format
	if format == "" {
		format = FormatDOT
	}

	var output string
	if format == FormatMermaid {
		output = dot.MermaidGraph(graph, dot.MermaidTopToBottom)
	} else {
		output = graph.String()
	}

	_, err := w.Write([]byte(output))
	return err
}

// GenerateString is a convenience method that returns the graph as a string.
func (g *Generator) GenerateString(tmpl *template.Template) (string, error) {
	var sb strings.Builder
	if err := g.Generate(tmpl, &sb); err != nil {
		return "", err
	}
	return sb.String(), nil
}

// buildGraph creates the dot.Graph structure from the template.
func (g *Generator) buildGraph(tmpl *template.Template) *dot.Graph {
	graph := dot.NewGraph(dot.Directed)
	graph.Attr("rankdir", "TB")

	// Set default node style using NodeInitializer
	graph.NodeInitializer(func(n dot.Node) {
		n.Attr("shape", "box")
		n.Attr("fontname", "Arial")
	})

	// Set default edge style using EdgeInitializer
	graph.EdgeInitializer(func(e dot.Edge) {
		e.Attr("fontname", "Arial")
		e.Attr("fontsize", "10")
	})

	edges := g.extractEdges(tmpl)

	if g.ClusterByType {
		g.addClusteredNodes(graph, tmpl)
	} else {
		g.addNodes(graph, tmpl)
	}

	// Add parameter nodes if requested
	if g.IncludeParameters {
		for name := range tmpl.Parameters {
			n := graph.Node(name)
			n.Attr("shape", "ellipse")
			n.Attr("style", "dashed")
			n.Label(name)
		}
	}

	// Add edges with appropriate styles
	for _, edge := range edges {
		from := graph.Node(edge.From)
		to := graph.Node(edge.To)
		e := graph.Edge(from, to)

		switch edge.Type {
		case "DependsOn":
			e.Dashed()
			e.Attr("color", "gray")
		case "GetAtt":
			e.Attr("color", "blue")
		default: // Ref
			e.Solid()
		}
	}

	return graph
}

// addNodes adds resource nodes without clustering.
func (g *Generator) addNodes(graph *dot.Graph, tmpl *template.Template) {
	for name, res := range tmpl.Resources {
		n := graph.Node(name)
		// Use HTML-like label for better formatting
		n.Label(name + "\\n[" + res.Type + "]")
	}
}

// addClusteredNodes adds resource nodes grouped by AWS service type.
func (g *Generator) addClusteredNodes(graph *dot.Graph, tmpl *template.Template) {
	// Group resources by service
	serviceResources := make(map[string][]string)
	resourceTypes := make(map[string]string)

	for name, res := range tmpl.Resources {
		service := extractService(res.Type)
		serviceResources[service] = append(serviceResources[service], name)
		resourceTypes[name] = res.Type
	}

	// Create clusters for each service with multiple resources
	for service, resources := range serviceResources {
		if len(resources) > 1 {
			cluster := graph.Subgraph("cluster_"+service, dot.ClusterOption{})
			cluster.Attr("label", service)
			cluster.Attr("style", "rounded")
			cluster.Attr("bgcolor", "lightyellow")

			for _, name := range resources {
				n := cluster.Node(name)
				n.Label(name + "\\n[" + resourceTypes[name] + "]")
			}
		} else {
			// Single resource, no cluster needed
			for _, name := range resources {
				n := graph.Node(name)
				n.Label(name + "\\n[" + resourceTypes[name] + "]")
			}
		}
	}
}

// extractService extracts the AWS service name from a resource type.
// e.g., "AWS::S3::Bucket" -> "S3"
func extractService(resourceType string) string {
	parts := strings.Split(resourceType, "::")
	if len(parts) >= 2 {
		return parts[1]
	}
	return "Other"
}

// extractEdges extracts all dependency edges from the template.
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
	for name, res := range tmpl.Resources {
		refs := findRefs(res.Properties)
		for _, ref := range refs.refs {
			if tmpl.HasResource(ref) {
				addEdge(name, ref, "Ref")
			} else if g.IncludeParameters && tmpl.HasParameter(ref) {
				addEdge(name, ref, "Ref")
			}
		}
		for _, ref := range refs.getAtts {
			if tmpl.HasResource(ref) {
				addEdge(name, ref, "GetAtt")
			}
		}
	}

	return edges
}

// refResults holds the results of findRefs
type refResults struct {
	refs    []string
	getAtts []string
}

// findRefs recursively searches for Ref and GetAtt intrinsics in a value.
func findRefs(v any) refResults {
	var result refResults
	findRefsRecursive(v, &result)
	return result
}

func findRefsRecursive(v any, result *refResults) {
	switch val := v.(type) {
	case map[string]any:
		if ref, ok := val["Ref"].(string); ok {
			result.refs = append(result.refs, ref)
		}
		if getAtt, ok := val["Fn::GetAtt"]; ok {
			switch ga := getAtt.(type) {
			case []any:
				if len(ga) > 0 {
					if s, ok := ga[0].(string); ok {
						result.getAtts = append(result.getAtts, s)
					}
				}
			case string:
				parts := strings.Split(ga, ".")
				if len(parts) > 0 {
					result.getAtts = append(result.getAtts, parts[0])
				}
			}
		}
		for _, child := range val {
			findRefsRecursive(child, result)
		}
	case []any:
		for _, child := range val {
			findRefsRecursive(child, result)
		}
	}
}
