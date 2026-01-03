package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
	"gopkg.in/yaml.v3"
)

func init() {
	rules.Register(&E3028{})
}

// E3028 validates resource metadata section.
type E3028 struct{}

func (r *E3028) ID() string {
	return "E3028"
}

func (r *E3028) ShortDesc() string {
	return "Validate the metadata section of a resource"
}

func (r *E3028) Description() string {
	return "Checks resource metadata structure where feasible"
}

func (r *E3028) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-attribute-metadata.html"
}

func (r *E3028) Tags() []string {
	return []string{"resources", "metadata"}
}

func (r *E3028) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Node == nil || res.Node.Kind != yaml.MappingNode {
			continue
		}

		// Look for Metadata in the resource node
		var metadataNode *yaml.Node
		for i := 0; i < len(res.Node.Content); i += 2 {
			key := res.Node.Content[i]
			value := res.Node.Content[i+1]
			if key.Value == "Metadata" {
				metadataNode = value
				break
			}
		}

		if metadataNode == nil {
			continue
		}

		// Metadata must be an object
		if metadataNode.Kind != yaml.MappingNode {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Resource '%s' has invalid Metadata (must be an object)", resName),
				Line:    metadataNode.Line,
				Column:  metadataNode.Column,
				Path:    []string{"Resources", resName, "Metadata"},
			})
		}
	}

	return matches
}
