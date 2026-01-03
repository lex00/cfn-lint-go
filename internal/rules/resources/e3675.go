// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3675{})
}

// E3675 validates EMR cluster instance type.
type E3675 struct{}

func (r *E3675) ID() string { return "E3675" }

func (r *E3675) ShortDesc() string {
	return "Validate EMR cluster instance type"
}

func (r *E3675) Description() string {
	return "Validates that AWS::EMR::Cluster and InstanceGroupConfig resources specify valid instance types."
}

func (r *E3675) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3675"
}

func (r *E3675) Tags() []string {
	return []string{"resources", "properties", "emr", "instancetype"}
}

func (r *E3675) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// EMR supports standard EC2 instance types
	validFamilies := []string{
		"m1.", "m2.", "m3.", "m4.", "m5.", "m5a.", "m5d.", "m5n.", "m6g.", "m6gd.", "m6i.", "m6id.",
		"c1.", "c3.", "c4.", "c5.", "c5d.", "c5n.", "c6g.", "c6gd.", "c6i.", "c6id.",
		"r3.", "r4.", "r5.", "r5a.", "r5b.", "r5d.", "r5n.", "r6g.", "r6gd.", "r6i.", "r6id.",
		"i2.", "i3.", "i3en.", "i4i.",
		"d2.", "d3.", "d3en.",
		"h1.",
		"g3.", "g4dn.", "g5.",
		"p2.", "p3.", "p4d.",
		"z1d.",
	}

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::EMR::Cluster" && res.Type != "AWS::EMR::InstanceGroupConfig" {
			continue
		}

		var instanceType any
		var hasInstanceType bool

		if res.Type == "AWS::EMR::Cluster" {
			// Check Instances configuration
			instances, hasInstances := res.Properties["Instances"]
			if !hasInstances || isIntrinsicFunction(instances) {
				continue
			}

			instancesMap, ok := instances.(map[string]any)
			if !ok {
				continue
			}

			// Check MasterInstanceType
			masterInstanceType, hasMaster := instancesMap["MasterInstanceType"]
			if hasMaster && !isIntrinsicFunction(masterInstanceType) {
				instanceType = masterInstanceType
				hasInstanceType = true
			}

			// Check CoreInstanceType
			coreInstanceType, hasCore := instancesMap["CoreInstanceType"]
			if hasCore && !isIntrinsicFunction(coreInstanceType) {
				instanceType = coreInstanceType
				hasInstanceType = true
			}
		} else {
			instanceType, hasInstanceType = res.Properties["InstanceType"]
		}

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
					"Resource '%s': Invalid EMR instance type '%s'. Must be a valid EC2 instance type",
					resName, instanceTypeStr,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties"},
			})
		}
	}

	return matches
}
