package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3029{})
}

// E3029 validates Route53 record set aliases.
type E3029 struct{}

func (r *E3029) ID() string {
	return "E3029"
}

func (r *E3029) ShortDesc() string {
	return "Validate Route53 record set aliases"
}

func (r *E3029) Description() string {
	return "Ensures alias records don't specify incompatible TTL or type combinations"
}

func (r *E3029) Source() string {
	return "https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/resource-record-sets-choosing-alias-non-alias.html"
}

func (r *E3029) Tags() []string {
	return []string{"resources", "route53", "alias"}
}

func (r *E3029) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Route53::RecordSet" {
			continue
		}

		// Check if this is an alias record
		aliasTargetRaw, hasAlias := res.Properties["AliasTarget"]
		if !hasAlias {
			continue
		}

		// Skip intrinsic functions
		if isIntrinsicFunction(aliasTargetRaw) {
			continue
		}

		aliasTarget, ok := aliasTargetRaw.(map[string]any)
		if !ok {
			continue
		}

		// Alias records must not have TTL
		if _, hasTTL := res.Properties["TTL"]; hasTTL {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Route53 RecordSet '%s' with AliasTarget cannot specify TTL", resName),
				Line:    res.Node.Line,
				Column:  res.Node.Column,
				Path:    []string{"Resources", resName, "Properties", "TTL"},
			})
		}

		// Alias records must have DNSName and HostedZoneId
		if _, hasDNSName := aliasTarget["DNSName"]; !hasDNSName {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Route53 RecordSet '%s' AliasTarget must have DNSName", resName),
				Line:    res.Node.Line,
				Column:  res.Node.Column,
				Path:    []string{"Resources", resName, "Properties", "AliasTarget"},
			})
		}

		if _, hasHostedZoneId := aliasTarget["HostedZoneId"]; !hasHostedZoneId {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Route53 RecordSet '%s' AliasTarget must have HostedZoneId", resName),
				Line:    res.Node.Line,
				Column:  res.Node.Column,
				Path:    []string{"Resources", resName, "Properties", "AliasTarget"},
			})
		}
	}

	return matches
}
