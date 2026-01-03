// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"encoding/json"
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3513{})
}

// E3513 validates ECR repository policies.
type E3513 struct{}

func (r *E3513) ID() string { return "E3513" }

func (r *E3513) ShortDesc() string {
	return "ECR repository policy"
}

func (r *E3513) Description() string {
	return "Validates that ECR repository policies have proper structure and required fields."
}

func (r *E3513) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3513"
}

func (r *E3513) Tags() []string {
	return []string{"resources", "properties", "ecr", "policy"}
}

func (r *E3513) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::ECR::Repository" {
			continue
		}

		// Check for RepositoryPolicyText
		policyText, hasPolicy := res.Properties["RepositoryPolicyText"]
		if !hasPolicy {
			continue
		}

		r.validateECRPolicy(policyText, resName, &matches, res.Node.Line, res.Node.Column, []string{"Resources", resName, "Properties", "RepositoryPolicyText"})
	}

	return matches
}

func (r *E3513) validateECRPolicy(policyDoc interface{}, resName string, matches *[]rules.Match, line, column int, path []string) {
	policyMap, ok := policyDoc.(map[string]interface{})
	if !ok {
		// Try to parse as JSON string
		if policyStr, ok := policyDoc.(string); ok {
			var parsedPolicy map[string]interface{}
			if err := json.Unmarshal([]byte(policyStr), &parsedPolicy); err != nil {
				*matches = append(*matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': RepositoryPolicyText must be valid JSON: %v",
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
				"Resource '%s': RepositoryPolicyText should include Version field (recommended: '2012-10-17')",
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
				"Resource '%s': RepositoryPolicyText must include Statement field",
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
					"Resource '%s': RepositoryPolicyText Statement must be an array",
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

				// Validate required statement fields
				if _, hasEffect := stmtMap["Effect"]; !hasEffect {
					*matches = append(*matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': RepositoryPolicyText Statement %d must include Effect",
							resName, i,
						),
						Line:   line,
						Column: column,
						Path:   append(path, "Statement", fmt.Sprintf("[%d]", i)),
					})
				}

				// ECR policies require Principal
				if _, hasPrincipal := stmtMap["Principal"]; !hasPrincipal {
					*matches = append(*matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': RepositoryPolicyText Statement %d must include Principal",
							resName, i,
						),
						Line:   line,
						Column: column,
						Path:   append(path, "Statement", fmt.Sprintf("[%d]", i)),
					})
				}

				// Must have Action
				if _, hasAction := stmtMap["Action"]; !hasAction {
					*matches = append(*matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': RepositoryPolicyText Statement %d must include Action",
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
