package sam

import (
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/template"

	// Import aws-sam-translator-go for future use in transformation.
	// This will be used in Phase 2 for SAM to CloudFormation transformation.
	_ "github.com/lex00/aws-sam-translator-go/pkg/translator"
)

// SAMTransform is the AWS SAM transform identifier.
const SAMTransform = "AWS::Serverless-2016-10-31"

// samResourceTypes contains all valid AWS SAM resource types.
var samResourceTypes = map[string]bool{
	"AWS::Serverless::Function":     true,
	"AWS::Serverless::Api":          true,
	"AWS::Serverless::HttpApi":      true,
	"AWS::Serverless::SimpleTable":  true,
	"AWS::Serverless::LayerVersion": true,
	"AWS::Serverless::Application":  true,
	"AWS::Serverless::StateMachine": true,
	"AWS::Serverless::Connector":    true,
	"AWS::Serverless::GraphQLApi":   true,
}

// IsSAMTemplate checks if the template is a SAM template.
// It returns true if the template contains the AWS::Serverless-2016-10-31 transform
// or any AWS::Serverless::* resource type.
func IsSAMTemplate(tmpl *template.Template) bool {
	if tmpl == nil {
		return false
	}

	// Check for SAM transform
	if hasServerlessTransform(tmpl.Transform) {
		return true
	}

	// Check for SAM resource types
	for _, res := range tmpl.Resources {
		if IsSAMResourceType(res.Type) {
			return true
		}
	}

	return false
}

// IsSAMResourceType checks if the given resource type is a SAM resource type.
func IsSAMResourceType(resourceType string) bool {
	return samResourceTypes[resourceType]
}

// hasServerlessTransform checks if the transform field contains the SAM transform.
func hasServerlessTransform(transform any) bool {
	if transform == nil {
		return false
	}

	switch t := transform.(type) {
	case string:
		return t == SAMTransform
	case []any:
		for _, item := range t {
			if str, ok := item.(string); ok && str == SAMTransform {
				return true
			}
		}
	case []string:
		for _, item := range t {
			if item == SAMTransform {
				return true
			}
		}
	}

	return false
}

// GetSAMResourceTypes returns a list of all supported SAM resource types.
func GetSAMResourceTypes() []string {
	types := make([]string, 0, len(samResourceTypes))
	for t := range samResourceTypes {
		types = append(types, t)
	}
	return types
}

// HasSAMResources checks if the template contains any SAM resources.
func HasSAMResources(tmpl *template.Template) bool {
	if tmpl == nil {
		return false
	}

	for _, res := range tmpl.Resources {
		if IsSAMResourceType(res.Type) {
			return true
		}
	}

	return false
}

// HasServerlessTransform checks if the template has the SAM transform declared.
func HasServerlessTransform(tmpl *template.Template) bool {
	if tmpl == nil {
		return false
	}
	return hasServerlessTransform(tmpl.Transform)
}

// GetSAMResources returns all SAM resources from the template.
func GetSAMResources(tmpl *template.Template) map[string]*template.Resource {
	if tmpl == nil {
		return nil
	}

	samResources := make(map[string]*template.Resource)
	for name, res := range tmpl.Resources {
		if IsSAMResourceType(res.Type) {
			samResources[name] = res
		}
	}

	return samResources
}

// IsSAMResourcePrefix checks if a resource type starts with the SAM prefix.
func IsSAMResourcePrefix(resourceType string) bool {
	return strings.HasPrefix(resourceType, "AWS::Serverless::")
}
