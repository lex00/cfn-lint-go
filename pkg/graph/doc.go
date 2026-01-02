// Package graph generates DOT format dependency graphs from CloudFormation templates.
//
// This package extracts resource dependencies (Ref, GetAtt, DependsOn) and generates
// Graphviz DOT format output for visualization.
//
// # Basic Usage
//
// Generate a graph to stdout:
//
//	tmpl, _ := template.ParseFile("template.yaml")
//	gen := &graph.Generator{}
//	gen.Generate(tmpl, os.Stdout)
//
// Generate to a string:
//
//	dot, err := gen.GenerateString(tmpl)
//
// # Including Parameters
//
// By default, only resource dependencies are shown. To include parameter references:
//
//	gen := &graph.Generator{
//	    IncludeParameters: true,
//	}
//
// # Output Format
//
// The output is valid Graphviz DOT format:
//
//	digraph G {
//	  rankdir=TB;
//	  node [shape=box];
//
//	  "MyBucket" [label="MyBucket\n[AWS::S3::Bucket]"];
//	  "MyRole" [label="MyRole\n[AWS::IAM::Role]"];
//
//	  // Parameters
//	  "Environment" [shape=ellipse, style=dashed];
//
//	  // Dependencies
//	  "MyRole" -> "MyBucket" [style=dashed];
//	}
//
// Render with Graphviz:
//
//	cfn-lint graph template.yaml | dot -Tpng -o graph.png
package graph
