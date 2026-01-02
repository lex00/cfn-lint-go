package graph

import (
	"strings"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

func TestGenerator_Generate(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket

  MyRole:
    Type: AWS::IAM::Role
    DependsOn: MyBucket

  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Role: !GetAtt MyRole.Arn
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	gen := &Generator{}
	dot, err := gen.GenerateString(tmpl)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Check basic DOT structure
	if !strings.Contains(dot, "digraph") {
		t.Error("Expected 'digraph' in output")
	}

	// Check for resource labels (emicklei/dot uses internal node IDs but our labels)
	if !strings.Contains(dot, "MyBucket") {
		t.Error("Expected MyBucket in output")
	}
	if !strings.Contains(dot, "MyRole") {
		t.Error("Expected MyRole in output")
	}
	if !strings.Contains(dot, "MyFunction") {
		t.Error("Expected MyFunction in output")
	}

	// Check for resource type labels
	if !strings.Contains(dot, "AWS::S3::Bucket") {
		t.Error("Expected AWS::S3::Bucket type in output")
	}
	if !strings.Contains(dot, "AWS::IAM::Role") {
		t.Error("Expected AWS::IAM::Role type in output")
	}

	// Check for edges (arrow syntax)
	if !strings.Contains(dot, "->") {
		t.Error("Expected edges (arrows) in output")
	}

	// Check for edge coloring (GetAtt should be blue)
	if !strings.Contains(dot, "blue") {
		t.Error("Expected blue color for GetAtt edge")
	}
}

func TestGenerator_WithParameters(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  BucketName:
    Type: String

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !Ref BucketName
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	gen := &Generator{IncludeParameters: true}
	dot, err := gen.GenerateString(tmpl)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Check for parameter node with ellipse shape
	if !strings.Contains(dot, "BucketName") {
		t.Error("Expected BucketName parameter in output")
	}
	if !strings.Contains(dot, "ellipse") {
		t.Error("Expected parameter to have ellipse shape")
	}
}

func TestGenerator_Mermaid(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket

  MyRole:
    Type: AWS::IAM::Role
    DependsOn: MyBucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	gen := &Generator{Format: FormatMermaid}
	out, err := gen.GenerateString(tmpl)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Check Mermaid structure (library uses "graph TB")
	if !strings.Contains(out, "graph") {
		t.Error("Expected 'graph' in Mermaid output")
	}

	// Check for resource nodes
	if !strings.Contains(out, "MyBucket") {
		t.Error("Expected MyBucket in output")
	}
	if !strings.Contains(out, "MyRole") {
		t.Error("Expected MyRole in output")
	}

	// Check for edges
	if !strings.Contains(out, "-->") && !strings.Contains(out, "-.->") {
		t.Error("Expected edge in Mermaid output")
	}
}

func TestGenerator_MermaidWithParameters(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Parameters:
  EnvName:
    Type: String

Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !Ref EnvName
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	gen := &Generator{Format: FormatMermaid, IncludeParameters: true}
	out, err := gen.GenerateString(tmpl)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Check for parameter node
	if !strings.Contains(out, "EnvName") {
		t.Error("Expected EnvName parameter in Mermaid output")
	}

	// Check for edge
	if !strings.Contains(out, "-->") {
		t.Error("Expected edge in Mermaid output")
	}
}

func TestGenerator_ClusterByType(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  Bucket1:
    Type: AWS::S3::Bucket
  Bucket2:
    Type: AWS::S3::Bucket
  MyRole:
    Type: AWS::IAM::Role
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	gen := &Generator{ClusterByType: true}
	dot, err := gen.GenerateString(tmpl)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Check for S3 cluster (since there are 2 S3 buckets)
	// emicklei/dot uses internal subgraph IDs, but the label should be "S3"
	if !strings.Contains(dot, `label="S3"`) {
		t.Error("Expected S3 cluster label in output")
	}
	if !strings.Contains(dot, "subgraph cluster_") {
		t.Error("Expected cluster subgraph in output")
	}

	// Check for cluster background color
	if !strings.Contains(dot, "lightyellow") {
		t.Error("Expected lightyellow background for cluster")
	}

	// IAM has only one resource, so no cluster expected
	// Just check that the role is present
	if !strings.Contains(dot, "AWS::IAM::Role") {
		t.Error("Expected IAM role in output")
	}
}

func TestFindRefs(t *testing.T) {
	tests := []struct {
		name            string
		input           any
		expectedRefs    []string
		expectedGetAtts []string
	}{
		{
			name:            "simple ref",
			input:           map[string]any{"Ref": "MyResource"},
			expectedRefs:    []string{"MyResource"},
			expectedGetAtts: nil,
		},
		{
			name: "getatt array",
			input: map[string]any{
				"Fn::GetAtt": []any{"MyResource", "Arn"},
			},
			expectedRefs:    nil,
			expectedGetAtts: []string{"MyResource"},
		},
		{
			name: "getatt string",
			input: map[string]any{
				"Fn::GetAtt": "MyResource.Arn",
			},
			expectedRefs:    nil,
			expectedGetAtts: []string{"MyResource"},
		},
		{
			name: "nested ref",
			input: map[string]any{
				"Key": map[string]any{"Ref": "NestedRef"},
			},
			expectedRefs:    []string{"NestedRef"},
			expectedGetAtts: nil,
		},
		{
			name: "refs in array",
			input: []any{
				map[string]any{"Ref": "First"},
				map[string]any{"Ref": "Second"},
			},
			expectedRefs:    []string{"First", "Second"},
			expectedGetAtts: nil,
		},
		{
			name: "mixed ref and getatt",
			input: map[string]any{
				"RoleArn": map[string]any{"Fn::GetAtt": []any{"MyRole", "Arn"}},
				"Bucket":  map[string]any{"Ref": "MyBucket"},
			},
			expectedRefs:    []string{"MyBucket"},
			expectedGetAtts: []string{"MyRole"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := findRefs(tc.input)

			if len(result.refs) != len(tc.expectedRefs) {
				t.Errorf("Expected %d refs, got %d", len(tc.expectedRefs), len(result.refs))
			} else {
				for i, exp := range tc.expectedRefs {
					if result.refs[i] != exp {
						t.Errorf("Expected ref[%d] = %q, got %q", i, exp, result.refs[i])
					}
				}
			}

			if len(result.getAtts) != len(tc.expectedGetAtts) {
				t.Errorf("Expected %d getAtts, got %d", len(tc.expectedGetAtts), len(result.getAtts))
			} else {
				for i, exp := range tc.expectedGetAtts {
					if result.getAtts[i] != exp {
						t.Errorf("Expected getAtt[%d] = %q, got %q", i, exp, result.getAtts[i])
					}
				}
			}
		})
	}
}

func TestExtractEdges(t *testing.T) {
	yaml := `
Resources:
  A:
    Type: AWS::S3::Bucket
  B:
    Type: AWS::S3::Bucket
    DependsOn: A
  C:
    Type: AWS::S3::Bucket
    DependsOn:
      - A
      - B
`
	tmpl, _ := template.Parse([]byte(yaml))
	gen := &Generator{}
	edges := gen.extractEdges(tmpl)

	// Should have 3 edges: B->A, C->A, C->B
	if len(edges) != 3 {
		t.Errorf("Expected 3 edges, got %d", len(edges))
	}

	// Check edge types
	for _, e := range edges {
		if e.Type != "DependsOn" {
			t.Errorf("Expected DependsOn type, got %s", e.Type)
		}
	}
}

func TestExtractEdgesWithGetAtt(t *testing.T) {
	yaml := `
Resources:
  MyRole:
    Type: AWS::IAM::Role
  MyFunction:
    Type: AWS::Lambda::Function
    Properties:
      Role: !GetAtt MyRole.Arn
`
	tmpl, _ := template.Parse([]byte(yaml))
	gen := &Generator{}
	edges := gen.extractEdges(tmpl)

	// Should have 1 edge: MyFunction->MyRole (GetAtt)
	if len(edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(edges))
		return
	}

	if edges[0].Type != "GetAtt" {
		t.Errorf("Expected GetAtt type, got %s", edges[0].Type)
	}
	if edges[0].From != "MyFunction" {
		t.Errorf("Expected edge from MyFunction, got %s", edges[0].From)
	}
	if edges[0].To != "MyRole" {
		t.Errorf("Expected edge to MyRole, got %s", edges[0].To)
	}
}

func TestExtractService(t *testing.T) {
	tests := []struct {
		resourceType string
		expected     string
	}{
		{"AWS::S3::Bucket", "S3"},
		{"AWS::Lambda::Function", "Lambda"},
		{"AWS::IAM::Role", "IAM"},
		{"AWS::EC2::Instance", "EC2"},
		{"Custom::MyResource", "MyResource"}, // parts[1] is "MyResource"
		{"InvalidType", "Other"},
	}

	for _, tc := range tests {
		t.Run(tc.resourceType, func(t *testing.T) {
			result := extractService(tc.resourceType)
			if result != tc.expected {
				t.Errorf("extractService(%q) = %q, want %q", tc.resourceType, result, tc.expected)
			}
		})
	}
}

func TestEdgeInitializer(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
  MyRole:
    Type: AWS::IAM::Role
    DependsOn: MyBucket
`
	tmpl, err := template.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	gen := &Generator{}
	dot, err := gen.GenerateString(tmpl)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Check that EdgeInitializer set default font
	if !strings.Contains(dot, "Arial") {
		t.Error("Expected Arial font from EdgeInitializer")
	}
}
