// Package formats contains format validation rules (E11xx).
package formats

import (
	"fmt"
	"regexp"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1150{})
}

// E1150 validates security group ID format.
type E1150 struct{}

func (r *E1150) ID() string { return "E1150" }

func (r *E1150) ShortDesc() string {
	return "Security group format validation"
}

func (r *E1150) Description() string {
	return "Validates that security group IDs match the format sg-[0-9a-zA-Z]*."
}

func (r *E1150) Source() string {
	return "https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-security-groups.html"
}

func (r *E1150) Tags() []string {
	return []string{"format", "security-group"}
}

var securityGroupPattern = regexp.MustCompile(`^sg-[0-9a-zA-Z]+$`)

func (r *E1150) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check all resources for properties that should contain security group IDs
	for resName, res := range tmpl.Resources {
		sgRefs := findSecurityGroupReferences(res.Properties, res.Type)
		for _, sg := range sgRefs {
			if sg.value != "" && !securityGroupPattern.MatchString(sg.value) {
				// Only validate if it looks like it's trying to be a security group ID
				// (starts with sg- or is a literal string that should be validated)
				if !isIntrinsicFunction(sg.rawValue) {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Invalid security group ID format '%s' in resource '%s', expected format: sg-[0-9a-zA-Z]+", sg.value, resName),
						Path:    append([]string{"Resources", resName, "Properties"}, sg.path...),
					})
				}
			}
		}
	}

	return matches
}

type securityGroupRef struct {
	value    string
	rawValue any
	path     []string
}

func findSecurityGroupReferences(v any, resourceType string) []securityGroupRef {
	var results []securityGroupRef
	findSecurityGroupReferencesRecursive(v, []string{}, &results)
	return results
}

func findSecurityGroupReferencesRecursive(v any, path []string, results *[]securityGroupRef) {
	switch val := v.(type) {
	case string:
		// Check if this looks like a security group ID
		if len(val) > 3 && val[:3] == "sg-" {
			*results = append(*results, securityGroupRef{
				value:    val,
				rawValue: v,
				path:     path,
			})
		}
	case map[string]any:
		// Skip intrinsic functions
		if isIntrinsicFunction(val) {
			return
		}
		for key, child := range val {
			// Look for common property names that contain security group IDs
			if isSecurityGroupProperty(key) {
				findSecurityGroupReferencesRecursive(child, append(path, key), results)
			} else {
				findSecurityGroupReferencesRecursive(child, append(path, key), results)
			}
		}
	case []any:
		for i, child := range val {
			findSecurityGroupReferencesRecursive(child, append(path, fmt.Sprintf("[%d]", i)), results)
		}
	}
}

func isSecurityGroupProperty(propName string) bool {
	sgProps := map[string]bool{
		"SecurityGroupId":       true,
		"SecurityGroupIds":      true,
		"SecurityGroups":        true,
		"GroupId":               true,
		"SourceSecurityGroupId": true,
	}
	return sgProps[propName]
}

func isIntrinsicFunction(v any) bool {
	if m, ok := v.(map[string]any); ok {
		for key := range m {
			if key == "Ref" || key[:4] == "Fn::" {
				return true
			}
		}
	}
	return false
}
