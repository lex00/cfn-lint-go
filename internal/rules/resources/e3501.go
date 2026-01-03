// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3501{})
}

// E3501 validates SQS queue properties based on queue type.
type E3501 struct{}

func (r *E3501) ID() string { return "E3501" }

func (r *E3501) ShortDesc() string {
	return "SQS queue properties"
}

func (r *E3501) Description() string {
	return "Validates that SQS queue properties are appropriate for the queue type (FIFO vs standard)."
}

func (r *E3501) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3501"
}

func (r *E3501) Tags() []string {
	return []string{"resources", "properties", "sqs", "queue"}
}

func (r *E3501) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::SQS::Queue" {
			continue
		}

		// Determine if this is a FIFO queue
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

		// FIFO-specific validations
		if isFIFO {
			// ContentBasedDeduplication is only valid for FIFO queues
			// (this is allowed)

			// Check for properties that are NOT allowed for FIFO queues
			// (Note: Most properties are allowed for both types)
		} else {
			// Standard queue validations
			// ContentBasedDeduplication should not be set for standard queues
			if _, hasDedup := res.Properties["ContentBasedDeduplication"]; hasDedup {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': ContentBasedDeduplication is only valid for FIFO queues",
						resName,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", "ContentBasedDeduplication"},
				})
			}

			// DeduplicationScope should not be set for standard queues
			if _, hasScope := res.Properties["DeduplicationScope"]; hasScope {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': DeduplicationScope is only valid for FIFO queues",
						resName,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", "DeduplicationScope"},
				})
			}

			// FifoThroughputLimit should not be set for standard queues
			if _, hasThroughput := res.Properties["FifoThroughputLimit"]; hasThroughput {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': FifoThroughputLimit is only valid for FIFO queues",
						resName,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", "FifoThroughputLimit"},
				})
			}
		}

		// Validate MessageRetentionPeriod (both types)
		if retention, hasRetention := res.Properties["MessageRetentionPeriod"]; hasRetention {
			if retentionInt, ok := r.toInt(retention); ok {
				if retentionInt < 60 || retentionInt > 1209600 {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': MessageRetentionPeriod must be between 60 and 1209600 seconds (got %d)",
							resName, retentionInt,
						),
						Line:   res.Node.Line,
						Column: res.Node.Column,
						Path:   []string{"Resources", resName, "Properties", "MessageRetentionPeriod"},
					})
				}
			}
		}

		// Validate VisibilityTimeout (both types)
		if visibility, hasVisibility := res.Properties["VisibilityTimeout"]; hasVisibility {
			if visibilityInt, ok := r.toInt(visibility); ok {
				if visibilityInt < 0 || visibilityInt > 43200 {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': VisibilityTimeout must be between 0 and 43200 seconds (got %d)",
							resName, visibilityInt,
						),
						Line:   res.Node.Line,
						Column: res.Node.Column,
						Path:   []string{"Resources", resName, "Properties", "VisibilityTimeout"},
					})
				}
			}
		}
	}

	return matches
}

func (r *E3501) toInt(value interface{}) (int, bool) {
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
