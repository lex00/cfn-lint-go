package resources

import (
	"fmt"
	"regexp"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3027{})
}

// E3027 validates AWS Event ScheduleExpression format.
type E3027 struct{}

func (r *E3027) ID() string {
	return "E3027"
}

func (r *E3027) ShortDesc() string {
	return "Validate AWS Event ScheduleExpression format"
}

func (r *E3027) Description() string {
	return "Confirms proper syntax for CloudWatch event scheduling expressions"
}

func (r *E3027) Source() string {
	return "https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-create-rule-schedule.html"
}

func (r *E3027) Tags() []string {
	return []string{"resources", "events", "schedule"}
}

// Basic regex patterns for schedule expressions
var (
	rateExpressionRegex = regexp.MustCompile(`^rate\(\d+\s+(minute|minutes|hour|hours|day|days)\)$`)
	cronExpressionRegex = regexp.MustCompile(`^cron\([^\)]+\)$`)
)

func (r *E3027) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Events::Rule" {
			continue
		}

		scheduleExprRaw, hasSchedule := res.Properties["ScheduleExpression"]
		if !hasSchedule {
			continue
		}

		// Skip intrinsic functions
		if isIntrinsicFunction(scheduleExprRaw) {
			continue
		}

		scheduleExpr, ok := scheduleExprRaw.(string)
		if !ok {
			continue
		}

		// Validate schedule expression format
		if !rateExpressionRegex.MatchString(scheduleExpr) && !cronExpressionRegex.MatchString(scheduleExpr) {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Events Rule '%s' has invalid ScheduleExpression '%s' (must be rate(...) or cron(...))", resName, scheduleExpr),
				Line:    res.Node.Line,
				Column:  res.Node.Column,
				Path:    []string{"Resources", resName, "Properties", "ScheduleExpression"},
			})
		}
	}

	return matches
}
