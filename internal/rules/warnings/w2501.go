package warnings

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W2501{})
}

// W2501 warns about password-related parameter configuration issues.
type W2501 struct{}

func (r *W2501) ID() string { return "W2501" }

func (r *W2501) ShortDesc() string {
	return "Password properties configuration"
}

func (r *W2501) Description() string {
	return "Warns about password-related parameters that may not have NoEcho enabled or have insufficient constraints."
}

func (r *W2501) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/parameters-section-structure.html"
}

func (r *W2501) Tags() []string {
	return []string{"warnings", "parameters", "security", "passwords"}
}

// passwordKeywords are words that suggest a parameter contains a password
var passwordKeywords = []string{
	"password", "passwd", "secret", "credential", "apikey", "api_key",
	"accesskey", "access_key", "privatekey", "private_key", "token",
}

func (r *W2501) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for paramName, param := range tmpl.Parameters {
		lowerName := strings.ToLower(paramName)

		// Check if parameter name suggests it's a password/secret
		isPasswordLike := false
		for _, keyword := range passwordKeywords {
			if strings.Contains(lowerName, keyword) {
				isPasswordLike = true
				break
			}
		}

		if !isPasswordLike {
			continue
		}

		// Check if NoEcho is enabled
		if !param.NoEcho {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Parameter '%s' appears to be a password/secret but NoEcho is not enabled", paramName),
				Path:    []string{"Parameters", paramName},
			})
		}

		// Check if there's a minimum length constraint
		hasMinLength := param.MinLength != nil && *param.MinLength > 0

		if !hasMinLength && param.AllowedPattern == "" {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Parameter '%s' appears to be a password but has no MinLength or AllowedPattern constraint", paramName),
				Path:    []string{"Parameters", paramName},
			})
		}

		// Check if there's a default value (passwords shouldn't have defaults)
		if param.Default != nil {
			defaultStr := fmt.Sprintf("%v", param.Default)
			if defaultStr != "" {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Parameter '%s' appears to be a password but has a default value; consider removing the default", paramName),
					Path:    []string{"Parameters", paramName, "Default"},
				})
			}
		}
	}

	return matches
}
