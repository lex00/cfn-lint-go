// Package template provides CloudFormation template parsing with line number tracking.
package template

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Template represents a parsed CloudFormation template with position information.
type Template struct {
	// Root is the raw parsed YAML node tree (preserves line numbers).
	Root *yaml.Node

	// Parsed sections for convenient access.
	AWSTemplateFormatVersion string
	Description              string
	Parameters               map[string]*Parameter
	Mappings                 map[string]any
	Conditions               map[string]any
	Resources                map[string]*Resource
	Outputs                  map[string]*Output

	// Filename for error reporting.
	Filename string
}

// Parameter represents a CloudFormation parameter.
type Parameter struct {
	Node          *yaml.Node
	Type          string
	Default       any
	AllowedValues []any
	Description   string
}

// Resource represents a CloudFormation resource.
type Resource struct {
	Node       *yaml.Node
	Type       string
	Properties map[string]any
	DependsOn  []string
	Condition  string
	Metadata   map[string]any
}

// Output represents a CloudFormation output.
type Output struct {
	Node        *yaml.Node
	Value       any
	Description string
	Export      map[string]any
	Condition   string
}

// ParseFile parses a CloudFormation template from a file.
func ParseFile(path string) (*Template, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	tmpl, err := Parse(data)
	if err != nil {
		return nil, err
	}
	tmpl.Filename = path
	return tmpl, nil
}

// Parse parses a CloudFormation template from bytes.
func Parse(data []byte) (*Template, error) {
	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}

	tmpl := &Template{
		Root:       &root,
		Parameters: make(map[string]*Parameter),
		Resources:  make(map[string]*Resource),
		Outputs:    make(map[string]*Output),
	}

	if err := tmpl.parseRoot(); err != nil {
		return nil, err
	}

	return tmpl, nil
}

func (t *Template) parseRoot() error {
	if t.Root.Kind != yaml.DocumentNode || len(t.Root.Content) == 0 {
		return fmt.Errorf("expected document node")
	}

	doc := t.Root.Content[0]
	if doc.Kind != yaml.MappingNode {
		return fmt.Errorf("expected mapping at document root")
	}

	for i := 0; i < len(doc.Content); i += 2 {
		key := doc.Content[i]
		value := doc.Content[i+1]

		switch key.Value {
		case "AWSTemplateFormatVersion":
			t.AWSTemplateFormatVersion = value.Value
		case "Description":
			t.Description = value.Value
		case "Parameters":
			if err := t.parseParameters(value); err != nil {
				return err
			}
		case "Resources":
			if err := t.parseResources(value); err != nil {
				return err
			}
		case "Outputs":
			if err := t.parseOutputs(value); err != nil {
				return err
			}
		case "Mappings":
			// TODO: Parse mappings
		case "Conditions":
			// TODO: Parse conditions
		}
	}

	return nil
}

func (t *Template) parseParameters(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return nil
	}

	for i := 0; i < len(node.Content); i += 2 {
		name := node.Content[i].Value
		paramNode := node.Content[i+1]

		param := &Parameter{Node: paramNode}
		if paramNode.Kind == yaml.MappingNode {
			for j := 0; j < len(paramNode.Content); j += 2 {
				key := paramNode.Content[j].Value
				val := paramNode.Content[j+1]
				switch key {
				case "Type":
					param.Type = val.Value
				case "Description":
					param.Description = val.Value
				}
			}
		}
		t.Parameters[name] = param
	}
	return nil
}

func (t *Template) parseResources(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return nil
	}

	for i := 0; i < len(node.Content); i += 2 {
		name := node.Content[i].Value
		resNode := node.Content[i+1]

		res := &Resource{
			Node:       resNode,
			Properties: make(map[string]any),
		}

		if resNode.Kind == yaml.MappingNode {
			for j := 0; j < len(resNode.Content); j += 2 {
				key := resNode.Content[j].Value
				val := resNode.Content[j+1]
				switch key {
				case "Type":
					res.Type = val.Value
				case "Properties":
					// Store raw properties for now
					var props map[string]any
					if err := val.Decode(&props); err == nil {
						res.Properties = props
					}
				case "Condition":
					res.Condition = val.Value
				case "DependsOn":
					var deps []string
					if val.Kind == yaml.SequenceNode {
						for _, d := range val.Content {
							deps = append(deps, d.Value)
						}
					} else {
						deps = []string{val.Value}
					}
					res.DependsOn = deps
				}
			}
		}
		t.Resources[name] = res
	}
	return nil
}

func (t *Template) parseOutputs(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return nil
	}

	for i := 0; i < len(node.Content); i += 2 {
		name := node.Content[i].Value
		outNode := node.Content[i+1]

		out := &Output{Node: outNode}
		if outNode.Kind == yaml.MappingNode {
			for j := 0; j < len(outNode.Content); j += 2 {
				key := outNode.Content[j].Value
				val := outNode.Content[j+1]
				switch key {
				case "Description":
					out.Description = val.Value
				case "Condition":
					out.Condition = val.Value
				}
			}
		}
		t.Outputs[name] = out
	}
	return nil
}

// GetResourceNames returns all resource logical IDs.
func (t *Template) GetResourceNames() []string {
	names := make([]string, 0, len(t.Resources))
	for name := range t.Resources {
		names = append(names, name)
	}
	return names
}

// GetParameterNames returns all parameter names.
func (t *Template) GetParameterNames() []string {
	names := make([]string, 0, len(t.Parameters))
	for name := range t.Parameters {
		names = append(names, name)
	}
	return names
}

// HasResource checks if a resource with the given logical ID exists.
func (t *Template) HasResource(name string) bool {
	_, ok := t.Resources[name]
	return ok
}

// HasParameter checks if a parameter with the given name exists.
func (t *Template) HasParameter(name string) bool {
	_, ok := t.Parameters[name]
	return ok
}
