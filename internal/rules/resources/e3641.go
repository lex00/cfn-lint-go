// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3641{})
}

// E3641 validates GameLift Fleet instance type.
type E3641 struct{}

func (r *E3641) ID() string { return "E3641" }

func (r *E3641) ShortDesc() string {
	return "Validate GameLift Fleet instance type"
}

func (r *E3641) Description() string {
	return "Validates that AWS::GameLift::Fleet resources specify valid EC2 instance types."
}

func (r *E3641) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3641"
}

func (r *E3641) Tags() []string {
	return []string{"resources", "properties", "gamelift", "instancetype"}
}

func (r *E3641) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// GameLift supports specific EC2 instance types
	validFamilies := []string{
		"c3.", "c4.", "c5.", "c5a.", "c5n.", "c6g.", "c6i.",
		"m3.", "m4.", "m5.", "m5a.", "m5n.", "m6g.", "m6i.",
		"r3.", "r4.", "r5.", "r5a.", "r5n.", "r6g.", "r6i.",
	}

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::GameLift::Fleet" {
			continue
		}

		instanceType, hasInstanceType := res.Properties["EC2InstanceType"]
		if !hasInstanceType || isIntrinsicFunction(instanceType) {
			continue
		}

		instanceTypeStr, ok := instanceType.(string)
		if !ok {
			continue
		}

		// Check if instance type starts with a valid family
		isValid := false
		for _, family := range validFamilies {
			if strings.HasPrefix(instanceTypeStr, family) {
				isValid = true
				break
			}
		}

		if !isValid {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Invalid GameLift Fleet instance type '%s'. Must be a supported EC2 instance type",
					resName, instanceTypeStr,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties", "EC2InstanceType"},
			})
		}
	}

	return matches
}
