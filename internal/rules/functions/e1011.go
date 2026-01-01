// Package functions contains intrinsic function validation rules (E1xxx).
package functions

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1011{})
}

// E1011 checks that FindInMap references valid mappings and keys.
type E1011 struct{}

func (r *E1011) ID() string { return "E1011" }

func (r *E1011) ShortDesc() string {
	return "FindInMap references undefined mapping"
}

func (r *E1011) Description() string {
	return "Checks that Fn::FindInMap references valid mappings with correct structure."
}

func (r *E1011) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-findinmap.html"
}

func (r *E1011) Tags() []string {
	return []string{"functions", "findinmap"}
}

func (r *E1011) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check all resources for invalid FindInMap references
	for resName, res := range tmpl.Resources {
		findInMaps := findAllFindInMaps(res.Properties)
		for _, fim := range findInMaps {
			if err := r.validateFindInMap(tmpl, fim); err != "" {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("%s in resource '%s'", err, resName),
					Line:    fim.line,
					Column:  fim.column,
					Path:    []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		findInMaps := findAllFindInMaps(out.Value)
		for _, fim := range findInMaps {
			if err := r.validateFindInMap(tmpl, fim); err != "" {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("%s in output '%s'", err, outName),
					Line:    fim.line,
					Column:  fim.column,
					Path:    []string{"Outputs", outName, "Value"},
				})
			}
		}
	}

	return matches
}

func (r *E1011) validateFindInMap(tmpl *template.Template, fim findInMapInfo) string {
	// Check structure - must have 3 elements
	if fim.argCount < 3 {
		return fmt.Sprintf("Fn::FindInMap requires 3 arguments, got %d", fim.argCount)
	}

	// If mapName is not a string (it's a Ref or other intrinsic), skip validation
	if fim.mapName == "" {
		return ""
	}

	// Check mapping exists
	mapping, ok := tmpl.Mappings[fim.mapName]
	if !ok {
		return fmt.Sprintf("Fn::FindInMap references undefined mapping '%s'", fim.mapName)
	}

	// If topLevelKey is not a string (dynamic), skip key validation
	if fim.topLevelKey == "" {
		return ""
	}

	// Check top-level key exists
	if _, ok := mapping.Values[fim.topLevelKey]; !ok {
		return fmt.Sprintf("Fn::FindInMap references undefined top-level key '%s' in mapping '%s'", fim.topLevelKey, fim.mapName)
	}

	// If secondLevelKey is not a string (dynamic), skip key validation
	if fim.secondLevelKey == "" {
		return ""
	}

	// Check second-level key exists
	if _, ok := mapping.Values[fim.topLevelKey][fim.secondLevelKey]; !ok {
		return fmt.Sprintf("Fn::FindInMap references undefined second-level key '%s' in mapping '%s'.'%s'", fim.secondLevelKey, fim.mapName, fim.topLevelKey)
	}

	return ""
}

type findInMapInfo struct {
	mapName        string
	topLevelKey    string
	secondLevelKey string
	argCount       int
	line           int
	column         int
}

func findAllFindInMaps(v any) []findInMapInfo {
	var results []findInMapInfo
	findFindInMapsRecursive(v, &results)
	return results
}

func findFindInMapsRecursive(v any, results *[]findInMapInfo) {
	switch val := v.(type) {
	case map[string]any:
		if fim, ok := val["Fn::FindInMap"]; ok {
			info := parseFindInMapValue(fim)
			*results = append(*results, info)
		}
		for _, child := range val {
			findFindInMapsRecursive(child, results)
		}
	case []any:
		for _, child := range val {
			findFindInMapsRecursive(child, results)
		}
	}
}

func parseFindInMapValue(v any) findInMapInfo {
	arr, ok := v.([]any)
	if !ok {
		return findInMapInfo{}
	}

	info := findInMapInfo{argCount: len(arr)}

	if len(arr) >= 1 {
		if name, ok := arr[0].(string); ok {
			info.mapName = name
		}
	}
	if len(arr) >= 2 {
		if key, ok := arr[1].(string); ok {
			info.topLevelKey = key
		}
	}
	if len(arr) >= 3 {
		if key, ok := arr[2].(string); ok {
			info.secondLevelKey = key
		}
	}

	return info
}
