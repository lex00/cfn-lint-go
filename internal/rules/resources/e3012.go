// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/schema"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3012{})
}

// E3012 checks that resource property values match expected types.
// Uses CloudFormation resource schemas from cloudformation-schema-go.
type E3012 struct{}

func (r *E3012) ID() string { return "E3012" }

func (r *E3012) ShortDesc() string {
	return "Property value type mismatch"
}

func (r *E3012) Description() string {
	return "Checks that resource property values match the expected types from CloudFormation resource schemas."
}

func (r *E3012) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3012"
}

func (r *E3012) Tags() []string {
	return []string{"resources", "properties", "type"}
}

func (r *E3012) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		rt, err := schema.GetResourceType(res.Type)
		if err != nil || rt == nil {
			// Schema loading error or unknown resource type - skip validation
			continue
		}

		for propName, propValue := range res.Properties {
			prop := rt.GetProperty(propName)
			if prop == nil {
				// Unknown property - handled by E3002
				continue
			}

			// Skip validation if value is an intrinsic function
			if isIntrinsicFunction(propValue) {
				continue
			}

			if err := validatePropertyType(propValue, prop.PrimitiveType, prop.Type); err != nil {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Property '%s' in resource '%s' (%s): %s", propName, resName, res.Type, err.Error()),
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", resName, "Properties", propName},
				})
			}
		}
	}

	return matches
}

// isIntrinsicFunction checks if a value is a CloudFormation intrinsic function.
func isIntrinsicFunction(value any) bool {
	m, ok := value.(map[string]any)
	if !ok {
		return false
	}
	// Check for intrinsic function keys
	for key := range m {
		if strings.HasPrefix(key, "Fn::") || key == "Ref" || key == "Condition" {
			return true
		}
	}
	return false
}

// validatePropertyType validates that a value matches the expected CloudFormation type.
func validatePropertyType(value any, primitiveType, cfnType string) error {
	// Handle primitive types
	if primitiveType != "" {
		return validatePrimitiveType(value, primitiveType)
	}

	// Handle complex types
	switch cfnType {
	case "List":
		if _, ok := value.([]any); !ok {
			return fmt.Errorf("expected list, got %T", value)
		}
	case "Map":
		if _, ok := value.(map[string]any); !ok {
			return fmt.Errorf("expected map, got %T", value)
		}
	default:
		// Named type (e.g., "CorsConfiguration") - should be a map
		if cfnType != "" {
			if _, ok := value.(map[string]any); !ok {
				return fmt.Errorf("expected object (type %s), got %T", cfnType, value)
			}
		}
	}

	return nil
}

// validatePrimitiveType validates a value against a CloudFormation primitive type.
func validatePrimitiveType(value any, primitiveType string) error {
	switch primitiveType {
	case "String":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("expected string, got %T", value)
		}
	case "Integer":
		// YAML parses integers as int or int64, but sometimes as float64
		switch v := value.(type) {
		case int, int64:
			// OK
		case float64:
			// Check if it's actually an integer value
			if v != float64(int64(v)) {
				return fmt.Errorf("expected integer, got float %v", v)
			}
		default:
			return fmt.Errorf("expected integer, got %T", value)
		}
	case "Long":
		// Same handling as Integer
		switch v := value.(type) {
		case int, int64:
			// OK
		case float64:
			if v != float64(int64(v)) {
				return fmt.Errorf("expected long, got float %v", v)
			}
		default:
			return fmt.Errorf("expected long, got %T", value)
		}
	case "Double":
		switch value.(type) {
		case float64, int, int64:
			// OK - integers are valid as doubles
		default:
			return fmt.Errorf("expected double, got %T", value)
		}
	case "Boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("expected boolean, got %T", value)
		}
	case "Timestamp":
		// Timestamps are typically strings in CloudFormation
		if _, ok := value.(string); !ok {
			return fmt.Errorf("expected timestamp (string), got %T", value)
		}
	case "Json":
		// Json type can be string, map, or list
		switch value.(type) {
		case string, map[string]any, []any:
			// OK
		default:
			return fmt.Errorf("expected JSON (string, map, or list), got %T", value)
		}
	}

	return nil
}
