package warnings

import (
	"fmt"
	"regexp"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W2031{})
}

// W2031 warns about potential issues with AllowedPattern constraints.
type W2031 struct{}

func (r *W2031) ID() string { return "W2031" }

func (r *W2031) ShortDesc() string {
	return "Parameter AllowedPattern check"
}

func (r *W2031) Description() string {
	return "Warns about potential issues with AllowedPattern regex constraints, such as patterns that are too permissive or may not match intended values."
}

func (r *W2031) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/parameters-section-structure.html"
}

func (r *W2031) Tags() []string {
	return []string{"warnings", "parameters", "validation", "regex"}
}

func (r *W2031) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for paramName, param := range tmpl.Parameters {
		if param.AllowedPattern == "" {
			continue
		}

		// Check if pattern compiles
		_, err := regexp.Compile(param.AllowedPattern)
		if err != nil {
			// This is actually an error, but we report it as a warning with helpful info
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Parameter '%s' has invalid AllowedPattern: %v", paramName, err),
				Path:    []string{"Parameters", paramName, "AllowedPattern"},
			})
			continue
		}

		// Check for overly permissive patterns
		if param.AllowedPattern == ".*" || param.AllowedPattern == ".+" || param.AllowedPattern == "^.*$" {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Parameter '%s' has AllowedPattern '%s' which matches almost anything; consider a more specific pattern", paramName, param.AllowedPattern),
				Path:    []string{"Parameters", paramName, "AllowedPattern"},
			})
		}

		// Check for patterns without anchors
		if len(param.AllowedPattern) > 2 && !startsWithAnchor(param.AllowedPattern) && !endsWithAnchor(param.AllowedPattern) {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Parameter '%s' AllowedPattern '%s' has no anchors (^ or $); pattern may match unintended substrings", paramName, param.AllowedPattern),
				Path:    []string{"Parameters", paramName, "AllowedPattern"},
			})
		}

		// Check if default matches the pattern
		if param.Default != nil {
			defaultStr := fmt.Sprintf("%v", param.Default)
			re, _ := regexp.Compile(param.AllowedPattern)
			if re != nil && !re.MatchString(defaultStr) {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Parameter '%s' default value '%s' does not match AllowedPattern '%s'", paramName, defaultStr, param.AllowedPattern),
					Path:    []string{"Parameters", paramName, "Default"},
				})
			}
		}
	}

	return matches
}

func startsWithAnchor(pattern string) bool {
	return len(pattern) > 0 && (pattern[0] == '^' || pattern[0] == '\\')
}

func endsWithAnchor(pattern string) bool {
	return len(pattern) > 0 && pattern[len(pattern)-1] == '$'
}
