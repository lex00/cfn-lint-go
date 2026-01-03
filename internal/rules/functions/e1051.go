// Package functions contains intrinsic function validation rules (E1xxx).
package functions

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1051{})
}

// E1051 validates Secrets Manager dynamic references.
type E1051 struct{}

func (r *E1051) ID() string { return "E1051" }

func (r *E1051) ShortDesc() string {
	return "Secrets Manager dynamic reference validation"
}

func (r *E1051) Description() string {
	return "Validates that dynamic references to Secrets Manager are only used in resource properties and have proper format."
}

func (r *E1051) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/dynamic-references-secretsmanager.html"
}

func (r *E1051) Tags() []string {
	return []string{"functions", "dynamic-references", "secretsmanager"}
}

func (r *E1051) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check resources - this is where secretsmanager dynamic refs are allowed
	for resName, res := range tmpl.Resources {
		dynRefs := findAllDynamicRefs(res.Properties)
		for _, dr := range dynRefs {
			if strings.HasPrefix(dr.content, "secretsmanager:") {
				if err := r.validateSecretsManagerRef(dr); err != "" {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("%s in resource '%s'", err, resName),
						Line:    dr.line,
						Column:  dr.column,
						Path:    []string{"Resources", resName, "Properties"},
					})
				}
			}
		}
	}

	// Check outputs - secretsmanager dynamic refs should not be used in outputs
	// as they may expose secrets
	for outName, out := range tmpl.Outputs {
		dynRefs := findAllDynamicRefs(out.Value)
		for _, dr := range dynRefs {
			if strings.HasPrefix(dr.content, "secretsmanager:") {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Dynamic references from Secrets Manager can only be used in resource properties, not in output '%s'", outName),
					Line:    dr.line,
					Column:  dr.column,
					Path:    []string{"Outputs", outName, "Value"},
				})
			}
		}
	}

	return matches
}

func (r *E1051) validateSecretsManagerRef(dr dynamicRefInfo) string {
	// Format: secretsmanager:secret-id:SecretString:json-key:version-stage:version-id
	// Minimum: secretsmanager:secret-id
	parts := strings.Split(dr.content, ":")

	if len(parts) < 2 {
		return fmt.Sprintf("Secrets Manager dynamic reference '%s' is malformed, expected format: {{resolve:secretsmanager:secret-id[:SecretString:json-key:version-stage:version-id]}}", dr.fullRef)
	}

	secretID := parts[1]
	if secretID == "" {
		return fmt.Sprintf("Secrets Manager dynamic reference '%s' is missing secret ID", dr.fullRef)
	}

	// Check for trailing backslash (not supported)
	if strings.HasSuffix(dr.content, "\\") {
		return fmt.Sprintf("Secrets Manager dynamic reference '%s' cannot end with a backslash (\\)", dr.fullRef)
	}

	// If more parts are specified, validate them
	if len(parts) > 2 {
		// parts[2] should be "SecretString" if specified
		if parts[2] != "" && parts[2] != "SecretString" {
			return fmt.Sprintf("Secrets Manager dynamic reference '%s' has invalid secret string specifier, expected 'SecretString'", dr.fullRef)
		}

		// Check that json-key and version-stage don't contain colons
		if len(parts) > 3 && parts[3] != "" && strings.Contains(parts[3], ":") {
			return fmt.Sprintf("Secrets Manager dynamic reference '%s' json-key cannot contain colons", dr.fullRef)
		}
		if len(parts) > 4 && parts[4] != "" && strings.Contains(parts[4], ":") {
			return fmt.Sprintf("Secrets Manager dynamic reference '%s' version-stage cannot contain colons", dr.fullRef)
		}
	}

	return ""
}
