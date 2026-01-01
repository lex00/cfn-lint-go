// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3003{})
}

// E3003 checks that resources have required properties.
// Note: Full implementation requires CloudFormation resource schemas.
// This is a simplified version that checks common required properties.
type E3003 struct{}

func (r *E3003) ID() string { return "E3003" }

func (r *E3003) ShortDesc() string {
	return "Required properties are present"
}

func (r *E3003) Description() string {
	return "Checks that resources have their required properties. Full validation requires CloudFormation resource schemas."
}

func (r *E3003) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3003"
}

func (r *E3003) Tags() []string {
	return []string{"resources", "required", "properties"}
}

// requiredProperties maps resource types to their required properties.
// This is a subset of the most common resources. Full implementation
// would load this from CloudFormation resource schemas.
var requiredProperties = map[string][]string{
	"AWS::Lambda::Function":            {"Role", "Code"},
	"AWS::IAM::Role":                   {"AssumeRolePolicyDocument"},
	"AWS::S3::Bucket":                  {}, // No required properties
	"AWS::EC2::Instance":               {"ImageId"},
	"AWS::EC2::SecurityGroup":          {"GroupDescription"},
	"AWS::EC2::VPC":                    {"CidrBlock"},
	"AWS::EC2::Subnet":                 {"VpcId", "CidrBlock"},
	"AWS::SNS::Topic":                  {}, // No required properties
	"AWS::SQS::Queue":                  {}, // No required properties
	"AWS::DynamoDB::Table":             {"KeySchema", "AttributeDefinitions"},
	"AWS::RDS::DBInstance":             {"DBInstanceClass", "Engine"},
	"AWS::ECS::TaskDefinition":         {"ContainerDefinitions"},
	"AWS::ECS::Service":                {"TaskDefinition"},
	"AWS::CloudWatch::Alarm":           {"ComparisonOperator", "EvaluationPeriods", "MetricName", "Namespace", "Period", "Threshold"},
	"AWS::ApiGateway::RestApi":         {},
	"AWS::Logs::LogGroup":              {},
	"AWS::StepFunctions::StateMachine": {"RoleArn", "DefinitionString"},
}

func (r *E3003) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		required, known := requiredProperties[res.Type]
		if !known {
			// Unknown resource type - skip validation
			continue
		}

		for _, prop := range required {
			if _, exists := res.Properties[prop]; !exists {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Resource '%s' (%s) is missing required property '%s'", resName, res.Type, prop),
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	return matches
}
