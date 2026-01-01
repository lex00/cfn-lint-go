// Package warnings contains warning-level rules (Wxxx).
package warnings

import (
	"fmt"
	"regexp"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W3010{})
}

// W3010 warns about hardcoded availability zones.
type W3010 struct{}

func (r *W3010) ID() string { return "W3010" }

func (r *W3010) ShortDesc() string {
	return "Hardcoded availability zone"
}

func (r *W3010) Description() string {
	return "Warns when an availability zone is hardcoded instead of using Fn::GetAZs or a parameter."
}

func (r *W3010) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getavailabilityzones.html"
}

func (r *W3010) Tags() []string {
	return []string{"warnings", "resources", "availability-zones"}
}

// Pattern to match AWS availability zone format (e.g., us-east-1a, eu-west-2b)
var azPattern = regexp.MustCompile(`^[a-z]{2}-[a-z]+-\d[a-z]$`)

// Properties that typically contain availability zones
var azPropertyNames = map[string]bool{
	"availabilityzone":  true,
	"availabilityzones": true,
}

func (r *W3010) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		hardcodedAZs := findHardcodedAZs(res.Properties, nil)
		for _, az := range hardcodedAZs {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Hardcoded availability zone '%s' in resource '%s'. Use Fn::GetAZs or a parameter for portability.", az.value, resName),
				Path:    append([]string{"Resources", resName, "Properties"}, az.path...),
			})
		}
	}

	return matches
}

type azInfo struct {
	value string
	path  []string
}

func findHardcodedAZs(v any, path []string) []azInfo {
	var results []azInfo

	switch val := v.(type) {
	case string:
		if azPattern.MatchString(val) {
			results = append(results, azInfo{value: val, path: path})
		}
	case map[string]any:
		// Skip intrinsic functions
		for key := range val {
			if key == "Ref" || key == "Fn::GetAZs" || key == "Fn::Select" || key == "Fn::If" {
				return results
			}
		}
		for key, child := range val {
			newPath := append(append([]string{}, path...), key)
			results = append(results, findHardcodedAZs(child, newPath)...)
		}
	case []any:
		for i, child := range val {
			newPath := append(append([]string{}, path...), fmt.Sprintf("[%d]", i))
			results = append(results, findHardcodedAZs(child, newPath)...)
		}
	}

	return results
}
