package warnings

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W2506{})
}

// W2506 warns when ImageId parameters don't use the AWS::EC2::Image::Id type.
type W2506 struct{}

func (r *W2506) ID() string { return "W2506" }

func (r *W2506) ShortDesc() string {
	return "ImageId parameter type"
}

func (r *W2506) Description() string {
	return "Warns when a parameter that appears to be an AMI ID doesn't use the AWS::EC2::Image::Id parameter type."
}

func (r *W2506) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/parameters-section-structure.html#aws-specific-parameter-types"
}

func (r *W2506) Tags() []string {
	return []string{"warnings", "parameters", "ec2", "ami"}
}

func (r *W2506) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for paramName, param := range tmpl.Parameters {
		lowerName := strings.ToLower(paramName)

		// Check if parameter name suggests it's an AMI ID
		isImageIdLike := strings.Contains(lowerName, "imageid") ||
			strings.Contains(lowerName, "image_id") ||
			strings.Contains(lowerName, "amiid") ||
			strings.Contains(lowerName, "ami_id") ||
			strings.Contains(lowerName, "ami") && strings.Contains(lowerName, "id")

		if !isImageIdLike {
			continue
		}

		// Check if the type is not AWS::EC2::Image::Id
		if param.Type != "AWS::EC2::Image::Id" && param.Type != "AWS::SSM::Parameter::Value<AWS::EC2::Image::Id>" {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Parameter '%s' appears to be an AMI ID but uses type '%s'; consider using AWS::EC2::Image::Id for validation", paramName, param.Type),
				Path:    []string{"Parameters", paramName, "Type"},
			})
		}

		// Check if there's a default that looks like an AMI ID
		if param.Default != nil {
			defaultStr := fmt.Sprintf("%v", param.Default)
			if strings.HasPrefix(defaultStr, "ami-") && param.Type == "String" {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Parameter '%s' has a default AMI ID but uses String type; consider using AWS::EC2::Image::Id", paramName),
					Path:    []string{"Parameters", paramName},
				})
			}
		}
	}

	return matches
}
