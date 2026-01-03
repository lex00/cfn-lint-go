package informational

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&I3510{})
}

// I3510 suggests that IAM policy statement resources should match actions.
type I3510 struct{}

func (r *I3510) ID() string { return "I3510" }

func (r *I3510) ShortDesc() string {
	return "IAM statement resources match actions"
}

func (r *I3510) Description() string {
	return "Suggests that IAM policy statement resources should match the service specified in the actions for better clarity and security."
}

func (r *I3510) Source() string {
	return "https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_resource.html"
}

func (r *I3510) Tags() []string {
	return []string{"iam", "security", "best-practice"}
}

func (r *I3510) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		// Check IAM Roles
		if res.Type == "AWS::IAM::Role" {
			checkPolicyDocument(res.Properties, "AssumeRolePolicyDocument", resName, &matches)
			if policies, ok := res.Properties["Policies"].([]any); ok {
				for i, policy := range policies {
					if policyMap, ok := policy.(map[string]any); ok {
						checkPolicyDocument(policyMap, "PolicyDocument", fmt.Sprintf("%s.Policies[%d]", resName, i), &matches)
					}
				}
			}
		}

		// Check IAM Policies
		if res.Type == "AWS::IAM::Policy" || res.Type == "AWS::IAM::ManagedPolicy" {
			checkPolicyDocument(res.Properties, "PolicyDocument", resName, &matches)
		}

		// Check IAM Users
		if res.Type == "AWS::IAM::User" {
			if policies, ok := res.Properties["Policies"].([]any); ok {
				for i, policy := range policies {
					if policyMap, ok := policy.(map[string]any); ok {
						checkPolicyDocument(policyMap, "PolicyDocument", fmt.Sprintf("%s.Policies[%d]", resName, i), &matches)
					}
				}
			}
		}

		// Check IAM Groups
		if res.Type == "AWS::IAM::Group" {
			if policies, ok := res.Properties["Policies"].([]any); ok {
				for i, policy := range policies {
					if policyMap, ok := policy.(map[string]any); ok {
						checkPolicyDocument(policyMap, "PolicyDocument", fmt.Sprintf("%s.Policies[%d]", resName, i), &matches)
					}
				}
			}
		}
	}

	return matches
}

func checkPolicyDocument(props map[string]any, docKey string, resName string, matches *[]rules.Match) {
	policyDoc, ok := props[docKey].(map[string]any)
	if !ok {
		return
	}

	statements, ok := policyDoc["Statement"].([]any)
	if !ok {
		return
	}

	for i, stmt := range statements {
		stmtMap, ok := stmt.(map[string]any)
		if !ok {
			continue
		}

		// Skip Deny statements as they might intentionally have broad resources
		if effect, ok := stmtMap["Effect"].(string); ok && effect == "Deny" {
			continue
		}

		// Get actions
		actions := extractActions(stmtMap)
		if len(actions) == 0 {
			continue
		}

		// Get resources
		resources := extractResources(stmtMap)
		if len(resources) == 0 {
			continue
		}

		// Check if resources match actions
		actionServices := extractServices(actions)
		resourceServices := extractServicesFromResources(resources)

		if len(actionServices) > 0 && len(resourceServices) > 0 {
			if !servicesMatch(actionServices, resourceServices) {
				*matches = append(*matches, rules.Match{
					Message: fmt.Sprintf("IAM policy statement %d in %s has actions for services %v but resources for services %v. Consider aligning resources with actions for clarity.",
						i, resName, setToSlice(actionServices), setToSlice(resourceServices)),
					Path: []string{"Resources", resName, docKey, "Statement"},
				})
			}
		}
	}
}

func extractActions(stmt map[string]any) []string {
	var actions []string

	if action, ok := stmt["Action"]; ok {
		switch v := action.(type) {
		case string:
			actions = append(actions, v)
		case []any:
			for _, a := range v {
				if str, ok := a.(string); ok {
					actions = append(actions, str)
				}
			}
		}
	}

	return actions
}

func extractResources(stmt map[string]any) []string {
	var resources []string

	if resource, ok := stmt["Resource"]; ok {
		switch v := resource.(type) {
		case string:
			if v != "*" {
				resources = append(resources, v)
			}
		case []any:
			for _, r := range v {
				if str, ok := r.(string); ok && str != "*" {
					resources = append(resources, str)
				}
			}
		}
	}

	return resources
}

func extractServices(actions []string) map[string]bool {
	services := make(map[string]bool)
	for _, action := range actions {
		if action == "*" {
			continue
		}
		parts := strings.Split(action, ":")
		if len(parts) >= 1 {
			services[parts[0]] = true
		}
	}
	return services
}

func extractServicesFromResources(resources []string) map[string]bool {
	services := make(map[string]bool)
	for _, resource := range resources {
		if !strings.HasPrefix(resource, "arn:") {
			continue
		}
		// ARN format: arn:partition:service:region:account-id:...
		parts := strings.Split(resource, ":")
		if len(parts) >= 3 {
			services[parts[2]] = true
		}
	}
	return services
}

func servicesMatch(actionServices, resourceServices map[string]bool) bool {
	// Allow if either is empty (wildcards)
	if len(actionServices) == 0 || len(resourceServices) == 0 {
		return true
	}

	// Check if there's at least one matching service
	for service := range actionServices {
		if resourceServices[service] {
			return true
		}
	}

	return false
}

func setToSlice(set map[string]bool) []string {
	var result []string
	for key := range set {
		result = append(result, key)
	}
	return result
}
