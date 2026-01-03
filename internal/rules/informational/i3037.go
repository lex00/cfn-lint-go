package informational

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&I3037{})
}

// I3037 checks for duplicates in lists that allow duplicates.
type I3037 struct{}

func (r *I3037) ID() string { return "I3037" }

func (r *I3037) ShortDesc() string {
	return "Duplicates in list that allows duplicates"
}

func (r *I3037) Description() string {
	return "Checks for duplicate values in lists that technically allow duplicates but where duplicates are likely unintentional (security groups, subnets, etc.)."
}

func (r *I3037) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint"
}

func (r *I3037) Tags() []string {
	return []string{"resources", "best-practice", "duplicates"}
}

// Properties that commonly shouldn't have duplicates even if technically allowed
var duplicateCheckProperties = map[string][]string{
	"AWS::EC2::Instance": {
		"SecurityGroupIds",
		"SecurityGroups",
	},
	"AWS::RDS::DBInstance": {
		"VPCSecurityGroups",
		"DBSecurityGroups",
	},
	"AWS::RDS::DBCluster": {
		"VpcSecurityGroupIds",
	},
	"AWS::Lambda::Function": {
		"SubnetIds", // via VpcConfig
		"SecurityGroupIds",
	},
	"AWS::ECS::Service": {
		"Subnets",
		"SecurityGroups",
	},
	"AWS::ElasticLoadBalancingV2::LoadBalancer": {
		"Subnets",
		"SecurityGroups",
	},
	"AWS::AutoScaling::AutoScalingGroup": {
		"VPCZoneIdentifier",
		"AvailabilityZones",
	},
}

func (r *I3037) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		propsToCheck, ok := duplicateCheckProperties[res.Type]
		if !ok {
			continue
		}

		for _, propName := range propsToCheck {
			// For Lambda functions, check nested VpcConfig properties
			if res.Type == "AWS::Lambda::Function" && (propName == "SubnetIds" || propName == "SecurityGroupIds") {
				if vpcConfig, ok := res.Properties["VpcConfig"].(map[string]any); ok {
					checkForDuplicates(vpcConfig, propName, []string{"Resources", resName, "Properties", "VpcConfig", propName}, &matches)
				}
			} else {
				// For all other resources, check the property directly
				checkForDuplicates(res.Properties, propName, []string{"Resources", resName, "Properties", propName}, &matches)
			}
		}
	}

	return matches
}

func checkForDuplicates(props map[string]any, propName string, path []string, matches *[]rules.Match) {
	if propVal, exists := props[propName]; exists {
		if arr, ok := propVal.([]any); ok {
			seen := make(map[string]int)
			for i, val := range arr {
				if strVal, ok := val.(string); ok {
					if prevIdx, duplicate := seen[strVal]; duplicate {
						*matches = append(*matches, rules.Match{
							Message: fmt.Sprintf("Duplicate value '%s' found at indices %d and %d. This is likely unintentional.", strVal, prevIdx, i),
							Path:    path,
						})
					}
					seen[strVal] = i
				}
			}
		}
	}
}
