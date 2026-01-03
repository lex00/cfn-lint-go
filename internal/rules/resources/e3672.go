// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3672{})
}

// E3672 validates DAX Cluster node type.
type E3672 struct{}

func (r *E3672) ID() string { return "E3672" }

func (r *E3672) ShortDesc() string {
	return "Validate DAX Cluster node type"
}

func (r *E3672) Description() string {
	return "Validates that AWS::DAX::Cluster resources specify valid node types."
}

func (r *E3672) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3672"
}

func (r *E3672) Tags() []string {
	return []string{"resources", "properties", "dax", "nodetype"}
}

func (r *E3672) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::DAX::Cluster" {
			continue
		}

		nodeType, hasNodeType := res.Properties["NodeType"]
		if !hasNodeType || isIntrinsicFunction(nodeType) {
			continue
		}

		nodeTypeStr, ok := nodeType.(string)
		if !ok {
			continue
		}

		// DAX node types start with dax.
		if !strings.HasPrefix(nodeTypeStr, "dax.") {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Invalid DAX node type '%s'. Must start with 'dax.'",
					resName, nodeTypeStr,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties", "NodeType"},
			})
		}
	}

	return matches
}
