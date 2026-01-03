// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3503{})
}

// E3503 validates ACM certificate ValidationDomain.
type E3503 struct{}

func (r *E3503) ID() string { return "E3503" }

func (r *E3503) ShortDesc() string {
	return "Certificate ValidationDomain"
}

func (r *E3503) Description() string {
	return "Validates that ACM certificate ValidationDomain is a superdomain of or equal to the DomainName in validation options."
}

func (r *E3503) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3503"
}

func (r *E3503) Tags() []string {
	return []string{"resources", "properties", "acm", "certificate"}
}

func (r *E3503) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::CertificateManager::Certificate" {
			continue
		}

		domainValidationOptions, hasOptions := res.Properties["DomainValidationOptions"]
		if !hasOptions {
			continue
		}

		optionsList, ok := domainValidationOptions.([]interface{})
		if !ok {
			continue
		}

		for i, option := range optionsList {
			optionMap, ok := option.(map[string]interface{})
			if !ok {
				continue
			}

			domainName, hasDomain := optionMap["DomainName"]
			validationDomain, hasValidation := optionMap["ValidationDomain"]

			if !hasDomain || !hasValidation {
				continue
			}

			domainNameStr, ok1 := domainName.(string)
			validationDomainStr, ok2 := validationDomain.(string)

			if !ok1 || !ok2 {
				continue
			}

			// Normalize by ensuring trailing dots for comparison
			if !strings.HasSuffix(domainNameStr, ".") {
				domainNameStr += "."
			}
			if !strings.HasSuffix(validationDomainStr, ".") {
				validationDomainStr += "."
			}

			// Check if domainName ends with validationDomain or equals it
			if !strings.HasSuffix(domainNameStr, validationDomainStr) && domainNameStr != validationDomainStr {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': DomainValidationOption %d ValidationDomain '%s' must be a superdomain of or equal to DomainName '%s'",
						resName, i, validationDomain, domainName,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", "DomainValidationOptions", fmt.Sprintf("[%d]", i), "ValidationDomain"},
				})
			}
		}
	}

	return matches
}
