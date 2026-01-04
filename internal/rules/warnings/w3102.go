package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W3102{})
}

// W3102 warns about SAM Api resources missing StageName property.
type W3102 struct{}

func (r *W3102) ID() string { return "W3102" }

func (r *W3102) ShortDesc() string {
	return "SAM Api missing StageName"
}

func (r *W3102) Description() string {
	return "Warns when an AWS::Serverless::Api resource does not specify StageName. Without an explicit StageName, SAM will use a default stage name. Consider explicitly setting StageName for clarity and consistency."
}

func (r *W3102) Source() string {
	return "https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-api.html"
}

func (r *W3102) Tags() []string {
	return []string{"resources", "sam", "serverless", "api-gateway"}
}

func (r *W3102) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		// Only check AWS::Serverless::Api, not HttpApi
		if res.Type != "AWS::Serverless::Api" {
			continue
		}

		// Check if StageName is set
		if _, hasStageName := res.Properties["StageName"]; hasStageName {
			continue
		}

		line, column := 0, 0
		if res.Node != nil {
			line = res.Node.Line
			column = res.Node.Column
		}

		matches = append(matches, rules.Match{
			Message: fmt.Sprintf("SAM Api '%s' does not specify StageName. Consider setting an explicit stage name for clarity.", resName),
			Line:    line,
			Column:  column,
			Path:    []string{"Resources", resName, "Properties"},
		})
	}

	return matches
}
