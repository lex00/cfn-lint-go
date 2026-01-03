package informational

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&I3100{})
}

// I3100 warns about legacy EC2 instance type generations.
type I3100 struct{}

func (r *I3100) ID() string { return "I3100" }

func (r *I3100) ShortDesc() string {
	return "Legacy instance type generations"
}

func (r *I3100) Description() string {
	return "Warns when using legacy or previous generation EC2 instance types. Consider using current generation instances for better performance and cost efficiency."
}

func (r *I3100) Source() string {
	return "https://aws.amazon.com/ec2/instance-types/"
}

func (r *I3100) Tags() []string {
	return []string{"resources", "ec2", "instance-types", "best-practice"}
}

// legacyInstanceFamilies contains instance families that are previous generation
var legacyInstanceFamilies = map[string]bool{
	"t1": true,
	"t2": true, // t3/t3a are current gen
	"m1": true,
	"m2": true,
	"m3": true,
	"m4": true, // m5/m6/m7 are current gen
	"c1": true,
	"c3": true,
	"c4": true, // c5/c6/c7 are current gen
	"r3": true,
	"r4": true, // r5/r6/r7 are current gen
	"i2": true,
	"i3": true, // i4 is current gen
	"d2": true,
	"g2": true,
	"g3": true, // g4/g5 are current gen
	"p2": true, // p3/p4/p5 are current gen
	"x1": true, // x2 is current gen
	"f1": true,
	"h1": true,
}

func (r *I3100) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		// Check EC2 instances
		if res.Type == "AWS::EC2::Instance" {
			if instanceType, ok := res.Properties["InstanceType"].(string); ok {
				if isLegacyInstanceType(instanceType) {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Resource '%s' uses legacy instance type '%s'. Consider using current generation instances (e.g., t3, m6, c6, r6) for better performance and cost efficiency.", resName, instanceType),
						Path:    []string{"Resources", resName, "Properties", "InstanceType"},
					})
				}
			}
		}

		// Check Launch Templates
		if res.Type == "AWS::EC2::LaunchTemplate" {
			if launchTemplateData, ok := res.Properties["LaunchTemplateData"].(map[string]any); ok {
				if instanceType, ok := launchTemplateData["InstanceType"].(string); ok {
					if isLegacyInstanceType(instanceType) {
						matches = append(matches, rules.Match{
							Message: fmt.Sprintf("Launch template '%s' uses legacy instance type '%s'. Consider using current generation instances.", resName, instanceType),
							Path:    []string{"Resources", resName, "Properties", "LaunchTemplateData", "InstanceType"},
						})
					}
				}
			}
		}

		// Check Auto Scaling Launch Configurations
		if res.Type == "AWS::AutoScaling::LaunchConfiguration" {
			if instanceType, ok := res.Properties["InstanceType"].(string); ok {
				if isLegacyInstanceType(instanceType) {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Launch configuration '%s' uses legacy instance type '%s'. Consider using current generation instances.", resName, instanceType),
						Path:    []string{"Resources", resName, "Properties", "InstanceType"},
					})
				}
			}
		}
	}

	return matches
}

func isLegacyInstanceType(instanceType string) bool {
	// Extract instance family (e.g., "m5" from "m5.large")
	parts := strings.Split(instanceType, ".")
	if len(parts) < 2 {
		return false
	}

	family := parts[0]
	return legacyInstanceFamilies[family]
}
