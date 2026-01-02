// Package graph generates DOT and Mermaid format dependency graphs from CloudFormation templates.
//
// This package extracts resource dependencies (Ref, GetAtt, DependsOn) and generates
// graph output for visualization. It uses the emicklei/dot library for DOT generation.
//
// # Output Formats
//
// Two output formats are supported:
//
//   - DOT (default): Graphviz DOT format, render with `dot -Tpng`
//   - Mermaid: Renders natively in GitHub READMEs, VS Code, etc.
//
// # Basic Usage
//
// Generate a DOT graph to stdout:
//
//	tmpl, _ := template.ParseFile("template.yaml")
//	gen := &graph.Generator{}
//	gen.Generate(tmpl, os.Stdout)
//
// Generate to a string:
//
//	dot, err := gen.GenerateString(tmpl)
//
// # Mermaid Output
//
// Generate Mermaid format for embedding in markdown:
//
//	gen := &graph.Generator{
//	    Format: graph.FormatMermaid,
//	}
//	mermaid, _ := gen.GenerateString(tmpl)
//
// The output can be embedded directly in GitHub markdown:
//
//	```mermaid
//	graph TB
//	  MyBucket[MyBucket<br/>AWS::S3::Bucket]
//	  MyRole[MyRole<br/>AWS::IAM::Role]
//	  MyRole --> MyBucket
//	```
//
// # Including Parameters
//
// By default, only resource dependencies are shown. To include parameter references:
//
//	gen := &graph.Generator{
//	    IncludeParameters: true,
//	}
//
// # Clustering by Service
//
// Group resources by AWS service type (e.g., all S3 buckets together):
//
//	gen := &graph.Generator{
//	    ClusterByType: true,
//	}
//
// # DOT Output Format
//
// The DOT output is valid Graphviz format:
//
//	digraph {
//	  rankdir=TB;
//	  n1[label="MyBucket\n[AWS::S3::Bucket]",shape="box"];
//	  n2[label="MyRole\n[AWS::IAM::Role]",shape="box"];
//	  n2->n1[style="dashed"];
//	}
//
// Render with Graphviz:
//
//	cfn-lint graph template.yaml | dot -Tpng -o graph.png
//
// # Edge Styles
//
// Dependencies are styled and colored by type:
//   - Solid black lines: Ref references
//   - Solid blue lines: GetAtt references
//   - Dashed gray lines: DependsOn explicit dependencies
//
// # Library Features Used
//
// This package leverages emicklei/dot features including:
//   - NodeInitializer: Sets default node styles (box shape, Arial font)
//   - EdgeInitializer: Sets default edge styles (Arial font, size 10)
//   - Subgraphs: Creates clusters for service grouping
//   - MermaidGraph: Built-in Mermaid format export
package graph
