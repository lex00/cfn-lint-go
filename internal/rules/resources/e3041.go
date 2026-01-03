// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3041{})
}

// E3041 checks that RecordSet HostedZoneName is a superdomain of or equal to Name.
type E3041 struct{}

func (r *E3041) ID() string { return "E3041" }

func (r *E3041) ShortDesc() string {
	return "RecordSet HostedZoneName is superdomain"
}

func (r *E3041) Description() string {
	return "Validates that in a Route53 RecordSet, the HostedZoneName must be a superdomain of or equal to the Name being validated."
}

func (r *E3041) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3041"
}

func (r *E3041) Tags() []string {
	return []string{"resources", "properties", "route53"}
}

func (r *E3041) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Route53::RecordSet" {
			continue
		}

		hostedZoneName, hasHostedZone := res.Properties["HostedZoneName"]
		name, hasName := res.Properties["Name"]

		if !hasHostedZone || !hasName {
			continue
		}

		// Extract string values
		hostedZoneStr, ok1 := hostedZoneName.(string)
		nameStr, ok2 := name.(string)

		if !ok1 || !ok2 {
			// Skip if not simple strings (could be Ref, GetAtt, etc.)
			continue
		}

		// Normalize by ensuring trailing dots
		if !strings.HasSuffix(hostedZoneStr, ".") {
			hostedZoneStr += "."
		}
		if !strings.HasSuffix(nameStr, ".") {
			nameStr += "."
		}

		// Check if nameStr ends with hostedZoneStr or equals it
		if !strings.HasSuffix(nameStr, hostedZoneStr) && nameStr != hostedZoneStr {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': RecordSet Name '%s' must be a subdomain of or equal to HostedZoneName '%s'",
					resName, name, hostedZoneName,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties", "Name"},
			})
		}
	}

	return matches
}
