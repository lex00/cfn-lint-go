// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3628{})
}

// E3628 validates EC2 instance types.
type E3628 struct{}

func (r *E3628) ID() string { return "E3628" }

func (r *E3628) ShortDesc() string {
	return "Validate EC2 instance types"
}

func (r *E3628) Description() string {
	return "Validates that AWS::EC2::Instance resources specify valid instance types."
}

func (r *E3628) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3628"
}

func (r *E3628) Tags() []string {
	return []string{"resources", "properties", "ec2", "instancetype"}
}

func (r *E3628) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Common instance type families
	validFamilies := []string{
		"t2.", "t3.", "t3a.", "t4g.",
		"m4.", "m5.", "m5a.", "m5n.", "m5zn.", "m6i.", "m6a.", "m6g.", "m7g.",
		"c4.", "c5.", "c5a.", "c5n.", "c6i.", "c6a.", "c6g.", "c7g.",
		"r4.", "r5.", "r5a.", "r5b.", "r5n.", "r6i.", "r6a.", "r6g.", "r7g.",
		"x1.", "x1e.", "x2gd.", "x2idn.", "x2iedn.", "x2iezn.",
		"z1d.",
		"i3.", "i3en.", "i4i.",
		"d2.", "d3.", "d3en.",
		"h1.",
		"p2.", "p3.", "p4d.",
		"g3.", "g4dn.", "g5.",
		"f1.",
		"inf1.", "inf2.",
		"trn1.",
		"a1.",
		"u-6tb1.", "u-9tb1.", "u-12tb1.", "u-18tb1.", "u-24tb1.",
	}

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::EC2::Instance" {
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
					"Resource '%s': Invalid EC2 instance type '%s'. Must be a valid EC2 instance type",
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
