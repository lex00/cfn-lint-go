// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3660{})
}

// E3660 validates RestApi name requirement.
type E3660 struct{}

func (r *E3660) ID() string { return "E3660" }

func (r *E3660) ShortDesc() string {
	return "RestApi requires Name when no Body or BodyS3Location"
}

func (r *E3660) Description() string {
	return "Validates that AWS::ApiGateway::RestApi resources have a Name property when neither Body nor BodyS3Location is specified."
}

func (r *E3660) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3660"
}

func (r *E3660) Tags() []string {
	return []string{"resources", "properties", "apigateway", "restapi"}
}

func (r *E3660) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::ApiGateway::RestApi" {
			continue
		}

		_, hasName := res.Properties["Name"]
		_, hasBody := res.Properties["Body"]
		_, hasBodyS3Location := res.Properties["BodyS3Location"]

		// If neither Body nor BodyS3Location is specified, Name is required
		if !hasBody && !hasBodyS3Location && !hasName {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': RestApi must have a Name property when Body and BodyS3Location are not specified",
					resName,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties"},
			})
		}
	}

	return matches
}
