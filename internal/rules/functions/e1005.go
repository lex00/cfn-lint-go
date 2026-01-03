// Package functions contains intrinsic function and template validation rules (E1xxx).
package functions

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
	"gopkg.in/yaml.v3"
)

func init() {
	rules.Register(&E1005{})
}

// E1005 checks that Transform declarations are valid.
type E1005 struct{}

func (r *E1005) ID() string { return "E1005" }

func (r *E1005) ShortDesc() string {
	return "Transform configuration error"
}

func (r *E1005) Description() string {
	return "Checks that template Transform declarations use valid transform names and configuration."
}

func (r *E1005) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E1005"
}

func (r *E1005) Tags() []string {
	return []string{"template", "transform"}
}

// Known AWS transforms
var validTransforms = map[string]bool{
	"AWS::Serverless-2016-10-31":     true,
	"AWS::Include":                   true,
	"AWS::CodeDeployBlueGreen":       true,
	"AWS::SecretsManager-2020-07-23": true,
	"AWS::LanguageExtensions":        true,
}

func (r *E1005) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Find Transform in root
	if tmpl.Root == nil || tmpl.Root.Kind != yaml.DocumentNode || len(tmpl.Root.Content) == 0 {
		return nil
	}

	doc := tmpl.Root.Content[0]
	if doc.Kind != yaml.MappingNode {
		return nil
	}

	for i := 0; i < len(doc.Content); i += 2 {
		key := doc.Content[i]
		value := doc.Content[i+1]

		if key.Value == "Transform" {
			matches = append(matches, r.checkTransform(value)...)
		}
	}

	return matches
}

func (r *E1005) checkTransform(node *yaml.Node) []rules.Match {
	var matches []rules.Match

	switch node.Kind {
	case yaml.ScalarNode:
		// Single transform
		if !isValidTransform(node.Value) {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Unknown transform '%s'", node.Value),
				Line:    node.Line,
				Column:  node.Column,
				Path:    []string{"Transform"},
			})
		}

	case yaml.SequenceNode:
		// List of transforms
		for idx, item := range node.Content {
			switch item.Kind {
			case yaml.ScalarNode:
				if !isValidTransform(item.Value) {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Unknown transform '%s'", item.Value),
						Line:    item.Line,
						Column:  item.Column,
						Path:    []string{"Transform", fmt.Sprintf("[%d]", idx)},
					})
				}
			case yaml.MappingNode:
				// Inline transform with Name key
				matches = append(matches, r.checkInlineTransform(item, idx)...)
			}
		}

	case yaml.MappingNode:
		// Inline transform with Name key
		matches = append(matches, r.checkInlineTransform(node, -1)...)
	}

	return matches
}

func (r *E1005) checkInlineTransform(node *yaml.Node, idx int) []rules.Match {
	var matches []rules.Match
	var name string
	var nameNode *yaml.Node

	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		value := node.Content[i+1]

		if key.Value == "Name" {
			name = value.Value
			nameNode = value
		}
	}

	if name == "" {
		path := []string{"Transform"}
		if idx >= 0 {
			path = append(path, fmt.Sprintf("[%d]", idx))
		}
		matches = append(matches, rules.Match{
			Message: "Transform is missing required 'Name' property",
			Line:    node.Line,
			Column:  node.Column,
			Path:    path,
		})
	} else if !isValidTransform(name) {
		path := []string{"Transform"}
		if idx >= 0 {
			path = append(path, fmt.Sprintf("[%d]", idx))
		}
		path = append(path, "Name")
		matches = append(matches, rules.Match{
			Message: fmt.Sprintf("Unknown transform '%s'", name),
			Line:    nameNode.Line,
			Column:  nameNode.Column,
			Path:    path,
		})
	}

	return matches
}

func isValidTransform(name string) bool {
	// Check known transforms
	if validTransforms[name] {
		return true
	}
	// Allow macro transforms (AWS::CloudFormation::Macro::*)
	if strings.HasPrefix(name, "AWS::CloudFormation::Macro::") {
		return true
	}
	// Allow custom macro names (alphanumeric with dashes/underscores)
	// Macros can have custom names when registered
	return false
}
