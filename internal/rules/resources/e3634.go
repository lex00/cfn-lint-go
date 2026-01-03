// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3634{})
}

// E3634 validates Lambda event source mapping doesn't use StartingPosition with SQS.
type E3634 struct{}

func (r *E3634) ID() string { return "E3634" }

func (r *E3634) ShortDesc() string {
	return "Lambda SQS event StartingPosition not allowed"
}

func (r *E3634) Description() string {
	return "Validates that Lambda event source mapping does not specify StartingPosition property when using SQS event sources."
}

func (r *E3634) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3634"
}

func (r *E3634) Tags() []string {
	return []string{"resources", "properties", "lambda", "eventsourcemapping", "sqs"}
}

func (r *E3634) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Lambda::EventSourceMapping" {
			continue
		}

		eventSourceArn, hasArn := res.Properties["EventSourceArn"]
		if !hasArn || isIntrinsicFunction(eventSourceArn) {
			continue
		}

		arnStr, ok := eventSourceArn.(string)
		if !ok {
			continue
		}

		// Check if the ARN is for SQS
		if containsStr(arnStr, ":sqs:") {
			if _, hasStartingPosition := res.Properties["StartingPosition"]; hasStartingPosition {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': Lambda EventSourceMapping for SQS must not specify StartingPosition property",
						resName,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", "StartingPosition"},
				})
			}
		}
	}

	return matches
}
