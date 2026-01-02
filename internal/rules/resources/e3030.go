// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cloudformation-schema-go/enums"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3030{})
}

// E3030 checks that resource property values are valid enum values.
// Uses enum definitions from cloudformation-schema-go.
type E3030 struct{}

func (r *E3030) ID() string { return "E3030" }

func (r *E3030) ShortDesc() string {
	return "Invalid enum value"
}

func (r *E3030) Description() string {
	return "Checks that resource property values match allowed enum values from CloudFormation schemas."
}

func (r *E3030) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3030"
}

func (r *E3030) Tags() []string {
	return []string{"resources", "properties", "enum"}
}

func (r *E3030) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		// Extract service name from resource type (e.g., "lambda" from "AWS::Lambda::Function")
		service := extractServiceName(res.Type)
		if service == "" {
			continue
		}

		for propName, propValue := range res.Properties {
			// Skip non-string values and intrinsic functions
			strValue, ok := propValue.(string)
			if !ok {
				continue
			}

			// Check if this property has an enum mapping
			enumName := enums.GetEnumForProperty(service, propName)
			if enumName == "" {
				continue
			}

			// Validate the value against allowed enum values
			if !enums.IsValidValue(service, enumName, strValue) {
				allowedValues := enums.GetAllowedValues(service, enumName)
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Property '%s' in resource '%s' (%s) has invalid value '%s'. Allowed values: %v",
						propName, resName, res.Type, strValue, allowedValues,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", propName},
				})
			}
		}
	}

	return matches
}

// extractServiceName extracts the lowercase service name from a CloudFormation resource type.
// For example, "AWS::Lambda::Function" returns "lambda".
func extractServiceName(resourceType string) string {
	// Handle standard AWS resources: AWS::ServiceName::ResourceName
	if !strings.HasPrefix(resourceType, "AWS::") {
		return ""
	}

	parts := strings.Split(resourceType, "::")
	if len(parts) < 3 {
		return ""
	}

	return strings.ToLower(parts[1])
}
