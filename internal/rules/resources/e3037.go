// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"encoding/json"
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/schema"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3037{})
}

// E3037 checks that list properties have unique items where required.
type E3037 struct{}

func (r *E3037) ID() string { return "E3037" }

func (r *E3037) ShortDesc() string {
	return "List has duplicate items"
}

func (r *E3037) Description() string {
	return "Checks that list properties requiring unique items do not contain duplicates."
}

func (r *E3037) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3037"
}

func (r *E3037) Tags() []string {
	return []string{"resources", "properties", "list", "unique"}
}

func (r *E3037) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		for propName, propValue := range res.Properties {
			// Check if this property requires unique items
			if !schema.RequiresUniqueItems(res.Type, propName) {
				continue
			}

			// Must be a list
			list, ok := propValue.([]any)
			if !ok {
				continue
			}

			// Check for duplicates
			duplicates := findDuplicates(list)
			for _, dup := range duplicates {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Property '%s' in resource '%s' (%s) contains duplicate item: %v",
						propName, resName, res.Type, dup,
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

// findDuplicates returns a list of duplicate items in the slice.
func findDuplicates(items []any) []any {
	seen := make(map[string]bool)
	reported := make(map[string]bool)
	var duplicates []any

	for _, item := range items {
		// Skip intrinsic functions
		if isIntrinsic(item) {
			continue
		}

		// Serialize to JSON for comparison (handles nested structures)
		key := toJSONKey(item)
		if seen[key] && !reported[key] {
			duplicates = append(duplicates, item)
			reported[key] = true
		}
		seen[key] = true
	}

	return duplicates
}

// toJSONKey converts a value to a JSON string for map key comparison.
func toJSONKey(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(b)
}

// isIntrinsic checks if a value is a CloudFormation intrinsic function.
func isIntrinsic(v any) bool {
	m, ok := v.(map[string]any)
	if !ok {
		return false
	}
	for key := range m {
		if key == "Ref" || len(key) > 4 && key[:4] == "Fn::" {
			return true
		}
	}
	return false
}
