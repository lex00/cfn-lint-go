// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3667{})
}

// E3667 validates RedShift cluster node type.
type E3667 struct{}

func (r *E3667) ID() string { return "E3667" }

func (r *E3667) ShortDesc() string {
	return "Validate RedShift cluster node type"
}

func (r *E3667) Description() string {
	return "Validates that AWS::Redshift::Cluster resources specify valid node types."
}

func (r *E3667) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3667"
}

func (r *E3667) Tags() []string {
	return []string{"resources", "properties", "redshift", "nodetype"}
}

func (r *E3667) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Redshift node type families
	validFamilies := []string{
		"dc1.", "dc2.", "ds2.", "ra3.",
	}

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Redshift::Cluster" {
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

		// Check if node type starts with a valid family
		isValid := false
		for _, family := range validFamilies {
			if strings.HasPrefix(nodeTypeStr, family) {
				isValid = true
				break
			}
		}

		if !isValid {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Invalid Redshift node type '%s'. Must start with dc1., dc2., ds2., or ra3.",
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
