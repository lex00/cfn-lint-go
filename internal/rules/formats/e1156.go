// Package formats contains format validation rules (E11xx).
package formats

import (
	"fmt"
	"regexp"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1156{})
}

// E1156 validates IAM role ARN format.
type E1156 struct{}

func (r *E1156) ID() string { return "E1156" }

func (r *E1156) ShortDesc() string {
	return "IAM role ARN format validation"
}

func (r *E1156) Description() string {
	return "Validates that IAM role ARNs match the correct format."
}

func (r *E1156) Source() string {
	return "https://docs.aws.amazon.com/IAM/latest/UserGuide/reference-arns.html"
}

func (r *E1156) Tags() []string {
	return []string{"format", "iam", "arn"}
}

// IAM role ARN pattern: arn:partition:iam::account-id:role/role-name
var iamRoleARNPattern = regexp.MustCompile(`^arn:(aws|aws-cn|aws-us-gov):iam::\d{12}:role/[\w+=,.@\-/]+$`)

func (r *E1156) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		arnRefs := findIAMRoleARNReferences(res.Properties)
		for _, ref := range arnRefs {
			if ref.value != "" && !iamRoleARNPattern.MatchString(ref.value) && !isIntrinsicFunction(ref.rawValue) {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Invalid IAM role ARN format '%s' in resource '%s', expected format: arn:partition:iam::account-id:role/role-name", ref.value, resName),
					Path:    append([]string{"Resources", resName, "Properties"}, ref.path...),
				})
			}
		}
	}

	return matches
}

func findIAMRoleARNReferences(v any) []formatRef {
	var results []formatRef
	findIAMRoleARNReferencesRecursive(v, []string{}, &results)
	return results
}

func findIAMRoleARNReferencesRecursive(v any, path []string, results *[]formatRef) {
	switch val := v.(type) {
	case string:
		// Check if this looks like an IAM role ARN
		if len(val) > 9 && val[:9] == "arn:aws:i" || (len(val) > 13 && val[:13] == "arn:aws-cn:i") || (len(val) > 16 && val[:16] == "arn:aws-us-gov:i") {
			// Further check if it's specifically an IAM role ARN
			if len(val) > 15 && (val[:15] == "arn:aws:iam:::" || val[:19] == "arn:aws-cn:iam:::" || val[:22] == "arn:aws-us-gov:iam:::") {
				// This is likely an IAM ARN, check if it's a role
				matched, _ := regexp.MatchString(`:role/`, val)
				if matched {
					*results = append(*results, formatRef{
						value:    val,
						rawValue: v,
						path:     path,
					})
				}
			}
		}
	case map[string]any:
		// Skip intrinsic functions
		if isIntrinsicFunction(val) {
			return
		}
		for key, child := range val {
			// Look for common IAM role ARN property names
			if isIAMRoleARNProperty(key) {
				findIAMRoleARNReferencesRecursive(child, append(path, key), results)
			} else {
				findIAMRoleARNReferencesRecursive(child, append(path, key), results)
			}
		}
	case []any:
		for i, child := range val {
			findIAMRoleARNReferencesRecursive(child, append(path, fmt.Sprintf("[%d]", i)), results)
		}
	}
}

func isIAMRoleARNProperty(propName string) bool {
	roleProps := map[string]bool{
		"RoleArn":  true,
		"Role":     true,
		"RoleARN":  true,
		"IamRole":  true,
		"IAMRole":  true,
		"RoleName": false, // RoleName is just the name, not the ARN
	}
	return roleProps[propName]
}
