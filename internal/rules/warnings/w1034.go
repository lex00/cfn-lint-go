package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W1034{})
}

// W1034 warns about FindInMap function value issues.
type W1034 struct{}

func (r *W1034) ID() string { return "W1034" }

func (r *W1034) ShortDesc() string {
	return "FindInMap function value validation"
}

func (r *W1034) Description() string {
	return "Warns about potential issues with Fn::FindInMap values, such as hardcoded keys that could use parameters."
}

func (r *W1034) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-findinmap.html"
}

func (r *W1034) Tags() []string {
	return []string{"warnings", "functions", "findinmap", "mappings"}
}

func (r *W1034) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check resources
	for resName, res := range tmpl.Resources {
		r.checkValue(res.Properties, []string{"Resources", resName, "Properties"}, tmpl, &matches)
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		r.checkValue(out.Value, []string{"Outputs", outName, "Value"}, tmpl, &matches)
	}

	return matches
}

func (r *W1034) checkValue(v any, path []string, tmpl *template.Template, matches *[]rules.Match) {
	switch val := v.(type) {
	case map[string]any:
		if findInMap, ok := val["Fn::FindInMap"]; ok {
			r.checkFindInMap(findInMap, path, tmpl, matches)
		}
		for key, child := range val {
			r.checkValue(child, append(path, key), tmpl, matches)
		}
	case []any:
		for i, child := range val {
			r.checkValue(child, append(path, fmt.Sprintf("[%d]", i)), tmpl, matches)
		}
	}
}

func (r *W1034) checkFindInMap(findInMap any, path []string, tmpl *template.Template, matches *[]rules.Match) {
	arr, ok := findInMap.([]any)
	if !ok || len(arr) < 3 {
		return
	}

	// Get map name
	mapName, ok := arr[0].(string)
	if !ok {
		return
	}

	// Check if second level key is hardcoded when it could use AWS::Region
	if secondKey, ok := arr[1].(string); ok {
		// Check if the key looks like a region and could use !Ref AWS::Region
		regions := map[string]bool{
			"us-east-1": true, "us-east-2": true, "us-west-1": true, "us-west-2": true,
			"eu-west-1": true, "eu-west-2": true, "eu-west-3": true, "eu-central-1": true,
			"ap-northeast-1": true, "ap-northeast-2": true, "ap-southeast-1": true, "ap-southeast-2": true,
			"ap-south-1": true, "sa-east-1": true, "ca-central-1": true,
		}
		if regions[secondKey] {
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf("Fn::FindInMap uses hardcoded region '%s'; consider using { Ref: AWS::Region } for portability", secondKey),
				Path:    path,
			})
		}
	}

	// Check if mapping exists and has the specified keys
	if mapping, ok := tmpl.Mappings[mapName]; ok {
		if secondKey, ok := arr[1].(string); ok {
			if secondLevel, ok := mapping.Values[secondKey]; ok {
				if thirdKey, ok := arr[2].(string); ok {
					if _, exists := secondLevel[thirdKey]; !exists {
						// Key doesn't exist - but this is an error, not a warning
						// Just check for potential typos by looking at similar keys
						for key := range secondLevel {
							if r.isSimilar(key, thirdKey) {
								*matches = append(*matches, rules.Match{
									Message: fmt.Sprintf("Fn::FindInMap third key '%s' not found in mapping '%s.%s'; did you mean '%s'?", thirdKey, mapName, secondKey, key),
									Path:    path,
								})
								break
							}
						}
					}
				}
			}
		}
	}
}

func (r *W1034) isSimilar(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	diffs := 0
	for i := range a {
		if a[i] != b[i] {
			diffs++
		}
	}
	return diffs == 1 // One character difference
}
