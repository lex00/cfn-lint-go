// Package formats contains format validation rules (E11xx).
package formats

import (
	"fmt"
	"regexp"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1153{})
}

// E1153 validates security group name format.
type E1153 struct{}

func (r *E1153) ID() string { return "E1153" }

func (r *E1153) ShortDesc() string {
	return "Security group name format validation"
}

func (r *E1153) Description() string {
	return "Validates that security group names don't start with 'sg-' (reserved for IDs)."
}

func (r *E1153) Source() string {
	return "https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-security-groups.html"
}

func (r *E1153) Tags() []string {
	return []string{"format", "security-group"}
}

var securityGroupNameInvalidPattern = regexp.MustCompile(`^sg-`)

func (r *E1153) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		// Only check EC2 SecurityGroup resources
		if res.Type == "AWS::EC2::SecurityGroup" {
			if groupName := getStringProperty(res.Properties, "GroupName"); groupName != "" {
				if securityGroupNameInvalidPattern.MatchString(groupName) {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Security group name '%s' in resource '%s' cannot start with 'sg-' (reserved for security group IDs)", groupName, resName),
						Path:    []string{"Resources", resName, "Properties", "GroupName"},
					})
				}
			}
		}
	}

	return matches
}

func getStringProperty(props map[string]any, propName string) string {
	if val, ok := props[propName]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
