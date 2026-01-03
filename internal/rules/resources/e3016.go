package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
	"gopkg.in/yaml.v3"
)

func init() {
	rules.Register(&E3016{})
}

// E3016 validates resource UpdatePolicy configuration.
type E3016 struct{}

func (r *E3016) ID() string {
	return "E3016"
}

func (r *E3016) ShortDesc() string {
	return "Check the configuration of a resources UpdatePolicy"
}

func (r *E3016) Description() string {
	return "Ensure resource UpdatePolicy attributes are properly structured"
}

func (r *E3016) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-attribute-updatepolicy.html"
}

func (r *E3016) Tags() []string {
	return []string{"resources", "updatepolicy"}
}

func (r *E3016) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Node == nil || res.Node.Kind != yaml.MappingNode {
			continue
		}

		// Look for UpdatePolicy in the resource node
		var updatePolicyNode *yaml.Node
		for i := 0; i < len(res.Node.Content); i += 2 {
			key := res.Node.Content[i]
			value := res.Node.Content[i+1]
			if key.Value == "UpdatePolicy" {
				updatePolicyNode = value
				break
			}
		}

		if updatePolicyNode == nil {
			continue
		}

		// UpdatePolicy must be an object
		if updatePolicyNode.Kind != yaml.MappingNode {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Resource '%s' has invalid UpdatePolicy (must be an object)", resName),
				Line:    updatePolicyNode.Line,
				Column:  updatePolicyNode.Column,
				Path:    []string{"Resources", resName, "UpdatePolicy"},
			})
			continue
		}

		// Valid UpdatePolicy keys based on resource type
		validKeys := map[string]bool{
			"AutoScalingReplacingUpdate":        true,
			"AutoScalingRollingUpdate":          true,
			"AutoScalingScheduledAction":        true,
			"CodeDeployLambdaAliasUpdate":       true,
			"EnableVersionUpgrade":              true,
			"UseOnlineResharding":               true,
			"AutoScalingReplicationGroupUpdate": true,
		}

		// Check for invalid keys
		for i := 0; i < len(updatePolicyNode.Content); i += 2 {
			key := updatePolicyNode.Content[i]
			if !validKeys[key.Value] {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Resource '%s' has invalid UpdatePolicy key '%s'", resName, key.Value),
					Line:    key.Line,
					Column:  key.Column,
					Path:    []string{"Resources", resName, "UpdatePolicy", key.Value},
				})
			}
		}
	}

	return matches
}
