package warnings

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W3037{})
}

// W3037 warns about IAM permission configuration issues.
type W3037 struct{}

func (r *W3037) ID() string { return "W3037" }

func (r *W3037) ShortDesc() string {
	return "IAM permission configuration"
}

func (r *W3037) Description() string {
	return "Warns about potential IAM permission misconfigurations, such as overly permissive managed policies or missing boundaries."
}

func (r *W3037) Source() string {
	return "https://docs.aws.amazon.com/IAM/latest/UserGuide/best-practices.html"
}

func (r *W3037) Tags() []string {
	return []string{"warnings", "iam", "security", "permissions"}
}

// Overly permissive managed policies
var overlyPermissivePolicies = map[string]string{
	"arn:aws:iam::aws:policy/AdministratorAccess": "grants full access to all AWS services",
	"arn:aws:iam::aws:policy/PowerUserAccess":     "grants broad access excluding IAM",
	"arn:aws:iam::aws:policy/IAMFullAccess":       "grants full IAM access",
}

func (r *W3037) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		switch res.Type {
		case "AWS::IAM::Role":
			r.checkRole(resName, res, &matches)
		case "AWS::IAM::User":
			r.checkUser(resName, res, &matches)
		case "AWS::IAM::Group":
			r.checkGroup(resName, res, &matches)
		}
	}

	return matches
}

func (r *W3037) checkRole(resName string, res *template.Resource, matches *[]rules.Match) {
	// Check for overly permissive managed policies
	if managedPolicies, ok := res.Properties["ManagedPolicyArns"].([]any); ok {
		for _, policy := range managedPolicies {
			if policyArn, ok := policy.(string); ok {
				if reason, found := overlyPermissivePolicies[policyArn]; found {
					*matches = append(*matches, rules.Match{
						Message: fmt.Sprintf("IAM Role '%s' uses '%s' which %s; consider using more restrictive policies", resName, policyArn, reason),
						Path:    []string{"Resources", resName, "Properties", "ManagedPolicyArns"},
					})
				}
			}
		}
	}

	// Check for missing permissions boundary
	_, hasBoundary := res.Properties["PermissionsBoundary"]
	if !hasBoundary {
		// Check if the role has admin-like policies
		if r.hasAdminLikePolicies(res) {
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf("IAM Role '%s' has broad permissions but no PermissionsBoundary; consider adding a permissions boundary", resName),
				Path:    []string{"Resources", resName, "Properties"},
			})
		}
	}

	// Check assume role policy for overly permissive principals
	if assumeRolePolicy, ok := res.Properties["AssumeRolePolicyDocument"].(map[string]any); ok {
		if statements, ok := assumeRolePolicy["Statement"].([]any); ok {
			for _, stmt := range statements {
				if stmtMap, ok := stmt.(map[string]any); ok {
					if principal, ok := stmtMap["Principal"]; ok {
						if principalStr, ok := principal.(string); ok && principalStr == "*" {
							*matches = append(*matches, rules.Match{
								Message: fmt.Sprintf("IAM Role '%s' allows any principal to assume it; consider restricting the Principal", resName),
								Path:    []string{"Resources", resName, "Properties", "AssumeRolePolicyDocument"},
							})
						}
					}
				}
			}
		}
	}
}

func (r *W3037) checkUser(resName string, res *template.Resource, matches *[]rules.Match) {
	// Warn about creating IAM users (prefer roles)
	*matches = append(*matches, rules.Match{
		Message: fmt.Sprintf("Resource '%s' creates an IAM User; consider using IAM Roles with temporary credentials instead", resName),
		Path:    []string{"Resources", resName},
	})

	// Check for overly permissive managed policies
	if managedPolicies, ok := res.Properties["ManagedPolicyArns"].([]any); ok {
		for _, policy := range managedPolicies {
			if policyArn, ok := policy.(string); ok {
				if reason, found := overlyPermissivePolicies[policyArn]; found {
					*matches = append(*matches, rules.Match{
						Message: fmt.Sprintf("IAM User '%s' uses '%s' which %s; consider using more restrictive policies", resName, policyArn, reason),
						Path:    []string{"Resources", resName, "Properties", "ManagedPolicyArns"},
					})
				}
			}
		}
	}

	// Check for inline policies with wildcards
	if policies, ok := res.Properties["Policies"].([]any); ok {
		r.checkInlinePolicies(resName, policies, matches)
	}
}

func (r *W3037) checkGroup(resName string, res *template.Resource, matches *[]rules.Match) {
	// Check for overly permissive managed policies
	if managedPolicies, ok := res.Properties["ManagedPolicyArns"].([]any); ok {
		for _, policy := range managedPolicies {
			if policyArn, ok := policy.(string); ok {
				if reason, found := overlyPermissivePolicies[policyArn]; found {
					*matches = append(*matches, rules.Match{
						Message: fmt.Sprintf("IAM Group '%s' uses '%s' which %s; consider using more restrictive policies", resName, policyArn, reason),
						Path:    []string{"Resources", resName, "Properties", "ManagedPolicyArns"},
					})
				}
			}
		}
	}
}

func (r *W3037) hasAdminLikePolicies(res *template.Resource) bool {
	if managedPolicies, ok := res.Properties["ManagedPolicyArns"].([]any); ok {
		for _, policy := range managedPolicies {
			if policyArn, ok := policy.(string); ok {
				if _, found := overlyPermissivePolicies[policyArn]; found {
					return true
				}
			}
		}
	}
	return false
}

func (r *W3037) checkInlinePolicies(resName string, policies []any, matches *[]rules.Match) {
	for _, policy := range policies {
		policyMap, ok := policy.(map[string]any)
		if !ok {
			continue
		}

		policyDoc, ok := policyMap["PolicyDocument"].(map[string]any)
		if !ok {
			continue
		}

		statements, ok := policyDoc["Statement"].([]any)
		if !ok {
			continue
		}

		for _, stmt := range statements {
			stmtMap, ok := stmt.(map[string]any)
			if !ok {
				continue
			}

			// Check for wildcard actions
			if action, ok := stmtMap["Action"]; ok {
				if actionStr, ok := action.(string); ok && actionStr == "*" {
					*matches = append(*matches, rules.Match{
						Message: fmt.Sprintf("Inline policy for '%s' allows all actions (*); consider using specific actions", resName),
						Path:    []string{"Resources", resName, "Properties", "Policies"},
					})
				}
			}

			// Check for wildcard resources with Allow effect
			if effect, ok := stmtMap["Effect"].(string); ok && strings.EqualFold(effect, "Allow") {
				if resource, ok := stmtMap["Resource"]; ok {
					if resStr, ok := resource.(string); ok && resStr == "*" {
						*matches = append(*matches, rules.Match{
							Message: fmt.Sprintf("Inline policy for '%s' allows access to all resources (*); consider restricting resources", resName),
							Path:    []string{"Resources", resName, "Properties", "Policies"},
						})
					}
				}
			}
		}
	}
}
