// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3676{})
}

// E3676 validates ELBv2 protocol certificates.
type E3676 struct{}

func (r *E3676) ID() string { return "E3676" }

func (r *E3676) ShortDesc() string {
	return "ELBv2 HTTPS/TLS listeners require certificates"
}

func (r *E3676) Description() string {
	return "Validates that AWS::ElasticLoadBalancingV2::Listener resources with HTTPS or TLS protocol specify certificates."
}

func (r *E3676) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3676"
}

func (r *E3676) Tags() []string {
	return []string{"resources", "properties", "elasticloadbalancingv2", "listener"}
}

func (r *E3676) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::ElasticLoadBalancingV2::Listener" {
			continue
		}

		protocol, hasProtocol := res.Properties["Protocol"]
		if !hasProtocol || isIntrinsicFunction(protocol) {
			continue
		}

		protocolStr, ok := protocol.(string)
		if !ok {
			continue
		}

		if protocolStr == "HTTPS" || protocolStr == "TLS" {
			_, hasCertificates := res.Properties["Certificates"]
			if !hasCertificates {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': Listener with Protocol '%s' must specify Certificates",
						resName, protocolStr,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	return matches
}
