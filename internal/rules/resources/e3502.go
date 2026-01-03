// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3502{})
}

// E3502 validates SQS DLQ queue type matching.
type E3502 struct{}

func (r *E3502) ID() string { return "E3502" }

func (r *E3502) ShortDesc() string {
	return "SQS DLQ queue type match"
}

func (r *E3502) Description() string {
	return "Validates that SQS dead-letter queue destination matches the source queue type (FIFO/standard)."
}

func (r *E3502) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3502"
}

func (r *E3502) Tags() []string {
	return []string{"resources", "properties", "sqs", "queue", "dlq"}
}

func (r *E3502) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Track queue types
	queueTypes := make(map[string]bool) // true = FIFO, false = standard

	// First pass: determine queue types
	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::SQS::Queue" {
			continue
		}

		isFIFO := false
		if queueName, hasName := res.Properties["QueueName"]; hasName {
			if queueNameStr, ok := queueName.(string); ok {
				if len(queueNameStr) >= 5 && queueNameStr[len(queueNameStr)-5:] == ".fifo" {
					isFIFO = true
				}
			}
		}

		if fifoQueue, hasFifo := res.Properties["FifoQueue"]; hasFifo {
			if fifoQueueBool, ok := fifoQueue.(bool); ok && fifoQueueBool {
				isFIFO = true
			}
		}

		queueTypes[resName] = isFIFO
	}

	// Second pass: validate DLQ configurations
	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::SQS::Queue" {
			continue
		}

		redrivePolicy, hasRedrive := res.Properties["RedrivePolicy"]
		if !hasRedrive {
			continue
		}

		redrivePolicyMap, ok := redrivePolicy.(map[string]interface{})
		if !ok {
			continue
		}

		deadLetterTargetArn, hasDLQ := redrivePolicyMap["deadLetterTargetArn"]
		if !hasDLQ {
			continue
		}

		// Extract DLQ reference
		dlqRef := ""
		if dlqMap, ok := deadLetterTargetArn.(map[string]interface{}); ok {
			if getAtt, hasGetAtt := dlqMap["Fn::GetAtt"]; hasGetAtt {
				if getAttList, ok := getAtt.([]interface{}); ok && len(getAttList) > 0 {
					if refStr, ok := getAttList[0].(string); ok {
						dlqRef = refStr
					}
				}
			}
		}

		if dlqRef == "" {
			continue
		}

		// Check if both queues exist and their types match
		sourceFIFO, sourceExists := queueTypes[resName]
		targetFIFO, targetExists := queueTypes[dlqRef]

		if sourceExists && targetExists && sourceFIFO != targetFIFO {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Dead-letter queue '%s' type must match source queue type (%s queue cannot use %s DLQ)",
					resName, dlqRef, r.queueTypeStr(sourceFIFO), r.queueTypeStr(targetFIFO),
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties", "RedrivePolicy", "deadLetterTargetArn"},
			})
		}
	}

	return matches
}

func (r *E3502) queueTypeStr(isFIFO bool) string {
	if isFIFO {
		return "FIFO"
	}
	return "standard"
}
