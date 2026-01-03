// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3505{})
}

// E3505 validates SQS VisibilityTimeout vs Lambda Timeout.
type E3505 struct{}

func (r *E3505) ID() string { return "E3505" }

func (r *E3505) ShortDesc() string {
	return "SQS VisibilityTimeout vs Lambda"
}

func (r *E3505) Description() string {
	return "Validates that SQS queue VisibilityTimeout is greater than or equal to the Lambda function timeout when used as an event source."
}

func (r *E3505) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3505"
}

func (r *E3505) Tags() []string {
	return []string{"resources", "properties", "sqs", "lambda", "eventsource"}
}

func (r *E3505) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Track Lambda function timeouts
	lambdaTimeouts := make(map[string]int)

	for resName, res := range tmpl.Resources {
		if res.Type == "AWS::Lambda::Function" {
			timeout := 3 // Default Lambda timeout is 3 seconds
			if timeoutVal, hasTimeout := res.Properties["Timeout"]; hasTimeout {
				if timeoutInt, ok := r.toInt(timeoutVal); ok {
					timeout = timeoutInt
				}
			}
			lambdaTimeouts[resName] = timeout
		}
	}

	// Track SQS queue visibility timeouts
	queueVisibility := make(map[string]int)

	for resName, res := range tmpl.Resources {
		if res.Type == "AWS::SQS::Queue" {
			visibility := 30 // Default SQS visibility timeout is 30 seconds
			if visibilityVal, hasVisibility := res.Properties["VisibilityTimeout"]; hasVisibility {
				if visibilityInt, ok := r.toInt(visibilityVal); ok {
					visibility = visibilityInt
				}
			}
			queueVisibility[resName] = visibility
		}
	}

	// Check Lambda event source mappings
	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Lambda::EventSourceMapping" {
			continue
		}

		eventSourceArn, hasEventSource := res.Properties["EventSourceArn"]
		functionName, hasFunction := res.Properties["FunctionName"]

		if !hasEventSource || !hasFunction {
			continue
		}

		// Extract queue reference from EventSourceArn
		queueRef := ""
		if arnMap, ok := eventSourceArn.(map[string]interface{}); ok {
			if getAtt, hasGetAtt := arnMap["Fn::GetAtt"]; hasGetAtt {
				if getAttList, ok := getAtt.([]interface{}); ok && len(getAttList) > 0 {
					if refStr, ok := getAttList[0].(string); ok {
						queueRef = refStr
					}
				}
			}
		}

		// Extract function reference
		funcRef := ""
		if funcMap, ok := functionName.(map[string]interface{}); ok {
			if ref, hasRef := funcMap["Ref"]; hasRef {
				if refStr, ok := ref.(string); ok {
					funcRef = refStr
				}
			} else if getAtt, hasGetAtt := funcMap["Fn::GetAtt"]; hasGetAtt {
				if getAttList, ok := getAtt.([]interface{}); ok && len(getAttList) > 0 {
					if refStr, ok := getAttList[0].(string); ok {
						funcRef = refStr
					}
				}
			}
		} else if funcStr, ok := functionName.(string); ok {
			funcRef = funcStr
		}

		if queueRef == "" || funcRef == "" {
			continue
		}

		// Compare timeouts
		queueVis, hasQueue := queueVisibility[queueRef]
		lambdaTimeout, hasLambda := lambdaTimeouts[funcRef]

		if hasQueue && hasLambda && queueVis < lambdaTimeout {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': SQS queue '%s' VisibilityTimeout (%d seconds) should be >= Lambda function '%s' Timeout (%d seconds)",
					resName, queueRef, queueVis, funcRef, lambdaTimeout,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName},
			})
		}
	}

	return matches
}

func (r *E3505) toInt(value interface{}) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case float64:
		return int(v), true
	case string:
		var i int
		if _, err := fmt.Sscanf(v, "%d", &i); err == nil {
			return i, true
		}
	}
	return 0, false
}
