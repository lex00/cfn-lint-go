package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W3101{})
}

// W3101 warns about SAM Function resources missing Timeout property.
type W3101 struct{}

func (r *W3101) ID() string { return "W3101" }

func (r *W3101) ShortDesc() string {
	return "SAM Function missing Timeout"
}

func (r *W3101) Description() string {
	return "Warns when an AWS::Serverless::Function resource does not specify Timeout. The default is 3 seconds which may be too short for many workloads. Consider explicitly setting Timeout based on your function's expected execution time."
}

func (r *W3101) Source() string {
	return "https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-function.html"
}

func (r *W3101) Tags() []string {
	return []string{"resources", "sam", "serverless", "lambda"}
}

func (r *W3101) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check if Globals has Timeout set
	hasGlobalTimeout := checkGlobalsFunctionProperty(tmpl, "Timeout")

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Serverless::Function" {
			continue
		}

		// Skip if function or Globals has Timeout
		if _, hasLocal := res.Properties["Timeout"]; hasLocal || hasGlobalTimeout {
			continue
		}

		line, column := 0, 0
		if res.Node != nil {
			line = res.Node.Line
			column = res.Node.Column
		}

		matches = append(matches, rules.Match{
			Message: fmt.Sprintf("SAM Function '%s' does not specify Timeout. Default is 3 seconds which may be too short. Consider setting an explicit value.", resName),
			Line:    line,
			Column:  column,
			Path:    []string{"Resources", resName, "Properties"},
		})
	}

	return matches
}
