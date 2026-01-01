// Package parameters contains parameter validation rules (E2xxx).
package parameters

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E2002{})
}

// E2002 checks that parameter types are valid CloudFormation types.
type E2002 struct{}

func (r *E2002) ID() string { return "E2002" }

func (r *E2002) ShortDesc() string {
	return "Invalid parameter type"
}

func (r *E2002) Description() string {
	return "Checks that parameter Type is a valid CloudFormation parameter type."
}

func (r *E2002) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E2002"
}

func (r *E2002) Tags() []string {
	return []string{"parameters", "type"}
}

// Valid CloudFormation parameter types
var validParamTypes = map[string]bool{
	// Standard types
	"String":             true,
	"Number":             true,
	"List<Number>":       true,
	"CommaDelimitedList": true,

	// AWS-specific parameter types
	"AWS::EC2::AvailabilityZone::Name":   true,
	"AWS::EC2::Image::Id":                true,
	"AWS::EC2::Instance::Id":             true,
	"AWS::EC2::KeyPair::KeyName":         true,
	"AWS::EC2::SecurityGroup::GroupName": true,
	"AWS::EC2::SecurityGroup::Id":        true,
	"AWS::EC2::Subnet::Id":               true,
	"AWS::EC2::Volume::Id":               true,
	"AWS::EC2::VPC::Id":                  true,
	"AWS::Route53::HostedZone::Id":       true,

	// List types
	"List<AWS::EC2::AvailabilityZone::Name>":   true,
	"List<AWS::EC2::Image::Id>":                true,
	"List<AWS::EC2::Instance::Id>":             true,
	"List<AWS::EC2::SecurityGroup::GroupName>": true,
	"List<AWS::EC2::SecurityGroup::Id>":        true,
	"List<AWS::EC2::Subnet::Id>":               true,
	"List<AWS::EC2::Volume::Id>":               true,
	"List<AWS::EC2::VPC::Id>":                  true,

	// SSM parameter types
	"AWS::SSM::Parameter::Name":                                            true,
	"AWS::SSM::Parameter::Value<String>":                                   true,
	"AWS::SSM::Parameter::Value<List<String>>":                             true,
	"AWS::SSM::Parameter::Value<CommaDelimitedList>":                       true,
	"AWS::SSM::Parameter::Value<AWS::EC2::AvailabilityZone::Name>":         true,
	"AWS::SSM::Parameter::Value<AWS::EC2::Image::Id>":                      true,
	"AWS::SSM::Parameter::Value<AWS::EC2::Instance::Id>":                   true,
	"AWS::SSM::Parameter::Value<AWS::EC2::KeyPair::KeyName>":               true,
	"AWS::SSM::Parameter::Value<AWS::EC2::SecurityGroup::GroupName>":       true,
	"AWS::SSM::Parameter::Value<AWS::EC2::SecurityGroup::Id>":              true,
	"AWS::SSM::Parameter::Value<AWS::EC2::Subnet::Id>":                     true,
	"AWS::SSM::Parameter::Value<AWS::EC2::Volume::Id>":                     true,
	"AWS::SSM::Parameter::Value<AWS::EC2::VPC::Id>":                        true,
	"AWS::SSM::Parameter::Value<List<AWS::EC2::AvailabilityZone::Name>>":   true,
	"AWS::SSM::Parameter::Value<List<AWS::EC2::Image::Id>>":                true,
	"AWS::SSM::Parameter::Value<List<AWS::EC2::Instance::Id>>":             true,
	"AWS::SSM::Parameter::Value<List<AWS::EC2::SecurityGroup::GroupName>>": true,
	"AWS::SSM::Parameter::Value<List<AWS::EC2::SecurityGroup::Id>>":        true,
	"AWS::SSM::Parameter::Value<List<AWS::EC2::Subnet::Id>>":               true,
	"AWS::SSM::Parameter::Value<List<AWS::EC2::Volume::Id>>":               true,
	"AWS::SSM::Parameter::Value<List<AWS::EC2::VPC::Id>>":                  true,
}

func (r *E2002) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for name, param := range tmpl.Parameters {
		if param.Type == "" {
			// E2001 handles missing Type
			continue
		}

		if !validParamTypes[param.Type] {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Parameter '%s' has invalid type '%s'", name, param.Type),
				Line:    param.Node.Line,
				Column:  param.Node.Column,
				Path:    []string{"Parameters", name, "Type"},
			})
		}
	}

	return matches
}
