package warnings

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W2511{})
}

// W2511 warns about potential IAM policy syntax issues.
type W2511 struct{}

func (r *W2511) ID() string { return "W2511" }

func (r *W2511) ShortDesc() string {
	return "IAM policy syntax"
}

func (r *W2511) Description() string {
	return "Warns about potential issues with IAM policy document syntax, such as missing Version or overly permissive statements."
}

func (r *W2511) Source() string {
	return "https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements.html"
}

func (r *W2511) Tags() []string {
	return []string{"warnings", "iam", "security", "policy"}
}

// IAM policy resource types
var iamPolicyResourceTypes = map[string][]string{
	"AWS::IAM::Policy":               {"PolicyDocument"},
	"AWS::IAM::Role":                 {"AssumeRolePolicyDocument", "Policies"},
	"AWS::IAM::User":                 {"Policies"},
	"AWS::IAM::Group":                {"Policies"},
	"AWS::IAM::ManagedPolicy":        {"PolicyDocument"},
	"AWS::S3::BucketPolicy":          {"PolicyDocument"},
	"AWS::SQS::QueuePolicy":          {"PolicyDocument"},
	"AWS::SNS::TopicPolicy":          {"PolicyDocument"},
	"AWS::Lambda::Permission":        {},
	"AWS::KMS::Key":                  {"KeyPolicy"},
	"AWS::SecretsManager::Secret":    {"ResourcePolicy"},
	"AWS::ECR::Repository":           {"RepositoryPolicyText"},
	"AWS::Elasticsearch::Domain":     {"AccessPolicies"},
	"AWS::OpenSearchService::Domain": {"AccessPolicies"},
}

func (r *W2511) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		props, ok := iamPolicyResourceTypes[res.Type]
		if !ok {
			continue
		}

		for _, propName := range props {
			if propName == "Policies" {
				// Handle inline policies array
				if policies, ok := res.Properties["Policies"].([]any); ok {
					for i, policy := range policies {
						if policyMap, ok := policy.(map[string]any); ok {
							if policyDoc, ok := policyMap["PolicyDocument"]; ok {
								r.checkPolicyDocument(policyDoc, []string{"Resources", resName, "Properties", "Policies", fmt.Sprintf("[%d]", i), "PolicyDocument"}, &matches)
							}
						}
					}
				}
			} else {
				if policyDoc, ok := res.Properties[propName]; ok {
					r.checkPolicyDocument(policyDoc, []string{"Resources", resName, "Properties", propName}, &matches)
				}
			}
		}
	}

	return matches
}

func (r *W2511) checkPolicyDocument(doc any, path []string, matches *[]rules.Match) {
	docMap, ok := doc.(map[string]any)
	if !ok {
		return
	}

	// Check for Version field
	version, hasVersion := docMap["Version"]
	if !hasVersion {
		*matches = append(*matches, rules.Match{
			Message: "IAM policy document missing Version field; recommend adding 'Version: \"2012-10-17\"'",
			Path:    path,
		})
	} else if versionStr, isStr := version.(string); isStr && versionStr != "2012-10-17" {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("IAM policy document uses Version '%s'; recommend using '2012-10-17' for latest policy features", versionStr),
			Path:    append(path, "Version"),
		})
	}

	// Check statements
	statements, ok := docMap["Statement"].([]any)
	if !ok {
		return
	}

	for i, stmt := range statements {
		stmtMap, ok := stmt.(map[string]any)
		if !ok {
			continue
		}

		stmtPath := append(path, "Statement", fmt.Sprintf("[%d]", i))

		// Check for Effect field
		if _, hasEffect := stmtMap["Effect"]; !hasEffect {
			*matches = append(*matches, rules.Match{
				Message: "IAM policy statement missing Effect field",
				Path:    stmtPath,
			})
		}

		// Check for overly permissive actions
		if actions, ok := stmtMap["Action"]; ok {
			r.checkActions(actions, stmtMap, stmtPath, matches)
		}

		// Check for overly permissive resources
		if resources, ok := stmtMap["Resource"]; ok {
			r.checkResources(resources, stmtMap, stmtPath, matches)
		}

		// Check for overly permissive principals
		if principal, ok := stmtMap["Principal"]; ok {
			r.checkPrincipal(principal, stmtMap, stmtPath, matches)
		}
	}
}

func (r *W2511) checkActions(actions any, stmt map[string]any, path []string, matches *[]rules.Match) {
	effect, _ := stmt["Effect"].(string)
	if effect != "Allow" {
		return
	}

	checkAction := func(action string) {
		if action == "*" {
			*matches = append(*matches, rules.Match{
				Message: "IAM policy allows all actions ('*'); consider restricting to specific actions",
				Path:    append(path, "Action"),
			})
		} else if strings.HasSuffix(action, ":*") {
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf("IAM policy allows all actions for service ('%s'); consider restricting to specific actions", action),
				Path:    append(path, "Action"),
			})
		}
	}

	switch a := actions.(type) {
	case string:
		checkAction(a)
	case []any:
		for _, act := range a {
			if actStr, ok := act.(string); ok {
				checkAction(actStr)
			}
		}
	}
}

func (r *W2511) checkResources(resources any, stmt map[string]any, path []string, matches *[]rules.Match) {
	effect, _ := stmt["Effect"].(string)
	if effect != "Allow" {
		return
	}

	checkResource := func(resource string) {
		if resource == "*" {
			*matches = append(*matches, rules.Match{
				Message: "IAM policy allows access to all resources ('*'); consider restricting to specific resources",
				Path:    append(path, "Resource"),
			})
		}
	}

	switch res := resources.(type) {
	case string:
		checkResource(res)
	case []any:
		for _, r := range res {
			if resStr, ok := r.(string); ok {
				checkResource(resStr)
			}
		}
	}
}

func (r *W2511) checkPrincipal(principal any, stmt map[string]any, path []string, matches *[]rules.Match) {
	effect, _ := stmt["Effect"].(string)
	if effect != "Allow" {
		return
	}

	if principalStr, ok := principal.(string); ok && principalStr == "*" {
		*matches = append(*matches, rules.Match{
			Message: "IAM policy allows access from any principal ('*'); consider restricting to specific principals",
			Path:    append(path, "Principal"),
		})
	}

	if principalMap, ok := principal.(map[string]any); ok {
		if aws, ok := principalMap["AWS"]; ok {
			if awsStr, ok := aws.(string); ok && awsStr == "*" {
				*matches = append(*matches, rules.Match{
					Message: "IAM policy allows access from any AWS principal; consider restricting to specific accounts or roles",
					Path:    append(path, "Principal", "AWS"),
				})
			}
		}
	}
}
