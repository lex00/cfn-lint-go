package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W3660{})
}

// W3660 warns when multiple resources modify the same RestApi.
type W3660 struct{}

func (r *W3660) ID() string { return "W3660" }

func (r *W3660) ShortDesc() string {
	return "Multiple resources modifying RestApi"
}

func (r *W3660) Description() string {
	return "Warns when multiple resources (such as Methods, Resources, or Deployments) reference the same RestApi, which may cause deployment ordering issues."
}

func (r *W3660) Source() string {
	return "https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-rest-api.html"
}

func (r *W3660) Tags() []string {
	return []string{"warnings", "apigateway", "rest-api", "deployment"}
}

// Resources that modify a RestApi
var restApiModifiers = map[string]string{
	"AWS::ApiGateway::Resource":         "RestApiId",
	"AWS::ApiGateway::Method":           "RestApiId",
	"AWS::ApiGateway::Deployment":       "RestApiId",
	"AWS::ApiGateway::Stage":            "RestApiId",
	"AWS::ApiGateway::Authorizer":       "RestApiId",
	"AWS::ApiGateway::Model":            "RestApiId",
	"AWS::ApiGateway::GatewayResponse":  "RestApiId",
	"AWS::ApiGateway::RequestValidator": "RestApiId",
}

func (r *W3660) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Build a map of RestApi -> modifying resources
	apiModifiers := make(map[string][]string)

	for resName, res := range tmpl.Resources {
		propName, isModifier := restApiModifiers[res.Type]
		if !isModifier {
			continue
		}

		apiId := r.getRestApiId(res.Properties[propName])
		if apiId != "" {
			apiModifiers[apiId] = append(apiModifiers[apiId], resName)
		}
	}

	// Check for RestApis with multiple modifiers but no explicit DependsOn
	for apiId, modifiers := range apiModifiers {
		if len(modifiers) <= 1 {
			continue
		}

		// Check for Deployment resources without proper DependsOn
		for _, modifier := range modifiers {
			res := tmpl.Resources[modifier]
			if res.Type != "AWS::ApiGateway::Deployment" {
				continue
			}

			// Check if deployment has DependsOn for other API resources
			hasDependsOn := len(res.DependsOn) > 0
			if !hasDependsOn {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("API Gateway Deployment '%s' modifies RestApi '%s' but has no DependsOn; this may cause deployment ordering issues", modifier, apiId),
					Path:    []string{"Resources", modifier},
				})
			}
		}

		// Warn about complex API configurations
		if len(modifiers) > 5 {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("RestApi '%s' is modified by %d resources; consider using AWS::ApiGateway::RestApi with embedded definitions or OpenAPI for complex APIs", apiId, len(modifiers)),
				Path:    []string{"Resources", apiId},
			})
		}
	}

	return matches
}

func (r *W3660) getRestApiId(v any) string {
	// Direct reference to a RestApi resource
	if ref, ok := v.(map[string]any); ok {
		if refName, ok := ref["Ref"].(string); ok {
			return refName
		}
	}

	// String value (could be a parameter or direct ID)
	if str, ok := v.(string); ok {
		return str
	}

	return ""
}
