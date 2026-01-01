// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
	"gopkg.in/yaml.v3"
)

func init() {
	rules.Register(&E3002{})
}

// E3002 checks that resource Properties is a valid structure.
type E3002 struct{}

func (r *E3002) ID() string { return "E3002" }

func (r *E3002) ShortDesc() string {
	return "Resource Properties structure error"
}

func (r *E3002) Description() string {
	return "Checks that resource Properties is a mapping (object), not a scalar or list."
}

func (r *E3002) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3002"
}

func (r *E3002) Tags() []string {
	return []string{"resources", "properties"}
}

func (r *E3002) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for name, res := range tmpl.Resources {
		if res.Node == nil || res.Node.Kind != yaml.MappingNode {
			continue
		}

		// Find Properties node
		for i := 0; i < len(res.Node.Content); i += 2 {
			key := res.Node.Content[i]
			value := res.Node.Content[i+1]

			if key.Value == "Properties" {
				// Properties must be a mapping
				if value.Kind != yaml.MappingNode {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Resource '%s' Properties must be an object, got %s", name, nodeKindName(value.Kind)),
						Line:    value.Line,
						Column:  value.Column,
						Path:    []string{"Resources", name, "Properties"},
					})
				}
				break
			}
		}
	}

	return matches
}

func nodeKindName(kind yaml.Kind) string {
	switch kind {
	case yaml.DocumentNode:
		return "document"
	case yaml.SequenceNode:
		return "list"
	case yaml.MappingNode:
		return "object"
	case yaml.ScalarNode:
		return "scalar"
	case yaml.AliasNode:
		return "alias"
	default:
		return "unknown"
	}
}
