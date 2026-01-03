// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3671{})
}

// E3671 validates block device mapping configuration.
type E3671 struct{}

func (r *E3671) ID() string { return "E3671" }

func (r *E3671) ShortDesc() string {
	return "Validate block device mapping Iops requirement"
}

func (r *E3671) Description() string {
	return "Validates that EC2 block device mappings with io1, io2, or gp3 volume types specify Iops."
}

func (r *E3671) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3671"
}

func (r *E3671) Tags() []string {
	return []string{"resources", "properties", "ec2", "blockdevice"}
}

func (r *E3671) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::EC2::Instance" && res.Type != "AWS::EC2::LaunchTemplate" {
			continue
		}

		var blockDeviceMappings any
		var hasBlockDeviceMappings bool

		if res.Type == "AWS::EC2::Instance" {
			blockDeviceMappings, hasBlockDeviceMappings = res.Properties["BlockDeviceMappings"]
		} else {
			launchTemplateData, hasLaunchTemplateData := res.Properties["LaunchTemplateData"]
			if !hasLaunchTemplateData || isIntrinsicFunction(launchTemplateData) {
				continue
			}
			dataMap, ok := launchTemplateData.(map[string]any)
			if !ok {
				continue
			}
			blockDeviceMappings, hasBlockDeviceMappings = dataMap["BlockDeviceMappings"]
		}

		if !hasBlockDeviceMappings || isIntrinsicFunction(blockDeviceMappings) {
			continue
		}

		bdmList, ok := blockDeviceMappings.([]any)
		if !ok {
			continue
		}

		for _, bdm := range bdmList {
			bdmMap, ok := bdm.(map[string]any)
			if !ok {
				continue
			}

			ebs, hasEbs := bdmMap["Ebs"]
			if !hasEbs || isIntrinsicFunction(ebs) {
				continue
			}

			ebsMap, ok := ebs.(map[string]any)
			if !ok {
				continue
			}

			volumeType, hasVolumeType := ebsMap["VolumeType"]
			if !hasVolumeType || isIntrinsicFunction(volumeType) {
				continue
			}

			volumeTypeStr, ok := volumeType.(string)
			if !ok {
				continue
			}

			// io1, io2, and gp3 require Iops
			if volumeTypeStr == "io1" || volumeTypeStr == "io2" || volumeTypeStr == "gp3" {
				_, hasIops := ebsMap["Iops"]
				if !hasIops {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': EBS volume with VolumeType '%s' must specify Iops",
							resName, volumeTypeStr,
						),
						Line:   res.Node.Line,
						Column: res.Node.Column,
						Path:   []string{"Resources", resName, "Properties", "BlockDeviceMappings"},
					})
				}
			}
		}
	}

	return matches
}
