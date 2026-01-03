package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/schema"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3019{})
}

// E3019 validates that all resources have unique primary identifiers.
type E3019 struct{}

func (r *E3019) ID() string {
	return "E3019"
}

func (r *E3019) ShortDesc() string {
	return "Validate that all resources have unique primary identifiers"
}

func (r *E3019) Description() string {
	return "Uses schema primary identifiers to confirm resource uniqueness within templates"
}

func (r *E3019) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint"
}

func (r *E3019) Tags() []string {
	return []string{"resources", "unique", "identifiers"}
}

func (r *E3019) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Track identifiers by resource type
	identifiersByType := make(map[string]map[string]string) // resourceType -> identifier -> resourceName

	for resName, res := range tmpl.Resources {
		// Skip resources without schemas or with intrinsic functions
		rt, err := schema.GetResourceType(res.Type)
		if err != nil || rt == nil {
			continue
		}

		// For now, we'll use a simple heuristic: check common identifier properties
		// A more complete implementation would use the schema's primaryIdentifier field
		identifierProps := []string{
			"Name", "BucketName", "TableName", "FunctionName", "QueueName",
			"TopicName", "RoleName", "PolicyName", "GroupName", "UserName",
			"KeyName", "StreamName", "ClusterName", "DBInstanceIdentifier",
		}

		var identifier string
		for _, prop := range identifierProps {
			if val, ok := res.Properties[prop]; ok {
				// Skip intrinsic functions
				if isIntrinsicFunction(val) {
					continue
				}
				if strVal, ok := val.(string); ok {
					identifier = strVal
					break
				}
			}
		}

		if identifier == "" {
			continue
		}

		// Check for duplicates within same resource type
		if identifiersByType[res.Type] == nil {
			identifiersByType[res.Type] = make(map[string]string)
		}

		if existingRes, exists := identifiersByType[res.Type][identifier]; exists {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Resource '%s' has duplicate identifier '%s' with resource '%s'", resName, identifier, existingRes),
				Line:    res.Node.Line,
				Column:  res.Node.Column,
				Path:    []string{"Resources", resName},
			})
		} else {
			identifiersByType[res.Type][identifier] = resName
		}
	}

	return matches
}
