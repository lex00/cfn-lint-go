// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"encoding/json"
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3510{})
}

// E3510 validates identity-based IAM policies.
type E3510 struct{}

func (r *E3510) ID() string { return "E3510" }

func (r *E3510) ShortDesc() string {
	return "Identity-based IAM policies"
}

func (r *E3510) Description() string {
	return "Validates that embedded IAM policies for identity-based access control have proper structure and required fields."
}

func (r *E3510) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3510"
}

func (r *E3510) Tags() []string {
	return []string{"resources", "properties", "iam", "policy"}
}

func (r *E3510) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// IAM resources with identity-based policies
	identityResources := map[string]string{
		"AWS::IAM::User":   "Policies",
		"AWS::IAM::Group":  "Policies",
		"AWS::IAM::Role":   "Policies",
		"AWS::IAM::Policy": "PolicyDocument",
	}

	for resName, res := range tmpl.Resources {
		policyProp, hasType := identityResources[res.Type]
		if !hasType {
			continue
		}

		// Get the policies property
		policies, hasPolicies := res.Properties[policyProp]
		if !hasPolicies {
			continue
		}

		// Handle different structures
		if policyProp == "PolicyDocument" {
			// AWS::IAM::Policy has PolicyDocument as a single policy
			r.validatePolicyDocument(policies, resName, &matches, res.Node.Line, res.Node.Column, []string{"Resources", resName, "Properties", "PolicyDocument"})
		} else {
			// Other resources have Policies as an array
			policiesList, ok := policies.([]interface{})
			if !ok {
				continue
			}

			for i, policy := range policiesList {
				policyMap, ok := policy.(map[string]interface{})
				if !ok {
					continue
				}

				if policyDoc, hasDoc := policyMap["PolicyDocument"]; hasDoc {
					r.validatePolicyDocument(policyDoc, resName, &matches, res.Node.Line, res.Node.Column, []string{"Resources", resName, "Properties", "Policies", fmt.Sprintf("[%d]", i), "PolicyDocument"})
				}
			}
		}
	}

	return matches
}

func (r *E3510) validatePolicyDocument(policyDoc interface{}, resName string, matches *[]rules.Match, line, column int, path []string) {
	policyMap, ok := policyDoc.(map[string]interface{})
	if !ok {
		// Try to parse as JSON string
		if policyStr, ok := policyDoc.(string); ok {
			var parsedPolicy map[string]interface{}
			if err := json.Unmarshal([]byte(policyStr), &parsedPolicy); err != nil {
				*matches = append(*matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': PolicyDocument must be valid JSON: %v",
						resName, err,
					),
					Line:   line,
					Column: column,
					Path:   path,
				})
				return
			}
			policyMap = parsedPolicy
		} else {
			return
		}
	}

	// Validate required fields
	if _, hasVersion := policyMap["Version"]; !hasVersion {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf(
				"Resource '%s': PolicyDocument should include Version field (recommended: '2012-10-17')",
				resName,
			),
			Line:   line,
			Column: column,
			Path:   path,
		})
	}

	if _, hasStatement := policyMap["Statement"]; !hasStatement {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf(
				"Resource '%s': PolicyDocument must include Statement field",
				resName,
			),
			Line:   line,
			Column: column,
			Path:   path,
		})
	} else {
		// Validate Statement structure
		statements, ok := policyMap["Statement"].([]interface{})
		if !ok {
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': PolicyDocument Statement must be an array",
					resName,
				),
				Line:   line,
				Column: column,
				Path:   append(path, "Statement"),
			})
		} else {
			for i, stmt := range statements {
				stmtMap, ok := stmt.(map[string]interface{})
				if !ok {
					continue
				}

				// Each statement must have Effect
				if _, hasEffect := stmtMap["Effect"]; !hasEffect {
					*matches = append(*matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': PolicyDocument Statement %d must include Effect",
							resName, i,
						),
						Line:   line,
						Column: column,
						Path:   append(path, "Statement", fmt.Sprintf("[%d]", i)),
					})
				}

				// Each statement must have Action or NotAction
				if _, hasAction := stmtMap["Action"]; !hasAction {
					if _, hasNotAction := stmtMap["NotAction"]; !hasNotAction {
						*matches = append(*matches, rules.Match{
							Message: fmt.Sprintf(
								"Resource '%s': PolicyDocument Statement %d must include Action or NotAction",
								resName, i,
							),
							Line:   line,
							Column: column,
							Path:   append(path, "Statement", fmt.Sprintf("[%d]", i)),
						})
					}
				}

				// Identity-based policies should have Resource or NotResource
				if _, hasResource := stmtMap["Resource"]; !hasResource {
					if _, hasNotResource := stmtMap["NotResource"]; !hasNotResource {
						*matches = append(*matches, rules.Match{
							Message: fmt.Sprintf(
								"Resource '%s': PolicyDocument Statement %d should include Resource or NotResource",
								resName, i,
							),
							Line:   line,
							Column: column,
							Path:   append(path, "Statement", fmt.Sprintf("[%d]", i)),
						})
					}
				}
			}
		}
	}
}
