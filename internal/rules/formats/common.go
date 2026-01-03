// Package formats contains format validation rules (E11xx).
package formats

import "fmt"

type formatRef struct {
	value    string
	rawValue any
	path     []string
}

// findFormatReferences finds string values that start with a given prefix
func findFormatReferences(v any, prefix string, propertyNames []string) []formatRef {
	var results []formatRef
	findFormatReferencesRecursive(v, []string{}, prefix, propertyNames, &results)
	return results
}

func findFormatReferencesRecursive(v any, path []string, prefix string, propertyNames []string, results *[]formatRef) {
	switch val := v.(type) {
	case string:
		// Check if this looks like the expected format
		if len(val) >= len(prefix) && val[:len(prefix)] == prefix {
			*results = append(*results, formatRef{
				value:    val,
				rawValue: v,
				path:     path,
			})
		}
	case map[string]any:
		// Skip intrinsic functions
		if isIntrinsicFunction(val) {
			return
		}
		for key, child := range val {
			findFormatReferencesRecursive(child, append(path, key), prefix, propertyNames, results)
		}
	case []any:
		for i, child := range val {
			findFormatReferencesRecursive(child, append(path, fmt.Sprintf("[%d]", i)), prefix, propertyNames, results)
		}
	}
}
