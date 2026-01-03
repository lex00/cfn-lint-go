// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3621{})
}

// E3621 validates AppStream Fleet instance types.
type E3621 struct{}

func (r *E3621) ID() string { return "E3621" }

func (r *E3621) ShortDesc() string {
	return "Validate AppStream Fleet instance types"
}

func (r *E3621) Description() string {
	return "Validates that AWS::AppStream::Fleet resources specify valid instance types."
}

func (r *E3621) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3621"
}

func (r *E3621) Tags() []string {
	return []string{"resources", "properties", "appstream", "instancetype"}
}

func (r *E3621) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::AppStream::Fleet" {
			continue
		}

		instanceType, hasInstanceType := res.Properties["InstanceType"]
		if !hasInstanceType || isIntrinsicFunction(instanceType) {
			continue
		}

		instanceTypeStr, ok := instanceType.(string)
		if !ok {
			continue
		}

		// AppStream instance types start with stream.
		if !strings.HasPrefix(instanceTypeStr, "stream.") {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Invalid AppStream Fleet instance type '%s'. Must start with 'stream.'",
					resName, instanceTypeStr,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties", "InstanceType"},
			})
		}
	}

	return matches
}
