package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
	"gopkg.in/yaml.v3"
)

func init() {
	rules.Register(&E3038{})
}

// E3038 validates that templates with Serverless resources have the Serverless transform.
type E3038 struct{}

func (r *E3038) ID() string {
	return "E3038"
}

func (r *E3038) ShortDesc() string {
	return "Check if Serverless Resources have Serverless Transform"
}

func (r *E3038) Description() string {
	return "Confirms templates with Serverless resources include the Transform"
}

func (r *E3038) Source() string {
	return "https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-specification.html"
}

func (r *E3038) Tags() []string {
	return []string{"resources", "serverless", "transform"}
}

// Serverless resource types
var serverlessResourceTypes = map[string]bool{
	"AWS::Serverless::Function":     true,
	"AWS::Serverless::Api":          true,
	"AWS::Serverless::HttpApi":      true,
	"AWS::Serverless::SimpleTable":  true,
	"AWS::Serverless::Application":  true,
	"AWS::Serverless::LayerVersion": true,
	"AWS::Serverless::StateMachine": true,
}

func (r *E3038) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check if any serverless resources exist
	hasServerlessResources := false
	var serverlessResourceNames []string

	for resName, res := range tmpl.Resources {
		if serverlessResourceTypes[res.Type] {
			hasServerlessResources = true
			serverlessResourceNames = append(serverlessResourceNames, resName)
		}
	}

	if !hasServerlessResources {
		return matches
	}

	// Check if Transform contains AWS::Serverless
	hasServerlessTransform := false

	// Transform can be a string or array
	// Root is a DocumentNode, so we need to get the first content which is the mapping
	var mappingNode *yaml.Node
	if tmpl.Root != nil && len(tmpl.Root.Content) > 0 {
		mappingNode = tmpl.Root.Content[0]
	}

	if mappingNode != nil && mappingNode.Kind == yaml.MappingNode {
		for i := 0; i < len(mappingNode.Content); i += 2 {
			key := mappingNode.Content[i]
			value := mappingNode.Content[i+1]

			if key.Value == "Transform" {
				if value.Kind == yaml.ScalarNode {
					if strings.Contains(value.Value, "AWS::Serverless") {
						hasServerlessTransform = true
					}
				} else if value.Kind == yaml.SequenceNode {
					for _, item := range value.Content {
						if item.Kind == yaml.ScalarNode && strings.Contains(item.Value, "AWS::Serverless") {
							hasServerlessTransform = true
							break
						}
					}
				}
				break
			}
		}
	}

	if !hasServerlessTransform {
		for _, resName := range serverlessResourceNames {
			res := tmpl.Resources[resName]
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Serverless resource '%s' requires Transform: AWS::Serverless-2016-10-31 to be declared", resName),
				Line:    res.Node.Line,
				Column:  res.Node.Column,
				Path:    []string{"Resources", resName},
			})
		}
	}

	return matches
}
