package informational

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/sam"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&I3101{})
}

// I3101 provides informational messages about SAM resource expansion.
type I3101 struct{}

func (r *I3101) ID() string { return "I3101" }

func (r *I3101) ShortDesc() string {
	return "SAM resource expansion info"
}

func (r *I3101) Description() string {
	return "Provides informational messages about AWS::Serverless resources that will be expanded into multiple CloudFormation resources during SAM transformation. This helps users understand the resources that will be created."
}

func (r *I3101) Source() string {
	return "https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-specification-resources-and-properties.html"
}

func (r *I3101) Tags() []string {
	return []string{"resources", "sam", "serverless", "informational"}
}

// samResourceExpansions maps SAM resource types to their typical expansion descriptions.
var samResourceExpansions = map[string]string{
	"AWS::Serverless::Function":     "expands to Lambda Function, IAM Role, and optionally API Gateway resources",
	"AWS::Serverless::Api":          "expands to API Gateway RestApi, Deployment, and Stage resources",
	"AWS::Serverless::HttpApi":      "expands to API Gateway V2 Api and Stage resources",
	"AWS::Serverless::SimpleTable":  "expands to DynamoDB Table resource",
	"AWS::Serverless::LayerVersion": "expands to Lambda LayerVersion resource",
	"AWS::Serverless::Application":  "expands to nested CloudFormation stack",
	"AWS::Serverless::StateMachine": "expands to Step Functions StateMachine and IAM Role resources",
	"AWS::Serverless::Connector":    "expands to IAM policies connecting source and destination resources",
	"AWS::Serverless::GraphQLApi":   "expands to AppSync GraphQLApi and related resources",
}

func (r *I3101) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if !sam.IsSAMResourceType(res.Type) {
			continue
		}

		description, exists := samResourceExpansions[res.Type]
		if !exists {
			description = "expands to multiple CloudFormation resources"
		}

		line, column := 0, 0
		if res.Node != nil {
			line = res.Node.Line
			column = res.Node.Column
		}

		matches = append(matches, rules.Match{
			Message: fmt.Sprintf("SAM resource '%s' (%s) %s during transformation.", resName, res.Type, description),
			Line:    line,
			Column:  column,
			Path:    []string{"Resources", resName},
		})
	}

	return matches
}
