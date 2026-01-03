// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"encoding/json"
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3512{})
}

// E3512 validates resource-based IAM policies.
type E3512 struct{}

func (r *E3512) ID() string { return "E3512" }

func (r *E3512) ShortDesc() string {
	return "Resource-based IAM policies"
}

func (r *E3512) Description() string {
	return "Validates that embedded IAM policies for resource-based access control have proper structure including Principal."
}

func (r *E3512) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3512"
}

func (r *E3512) Tags() []string {
	return []string{"resources", "properties", "iam", "policy", "resource"}
}

func (r *E3512) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Resources with resource-based policies
	resourceBasedPolicies := map[string]string{
		"AWS::S3::BucketPolicy":                     "PolicyDocument",
		"AWS::SQS::QueuePolicy":                     "PolicyDocument",
		"AWS::SNS::TopicPolicy":                     "PolicyDocument",
		"AWS::KMS::Key":                             "KeyPolicy",
		"AWS::Lambda::Permission":                   "", // Special case - uses different properties
		"AWS::SecretsManager::SecretResourcePolicy": "ResourcePolicy",
	}

	for resName, res := range tmpl.Resources {
		policyProp, hasType := resourceBasedPolicies[res.Type]
		if !hasType {
			continue
		}

		// Skip Lambda::Permission (it uses different structure)
		if res.Type == "AWS::Lambda::Permission" {
			continue
		}

		// Get the policy property
		policyDoc, hasPolicy := res.Properties[policyProp]
		if !hasPolicy {
			continue
		}

		r.validateResourceBasedPolicy(policyDoc, resName, &matches, res.Node.Line, res.Node.Column, []string{"Resources", resName, "Properties", policyProp})
	}

	return matches
}

func (r *E3512) validateResourceBasedPolicy(policyDoc interface{}, resName string, matches *[]rules.Match, line, column int, path []string) {
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

				// Resource-based policies MUST have Principal
				if _, hasPrincipal := stmtMap["Principal"]; !hasPrincipal {
					if _, hasNotPrincipal := stmtMap["NotPrincipal"]; !hasNotPrincipal {
						*matches = append(*matches, rules.Match{
							Message: fmt.Sprintf(
								"Resource '%s': Resource-based PolicyDocument Statement %d must include Principal or NotPrincipal",
								resName, i,
							),
							Line:   line,
							Column: column,
							Path:   append(path, "Statement", fmt.Sprintf("[%d]", i)),
						})
					}
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
			}
		}
	}
}
