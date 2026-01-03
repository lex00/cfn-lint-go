package parameters

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E2004{})
}

// E2004 checks that NoEcho is used appropriately for sensitive parameters.
type E2004 struct{}

func (r *E2004) ID() string { return "E2004" }

func (r *E2004) ShortDesc() string {
	return "Parameter NoEcho configuration"
}

func (r *E2004) Description() string {
	return "Checks that sensitive parameters use NoEcho to prevent value exposure in console/API."
}

func (r *E2004) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E2004"
}

func (r *E2004) Tags() []string {
	return []string{"parameters", "security"}
}

// Sensitive keywords that should use NoEcho
var sensitiveKeywords = []string{
	"password",
	"secret",
	"apikey",
	"api_key",
	"token",
	"credential",
	"passphrase",
	"privatekey",
	"private_key",
}

func (r *E2004) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for name, param := range tmpl.Parameters {
		// Check if parameter name suggests it's sensitive
		nameLower := strings.ToLower(name)
		isSensitive := false

		for _, keyword := range sensitiveKeywords {
			if strings.Contains(nameLower, keyword) {
				isSensitive = true
				break
			}
		}

		// If parameter looks sensitive but doesn't have NoEcho=true
		if isSensitive && !param.NoEcho {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Parameter '%s' appears to contain sensitive data but NoEcho is not set to true", name),
				Line:    param.Node.Line,
				Column:  param.Node.Column,
				Path:    []string{"Parameters", name},
			})
		}
	}

	return matches
}
