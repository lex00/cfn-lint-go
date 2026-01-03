package parameters

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E2012{})
}

// E2012 checks that SSM parameter types are valid.
type E2012 struct{}

func (r *E2012) ID() string { return "E2012" }

func (r *E2012) ShortDesc() string {
	return "Parameter Type validation with SSM types"
}

func (r *E2012) Description() string {
	return "Checks that parameters using SSM types have valid format and constraints."
}

func (r *E2012) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E2012"
}

func (r *E2012) Tags() []string {
	return []string{"parameters", "type", "ssm"}
}

// Valid CloudFormation parameter types
var validParameterTypes = map[string]bool{
	// Basic types
	"String":             true,
	"Number":             true,
	"CommaDelimitedList": true,

	// AWS-specific types
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
	"List<String>":                             true,
	"List<Number>":                             true,
	"List<AWS::EC2::AvailabilityZone::Name>":   true,
	"List<AWS::EC2::Image::Id>":                true,
	"List<AWS::EC2::Instance::Id>":             true,
	"List<AWS::EC2::SecurityGroup::GroupName>": true,
	"List<AWS::EC2::SecurityGroup::Id>":        true,
	"List<AWS::EC2::Subnet::Id>":               true,
	"List<AWS::EC2::Volume::Id>":               true,
	"List<AWS::EC2::VPC::Id>":                  true,
	"List<AWS::Route53::HostedZone::Id>":       true,

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
	"AWS::SSM::Parameter::Value<AWS::Route53::HostedZone::Id>":             true,
	"AWS::SSM::Parameter::Value<List<AWS::EC2::AvailabilityZone::Name>>":   true,
	"AWS::SSM::Parameter::Value<List<AWS::EC2::Image::Id>>":                true,
	"AWS::SSM::Parameter::Value<List<AWS::EC2::Instance::Id>>":             true,
	"AWS::SSM::Parameter::Value<List<AWS::EC2::SecurityGroup::GroupName>>": true,
	"AWS::SSM::Parameter::Value<List<AWS::EC2::SecurityGroup::Id>>":        true,
	"AWS::SSM::Parameter::Value<List<AWS::EC2::Subnet::Id>>":               true,
	"AWS::SSM::Parameter::Value<List<AWS::EC2::Volume::Id>>":               true,
	"AWS::SSM::Parameter::Value<List<AWS::EC2::VPC::Id>>":                  true,
	"AWS::SSM::Parameter::Value<List<AWS::Route53::HostedZone::Id>>":       true,
}

func (r *E2012) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for name, param := range tmpl.Parameters {
		// Check if parameter type is valid
		if param.Type != "" && !validParameterTypes[param.Type] {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Parameter '%s' has invalid Type '%s'", name, param.Type),
				Line:    param.Node.Line,
				Column:  param.Node.Column,
				Path:    []string{"Parameters", name, "Type"},
			})
		}
	}

	return matches
}
