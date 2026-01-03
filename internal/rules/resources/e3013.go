package resources

import (
	"fmt"
	"regexp"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3013{})
}

// E3013 validates CloudFront aliases contain valid domain names.
type E3013 struct{}

func (r *E3013) ID() string {
	return "E3013"
}

func (r *E3013) ShortDesc() string {
	return "CloudFront Aliases"
}

func (r *E3013) Description() string {
	return "CloudFront aliases should contain valid domain names"
}

func (r *E3013) Source() string {
	return "https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/CNAMEs.html"
}

func (r *E3013) Tags() []string {
	return []string{"resources", "cloudfront", "aliases"}
}

// Domain name regex - supports wildcards and standard domain format
var domainNameRegex = regexp.MustCompile(`^(\*\.)?([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)

func (r *E3013) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::CloudFront::Distribution" {
			continue
		}

		// Navigate to DistributionConfig.Aliases
		distConfig, ok := res.Properties["DistributionConfig"].(map[string]any)
		if !ok {
			continue
		}

		aliasesRaw, hasAliases := distConfig["Aliases"]
		if !hasAliases {
			continue
		}

		// Skip intrinsic functions
		if isIntrinsicFunction(aliasesRaw) {
			continue
		}

		aliases, ok := aliasesRaw.([]any)
		if !ok {
			continue
		}

		// Validate each alias
		for i, aliasRaw := range aliases {
			// Skip intrinsic functions in individual aliases
			if isIntrinsicFunction(aliasRaw) {
				continue
			}

			alias, ok := aliasRaw.(string)
			if !ok {
				continue
			}

			if !domainNameRegex.MatchString(alias) {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Resource '%s' has invalid CloudFront alias '%s' (must be a valid domain name)", resName, alias),
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", resName, "Properties", "DistributionConfig", "Aliases", fmt.Sprintf("[%d]", i)},
				})
			}
		}
	}

	return matches
}
