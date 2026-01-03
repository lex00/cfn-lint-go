package resources

import (
	"fmt"
	"regexp"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3024{})
}

// E3024 validates tag configuration in resources.
type E3024 struct{}

func (r *E3024) ID() string {
	return "E3024"
}

func (r *E3024) ShortDesc() string {
	return "Validate tag configuration"
}

func (r *E3024) Description() string {
	return "Verifies tag uniqueness and pattern compliance in resource metadata"
}

func (r *E3024) Source() string {
	return "https://docs.aws.amazon.com/general/latest/gr/aws_tagging.html"
}

func (r *E3024) Tags() []string {
	return []string{"resources", "tags"}
}

// Tag key pattern: up to 128 characters, letters, numbers, spaces, and +-=._:/@
var tagKeyPattern = regexp.MustCompile(`^[\w\s\+\-=._:/@]{1,128}$`)

func (r *E3024) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		tagsRaw, hasTags := res.Properties["Tags"]
		if !hasTags {
			continue
		}

		// Skip intrinsic functions
		if isIntrinsicFunction(tagsRaw) {
			continue
		}

		tags, ok := tagsRaw.([]any)
		if !ok {
			continue
		}

		// Track tag keys for uniqueness
		tagKeys := make(map[string]bool)

		for i, tagRaw := range tags {
			// Skip intrinsic functions
			if isIntrinsicFunction(tagRaw) {
				continue
			}

			tag, ok := tagRaw.(map[string]any)
			if !ok {
				continue
			}

			// Get tag key
			keyRaw, hasKey := tag["Key"]
			if !hasKey {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Resource '%s' has tag without Key at index %d", resName, i),
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", resName, "Properties", "Tags", fmt.Sprintf("[%d]", i)},
				})
				continue
			}

			key, ok := keyRaw.(string)
			if !ok {
				continue
			}

			// Check for duplicate keys
			if tagKeys[key] {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Resource '%s' has duplicate tag key '%s'", resName, key),
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", resName, "Properties", "Tags", fmt.Sprintf("[%d]", i), "Key"},
				})
			}
			tagKeys[key] = true

			// Validate key pattern
			if !tagKeyPattern.MatchString(key) {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Resource '%s' has invalid tag key '%s' (must be 1-128 characters, letters, numbers, spaces, and +-=._:/@)", resName, key),
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", resName, "Properties", "Tags", fmt.Sprintf("[%d]", i), "Key"},
				})
			}

			// Check value exists
			if _, hasValue := tag["Value"]; !hasValue {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Resource '%s' has tag without Value at index %d", resName, i),
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", resName, "Properties", "Tags", fmt.Sprintf("[%d]", i)},
				})
			}
		}
	}

	return matches
}
