// Package template provides CloudFormation template parsing with line number tracking.
package template

import (
	"fmt"
	"os"
	"strings"

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
	Mappings                 map[string]*Mapping
	Conditions               map[string]*Condition
	Resources                map[string]*Resource
	Outputs                  map[string]*Output
	Metadata                 map[string]any
	Rules                    map[string]*Rule

	// Raw nodes for validation rules that need position info
	MetadataNode   *yaml.Node
	MappingsNode   *yaml.Node
	ConditionsNode *yaml.Node
	RulesNode      *yaml.Node

	// Filename for error reporting.
	Filename string
}

// Mapping represents a CloudFormation mapping.
type Mapping struct {
	Node   *yaml.Node
	Values map[string]map[string]any // TopLevelKey -> SecondLevelKey -> Value
}

// Condition represents a CloudFormation condition.
type Condition struct {
	Node       *yaml.Node
	Expression any // The condition expression (Fn::Equals, Fn::And, etc.)
}

// Parameter represents a CloudFormation parameter.
type Parameter struct {
	Node           *yaml.Node
	Type           string
	Default        any
	AllowedValues  []any
	AllowedPattern string
	MinValue       *float64
	MaxValue       *float64
	MinLength      *int
	MaxLength      *int
	Description    string
	NoEcho         bool
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

// Rule represents a CloudFormation rule for parameter validation.
type Rule struct {
	Node          *yaml.Node
	RuleCondition any             // Optional: determines when the rule takes effect
	Assertions    []RuleAssertion // Required: at least one assertion
}

// RuleAssertion represents an assertion within a CloudFormation rule.
type RuleAssertion struct {
	Node              *yaml.Node
	Assert            any    // Required: condition that must evaluate to true
	AssertDescription string // Optional: error message when assertion fails
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
		Mappings:   make(map[string]*Mapping),
		Conditions: make(map[string]*Condition),
		Resources:  make(map[string]*Resource),
		Outputs:    make(map[string]*Output),
		Metadata:   make(map[string]any),
		Rules:      make(map[string]*Rule),
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
			t.MappingsNode = value
			if err := t.parseMappings(value); err != nil {
				return err
			}
		case "Conditions":
			t.ConditionsNode = value
			if err := t.parseConditions(value); err != nil {
				return err
			}
		case "Metadata":
			t.MetadataNode = value
			if decoded := parseYAMLNode(value); decoded != nil {
				if meta, ok := decoded.(map[string]any); ok {
					t.Metadata = meta
				}
			}
		case "Rules":
			t.RulesNode = value
			if err := t.parseRules(value); err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *Template) parseMappings(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return nil
	}

	for i := 0; i < len(node.Content); i += 2 {
		name := node.Content[i].Value
		mapNode := node.Content[i+1]

		mapping := &Mapping{
			Node:   mapNode,
			Values: make(map[string]map[string]any),
		}

		if mapNode.Kind == yaml.MappingNode {
			for j := 0; j < len(mapNode.Content); j += 2 {
				topKey := mapNode.Content[j].Value
				topVal := mapNode.Content[j+1]

				mapping.Values[topKey] = make(map[string]any)
				if topVal.Kind == yaml.MappingNode {
					for k := 0; k < len(topVal.Content); k += 2 {
						secondKey := topVal.Content[k].Value
						secondVal := parseYAMLNode(topVal.Content[k+1])
						mapping.Values[topKey][secondKey] = secondVal
					}
				}
			}
		}
		t.Mappings[name] = mapping
	}
	return nil
}

func (t *Template) parseConditions(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return nil
	}

	for i := 0; i < len(node.Content); i += 2 {
		name := node.Content[i].Value
		condNode := node.Content[i+1]

		condition := &Condition{
			Node:       condNode,
			Expression: parseYAMLNode(condNode),
		}
		t.Conditions[name] = condition
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
				case "Default":
					param.Default = parseYAMLNode(val)
				case "AllowedValues":
					if val.Kind == yaml.SequenceNode {
						for _, item := range val.Content {
							param.AllowedValues = append(param.AllowedValues, parseYAMLNode(item))
						}
					}
				case "AllowedPattern":
					param.AllowedPattern = val.Value
				case "MinValue":
					var v float64
					if val.Decode(&v) == nil {
						param.MinValue = &v
					}
				case "MaxValue":
					var v float64
					if val.Decode(&v) == nil {
						param.MaxValue = &v
					}
				case "MinLength":
					var v int
					if val.Decode(&v) == nil {
						param.MinLength = &v
					}
				case "MaxLength":
					var v int
					if val.Decode(&v) == nil {
						param.MaxLength = &v
					}
				case "NoEcho":
					var v bool
					if val.Decode(&v) == nil {
						param.NoEcho = v
					}
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
					// Decode properties with CloudFormation tag handling
					if decoded := parseYAMLNode(val); decoded != nil {
						if props, ok := decoded.(map[string]any); ok {
							res.Properties = props
						}
					}
				case "Metadata":
					// Decode metadata
					if decoded := parseYAMLNode(val); decoded != nil {
						if meta, ok := decoded.(map[string]any); ok {
							res.Metadata = meta
						}
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
				case "Value":
					out.Value = parseYAMLNode(val)
				case "Export":
					if decoded := parseYAMLNode(val); decoded != nil {
						if exp, ok := decoded.(map[string]any); ok {
							out.Export = exp
						}
					}
				}
			}
		}
		t.Outputs[name] = out
	}
	return nil
}

func (t *Template) parseRules(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return nil
	}

	for i := 0; i < len(node.Content); i += 2 {
		name := node.Content[i].Value
		ruleNode := node.Content[i+1]

		rule := &Rule{
			Node:       ruleNode,
			Assertions: []RuleAssertion{},
		}

		if ruleNode.Kind == yaml.MappingNode {
			for j := 0; j < len(ruleNode.Content); j += 2 {
				key := ruleNode.Content[j].Value
				val := ruleNode.Content[j+1]

				switch key {
				case "RuleCondition":
					rule.RuleCondition = parseYAMLNode(val)
				case "Assertions":
					if val.Kind == yaml.SequenceNode {
						for _, assertNode := range val.Content {
							assertion := RuleAssertion{Node: assertNode}
							if assertNode.Kind == yaml.MappingNode {
								for k := 0; k < len(assertNode.Content); k += 2 {
									assertKey := assertNode.Content[k].Value
									assertVal := assertNode.Content[k+1]
									switch assertKey {
									case "Assert":
										assertion.Assert = parseYAMLNode(assertVal)
									case "AssertDescription":
										assertion.AssertDescription = assertVal.Value
									}
								}
							}
							rule.Assertions = append(rule.Assertions, assertion)
						}
					}
				}
			}
		}
		t.Rules[name] = rule
	}
	return nil
}

// parseYAMLNode recursively converts a yaml.Node to Go values, handling CF intrinsic tags.
func parseYAMLNode(node *yaml.Node) any {
	if node == nil {
		return nil
	}

	// Handle document node
	if node.Kind == yaml.DocumentNode {
		if len(node.Content) > 0 {
			return parseYAMLNode(node.Content[0])
		}
		return nil
	}

	// Check for CloudFormation intrinsic function tags (single !, not !! standard tags)
	if node.Tag != "" && strings.HasPrefix(node.Tag, "!") && !strings.HasPrefix(node.Tag, "!!") {
		return parseIntrinsicTag(node)
	}

	switch node.Kind {
	case yaml.ScalarNode:
		var val any
		_ = node.Decode(&val)
		return val

	case yaml.SequenceNode:
		result := make([]any, len(node.Content))
		for i, child := range node.Content {
			result[i] = parseYAMLNode(child)
		}
		return result

	case yaml.MappingNode:
		result := make(map[string]any)
		for i := 0; i < len(node.Content); i += 2 {
			key := node.Content[i].Value
			result[key] = parseYAMLNode(node.Content[i+1])
		}
		return result

	case yaml.AliasNode:
		return parseYAMLNode(node.Alias)
	}

	return nil
}

// parseNodeContents parses the contents of a tagged node without re-checking the tag.
// This prevents infinite recursion when an intrinsic wraps another structure.
func parseNodeContents(node *yaml.Node) any {
	switch node.Kind {
	case yaml.ScalarNode:
		var val any
		_ = node.Decode(&val)
		return val
	case yaml.SequenceNode:
		result := make([]any, len(node.Content))
		for i, child := range node.Content {
			result[i] = parseYAMLNode(child)
		}
		return result
	case yaml.MappingNode:
		result := make(map[string]any)
		for i := 0; i < len(node.Content); i += 2 {
			key := node.Content[i].Value
			result[key] = parseYAMLNode(node.Content[i+1])
		}
		return result
	}
	return nil
}

// parseIntrinsicTag handles CloudFormation intrinsic function YAML tags.
// Returns map[string]any in the long-form CloudFormation format.
func parseIntrinsicTag(node *yaml.Node) map[string]any {
	tag := strings.TrimPrefix(node.Tag, "!")

	switch tag {
	case "Ref":
		return map[string]any{"Ref": node.Value}

	case "GetAtt":
		if node.Kind == yaml.ScalarNode {
			// !GetAtt Resource.Attribute format
			return map[string]any{"Fn::GetAtt": node.Value}
		}
		// !GetAtt [Resource, Attribute] format
		if node.Kind == yaml.SequenceNode {
			parts := make([]any, len(node.Content))
			for i, child := range node.Content {
				parts[i] = child.Value
			}
			return map[string]any{"Fn::GetAtt": parts}
		}

	case "Sub":
		if node.Kind == yaml.ScalarNode {
			return map[string]any{"Fn::Sub": node.Value}
		}
		if node.Kind == yaml.SequenceNode {
			args := make([]any, len(node.Content))
			for i, child := range node.Content {
				args[i] = parseYAMLNode(child)
			}
			return map[string]any{"Fn::Sub": args}
		}

	case "Join":
		if node.Kind == yaml.SequenceNode && len(node.Content) >= 2 {
			args := make([]any, len(node.Content))
			for i, child := range node.Content {
				args[i] = parseYAMLNode(child)
			}
			return map[string]any{"Fn::Join": args}
		}

	case "Select":
		if node.Kind == yaml.SequenceNode && len(node.Content) >= 2 {
			args := make([]any, len(node.Content))
			for i, child := range node.Content {
				args[i] = parseYAMLNode(child)
			}
			return map[string]any{"Fn::Select": args}
		}

	case "If":
		if node.Kind == yaml.SequenceNode && len(node.Content) >= 3 {
			args := make([]any, len(node.Content))
			for i, child := range node.Content {
				args[i] = parseYAMLNode(child)
			}
			return map[string]any{"Fn::If": args}
		}

	case "Condition":
		return map[string]any{"Condition": node.Value}

	case "GetAZs":
		if node.Kind == yaml.ScalarNode {
			return map[string]any{"Fn::GetAZs": node.Value}
		}
		// Handle nested intrinsic
		return map[string]any{"Fn::GetAZs": parseNodeContents(node)}

	case "Base64":
		if node.Kind == yaml.ScalarNode {
			return map[string]any{"Fn::Base64": node.Value}
		}
		// For non-scalar (nested intrinsics), parse contents directly
		return map[string]any{"Fn::Base64": parseNodeContents(node)}

	case "Cidr":
		if node.Kind == yaml.SequenceNode && len(node.Content) >= 3 {
			args := make([]any, len(node.Content))
			for i, child := range node.Content {
				args[i] = parseYAMLNode(child)
			}
			return map[string]any{"Fn::Cidr": args}
		}

	case "ImportValue":
		if node.Kind == yaml.ScalarNode {
			return map[string]any{"Fn::ImportValue": node.Value}
		}
		return map[string]any{"Fn::ImportValue": parseNodeContents(node)}

	case "Split":
		if node.Kind == yaml.SequenceNode && len(node.Content) >= 2 {
			args := make([]any, len(node.Content))
			for i, child := range node.Content {
				args[i] = parseYAMLNode(child)
			}
			return map[string]any{"Fn::Split": args}
		}

	case "FindInMap":
		if node.Kind == yaml.SequenceNode && len(node.Content) >= 3 {
			args := make([]any, len(node.Content))
			for i, child := range node.Content {
				args[i] = parseYAMLNode(child)
			}
			return map[string]any{"Fn::FindInMap": args}
		}

	case "Equals":
		if node.Kind == yaml.SequenceNode && len(node.Content) >= 2 {
			args := make([]any, len(node.Content))
			for i, child := range node.Content {
				args[i] = parseYAMLNode(child)
			}
			return map[string]any{"Fn::Equals": args}
		}

	case "And":
		if node.Kind == yaml.SequenceNode {
			args := make([]any, len(node.Content))
			for i, child := range node.Content {
				args[i] = parseYAMLNode(child)
			}
			return map[string]any{"Fn::And": args}
		}

	case "Or":
		if node.Kind == yaml.SequenceNode {
			args := make([]any, len(node.Content))
			for i, child := range node.Content {
				args[i] = parseYAMLNode(child)
			}
			return map[string]any{"Fn::Or": args}
		}

	case "Not":
		if node.Kind == yaml.SequenceNode && len(node.Content) > 0 {
			return map[string]any{"Fn::Not": []any{parseYAMLNode(node.Content[0])}}
		}

	case "Transform":
		return map[string]any{"Fn::Transform": parseNodeContents(node)}
	}

	// Unknown tag - generic Fn:: handler
	fnName := "Fn::" + tag
	if node.Kind == yaml.ScalarNode {
		return map[string]any{fnName: node.Value}
	}
	return map[string]any{fnName: parseNodeContents(node)}
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
