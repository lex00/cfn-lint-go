// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/schema"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3034{})
}

// E3034 checks that numeric properties are within valid ranges.
type E3034 struct{}

func (r *E3034) ID() string { return "E3034" }

func (r *E3034) ShortDesc() string {
	return "Number out of range"
}

func (r *E3034) Description() string {
	return "Checks that numeric property values are within minimum and maximum constraints from CloudFormation schemas."
}

func (r *E3034) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3034"
}

func (r *E3034) Tags() []string {
	return []string{"resources", "properties", "number", "range"}
}

func (r *E3034) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if !schema.HasConstraints(res.Type) {
			continue
		}

		for propName, propValue := range res.Properties {
			// Extract numeric value
			numValue, ok := toFloat64(propValue)
			if !ok {
				continue
			}

			constraints := schema.GetPropertyConstraints(res.Type, propName)
			if constraints == nil {
				continue
			}

			// Check minimum value
			if constraints.MinValue != nil && numValue < *constraints.MinValue {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Property '%s' in resource '%s' (%s): value %v is less than minimum %v",
						propName, resName, res.Type, numValue, *constraints.MinValue,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", propName},
				})
			}

			// Check maximum value
			if constraints.MaxValue != nil && numValue > *constraints.MaxValue {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Property '%s' in resource '%s' (%s): value %v exceeds maximum %v",
						propName, resName, res.Type, numValue, *constraints.MaxValue,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", propName},
				})
			}
		}
	}

	return matches
}

// toFloat64 converts various numeric types to float64.
// Returns false if the value is not numeric.
func toFloat64(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint32:
		return float64(n), true
	case uint64:
		return float64(n), true
	default:
		return 0, false
	}
}
