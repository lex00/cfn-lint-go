// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3633{})
}

// E3633 validates Lambda event source mapping StartingPosition usage.
type E3633 struct{}

func (r *E3633) ID() string { return "E3633" }

func (r *E3633) ShortDesc() string {
	return "Lambda event source StartingPosition required"
}

func (r *E3633) Description() string {
	return "Validates that Lambda event source mapping has StartingPosition property when using Kinesis, Kafka, or DynamoDB stream event sources."
}

func (r *E3633) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3633"
}

func (r *E3633) Tags() []string {
	return []string{"resources", "properties", "lambda", "eventsourcemapping"}
}

func (r *E3633) Match(tmpl *template.Template) []rules.Match {
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

		// Check if the ARN is for Kinesis, Kafka, or DynamoDB stream
		requiresStartingPosition := false
		if containsStr(arnStr, ":kinesis:") ||
			containsStr(arnStr, ":kafka:") ||
			containsStr(arnStr, ":dynamodb:") {
			requiresStartingPosition = true
		}

		if requiresStartingPosition {
			if _, hasStartingPosition := res.Properties["StartingPosition"]; !hasStartingPosition {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': Lambda EventSourceMapping for Kinesis, Kafka, or DynamoDB stream requires StartingPosition property",
						resName,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	return matches
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
