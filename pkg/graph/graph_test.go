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
	if !strings.Contains(dot, "digraph G {") {
		t.Error("Expected 'digraph G {' in output")
	}

	// Check for resource nodes
	if !strings.Contains(dot, `"MyBucket"`) {
		t.Error("Expected MyBucket node in output")
	}
	if !strings.Contains(dot, `"MyRole"`) {
		t.Error("Expected MyRole node in output")
	}
	if !strings.Contains(dot, `"MyFunction"`) {
		t.Error("Expected MyFunction node in output")
	}

	// Check for DependsOn edge
	if !strings.Contains(dot, `"MyRole" -> "MyBucket"`) {
		t.Error("Expected DependsOn edge from MyRole to MyBucket")
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
	if !strings.Contains(dot, `"BucketName"`) {
		t.Error("Expected BucketName parameter in output")
	}
	if !strings.Contains(dot, "shape=ellipse") {
		t.Error("Expected parameter to have ellipse shape")
	}
}

func TestFindRefs(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected []string
	}{
		{
			name:     "simple ref",
			input:    map[string]any{"Ref": "MyResource"},
			expected: []string{"MyResource"},
		},
		{
			name: "getatt array",
			input: map[string]any{
				"Fn::GetAtt": []any{"MyResource", "Arn"},
			},
			expected: []string{"MyResource"},
		},
		{
			name: "getatt string",
			input: map[string]any{
				"Fn::GetAtt": "MyResource.Arn",
			},
			expected: []string{"MyResource"},
		},
		{
			name: "nested",
			input: map[string]any{
				"Key": map[string]any{"Ref": "NestedRef"},
			},
			expected: []string{"NestedRef"},
		},
		{
			name: "in array",
			input: []any{
				map[string]any{"Ref": "First"},
				map[string]any{"Ref": "Second"},
			},
			expected: []string{"First", "Second"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			refs := findRefs(tc.input)
			if len(refs) != len(tc.expected) {
				t.Errorf("Expected %d refs, got %d", len(tc.expected), len(refs))
				return
			}
			for i, exp := range tc.expected {
				if refs[i] != exp {
					t.Errorf("Expected ref[%d] = %q, got %q", i, exp, refs[i])
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
