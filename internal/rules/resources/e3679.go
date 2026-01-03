// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3679{})
}

// E3679 validates ELB protocol certificates.
type E3679 struct{}

func (r *E3679) ID() string { return "E3679" }

func (r *E3679) ShortDesc() string {
	return "ELB HTTPS/SSL listeners require certificates"
}

func (r *E3679) Description() string {
	return "Validates that AWS::ElasticLoadBalancing::LoadBalancer listeners with HTTPS or SSL protocol specify SSLCertificateId."
}

func (r *E3679) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3679"
}

func (r *E3679) Tags() []string {
	return []string{"resources", "properties", "elasticloadbalancing", "listener"}
}

func (r *E3679) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::ElasticLoadBalancing::LoadBalancer" {
			continue
		}

		listeners, hasListeners := res.Properties["Listeners"]
		if !hasListeners || isIntrinsicFunction(listeners) {
			continue
		}

		listenerList, ok := listeners.([]any)
		if !ok {
			continue
		}

		for _, listener := range listenerList {
			listenerMap, ok := listener.(map[string]any)
			if !ok {
				continue
			}

			protocol, hasProtocol := listenerMap["Protocol"]
			if !hasProtocol || isIntrinsicFunction(protocol) {
				continue
			}

			protocolStr, ok := protocol.(string)
			if !ok {
				continue
			}

			if protocolStr == "HTTPS" || protocolStr == "SSL" {
				_, hasSSLCertificateId := listenerMap["SSLCertificateId"]
				if !hasSSLCertificateId {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': Listener with Protocol '%s' must specify SSLCertificateId",
							resName, protocolStr,
						),
						Line:   res.Node.Line,
						Column: res.Node.Column,
						Path:   []string{"Resources", resName, "Properties", "Listeners"},
					})
				}
			}
		}
	}

	return matches
}
