// Package formats contains format validation rules (E11xx).
package formats

import (
	"fmt"
	"regexp"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1155{})
}

// E1155 validates CloudWatch Logs group name format.
type E1155 struct{}

func (r *E1155) ID() string { return "E1155" }

func (r *E1155) ShortDesc() string {
	return "CloudWatch logs group name format validation"
}

func (r *E1155) Description() string {
	return "Validates that CloudWatch Logs log group names match valid format."
}

func (r *E1155) Source() string {
	return "https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/Working-with-log-groups-and-streams.html"
}

func (r *E1155) Tags() []string {
	return []string{"format", "cloudwatch", "logs"}
}

// CloudWatch Logs log group name pattern: alphanumeric, hyphens, underscores, periods, forward slashes
// Length: 1-512 characters
var logGroupNamePattern = regexp.MustCompile(`^[\w\-\./]+$`)

func (r *E1155) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		// Check AWS::Logs::LogGroup resources
		if res.Type == "AWS::Logs::LogGroup" {
			if logGroupName := getStringProperty(res.Properties, "LogGroupName"); logGroupName != "" {
				if !logGroupNamePattern.MatchString(logGroupName) {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Invalid CloudWatch Logs log group name '%s' in resource '%s', must contain only alphanumeric, hyphens, underscores, periods, and forward slashes", logGroupName, resName),
						Path:    []string{"Resources", resName, "Properties", "LogGroupName"},
					})
				}
				if len(logGroupName) > 512 {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("CloudWatch Logs log group name '%s' in resource '%s' exceeds maximum length of 512 characters", logGroupName, resName),
						Path:    []string{"Resources", resName, "Properties", "LogGroupName"},
					})
				}
			}
		}
	}

	return matches
}
