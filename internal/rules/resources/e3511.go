// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"regexp"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3511{})
}

// E3511 validates IAM role ARN patterns.
type E3511 struct{}

func (r *E3511) ID() string { return "E3511" }

func (r *E3511) ShortDesc() string {
	return "IAM role ARN pattern"
}

func (r *E3511) Description() string {
	return "Validates that IAM role ARN references follow the correct format pattern."
}

func (r *E3511) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3511"
}

func (r *E3511) Tags() []string {
	return []string{"resources", "properties", "iam", "role", "arn"}
}

// IAM role ARN pattern: arn:partition:iam::account-id:role/role-name
var iamRoleARNPattern = regexp.MustCompile(`^arn:(aws|aws-cn|aws-us-gov):iam::\d{12}:role/[\w+=,.@-]+$`)

func (r *E3511) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		// Check properties that expect IAM role ARNs
		r.checkProperties(res.Properties, resName, &matches, res.Node.Line, res.Node.Column, []string{"Resources", resName, "Properties"})
	}

	return matches
}

func (r *E3511) checkProperties(props map[string]interface{}, resName string, matches *[]rules.Match, line, column int, path []string) {
	for key, value := range props {
		// Check for common IAM role ARN property names
		if r.isRoleARNProperty(key) {
			if arnStr, ok := value.(string); ok {
				// Only validate if it looks like an ARN (not a Ref or other intrinsic)
				if len(arnStr) > 4 && arnStr[:4] == "arn:" {
					if !iamRoleARNPattern.MatchString(arnStr) {
						*matches = append(*matches, rules.Match{
							Message: fmt.Sprintf(
								"Resource '%s': Property '%s' has invalid IAM role ARN format '%s'",
								resName, key, arnStr,
							),
							Line:   line,
							Column: column,
							Path:   append(path, key),
						})
					}
				}
			}
		}

		// Recursively check nested objects
		switch v := value.(type) {
		case map[string]interface{}:
			r.checkProperties(v, resName, matches, line, column, append(path, key))
		case []interface{}:
			for i, item := range v {
				if itemMap, ok := item.(map[string]interface{}); ok {
					r.checkProperties(itemMap, resName, matches, line, column, append(path, key, fmt.Sprintf("[%d]", i)))
				}
			}
		}
	}
}

func (r *E3511) isRoleARNProperty(propName string) bool {
	roleARNProps := []string{
		"RoleArn",
		"RoleARN",
		"ExecutionRoleArn",
		"TaskRoleArn",
		"ServiceRoleArn",
		"IamRoleArn",
	}

	for _, prop := range roleARNProps {
		if propName == prop {
			return true
		}
	}
	return false
}
