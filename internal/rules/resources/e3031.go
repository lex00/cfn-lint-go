// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/schema"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3031{})
}

// E3031 checks that string property values match required patterns.
type E3031 struct{}

func (r *E3031) ID() string { return "E3031" }

func (r *E3031) ShortDesc() string {
	return "Property value pattern mismatch"
}

func (r *E3031) Description() string {
	return "Checks that string property values match required regex patterns from CloudFormation schemas."
}

func (r *E3031) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3031"
}

func (r *E3031) Tags() []string {
	return []string{"resources", "properties", "pattern"}
}

// patternCache caches compiled regex patterns.
var patternCache = make(map[string]*regexp.Regexp)

func (r *E3031) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if !schema.HasConstraints(res.Type) {
			continue
		}

		for propName, propValue := range res.Properties {
			// Skip non-string values and intrinsic functions
			strValue, ok := propValue.(string)
			if !ok {
				continue
			}

			constraints := schema.GetPropertyConstraints(res.Type, propName)
			if constraints == nil || constraints.Pattern == "" {
				continue
			}

			// Get or compile the pattern
			pattern, ok := patternCache[constraints.Pattern]
			if !ok {
				var err error
				pattern, err = regexp.Compile(constraints.Pattern)
				if err != nil {
					// Invalid pattern - skip validation
					continue
				}
				patternCache[constraints.Pattern] = pattern
			}

			if !pattern.MatchString(strValue) {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Property '%s' in resource '%s' (%s): value '%s' does not match pattern '%s'",
						propName, resName, res.Type, truncateString(strValue, 50), constraints.Pattern,
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

// truncateString truncates a string to maxLen characters, adding "..." if truncated.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	// Remove newlines for cleaner output
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
