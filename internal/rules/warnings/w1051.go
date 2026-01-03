package warnings

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W1051{})
}

// W1051 warns when Secrets Manager ARN is used in dynamic reference instead of secret name.
type W1051 struct{}

func (r *W1051) ID() string { return "W1051" }

func (r *W1051) ShortDesc() string {
	return "Secrets Manager ARN in dynamic ref"
}

func (r *W1051) Description() string {
	return "Warns when a full Secrets Manager ARN is used in a dynamic reference instead of just the secret name or partial ARN."
}

func (r *W1051) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/dynamic-references.html"
}

func (r *W1051) Tags() []string {
	return []string{"warnings", "functions", "dynamic-references", "secrets-manager"}
}

// secretsManagerDynRefPattern matches {{resolve:secretsmanager:...}}
var secretsManagerDynRefPattern = regexp.MustCompile(`\{\{resolve:secretsmanager:([^}]+)\}\}`)

// secretsManagerARNPattern matches full ARN format
var secretsManagerARNPattern = regexp.MustCompile(`^arn:aws:secretsmanager:[a-z0-9-]+:\d{12}:secret:`)

func (r *W1051) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check resources
	for resName, res := range tmpl.Resources {
		r.checkValue(res.Properties, []string{"Resources", resName, "Properties"}, &matches)
	}

	return matches
}

func (r *W1051) checkValue(v any, path []string, matches *[]rules.Match) {
	switch val := v.(type) {
	case string:
		r.checkString(val, path, matches)
	case map[string]any:
		for key, child := range val {
			r.checkValue(child, append(path, key), matches)
		}
	case []any:
		for i, child := range val {
			r.checkValue(child, append(path, fmt.Sprintf("[%d]", i)), matches)
		}
	}
}

func (r *W1051) checkString(s string, path []string, matches *[]rules.Match) {
	// Find all dynamic references to secrets manager
	dynRefs := secretsManagerDynRefPattern.FindAllStringSubmatch(s, -1)
	for _, ref := range dynRefs {
		if len(ref) >= 2 {
			secretRef := ref[1]
			// Check if the reference uses a full ARN
			if secretsManagerARNPattern.MatchString(secretRef) {
				*matches = append(*matches, rules.Match{
					Message: "Dynamic reference uses full Secrets Manager ARN; consider using just the secret name for simplicity",
					Path:    path,
				})
			}

			// Check for potential issues with the reference format
			parts := strings.Split(secretRef, ":")
			if len(parts) > 0 {
				// Check if secret name contains spaces
				if strings.Contains(parts[0], " ") {
					*matches = append(*matches, rules.Match{
						Message: "Secrets Manager secret name contains spaces which may cause issues",
						Path:    path,
					})
				}
			}
		}
	}
}
