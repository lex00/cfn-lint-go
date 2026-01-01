// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"regexp"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3006{})
}

// E3006 checks that resource types follow the AWS naming convention.
type E3006 struct{}

func (r *E3006) ID() string { return "E3006" }

func (r *E3006) ShortDesc() string {
	return "Invalid resource type format"
}

func (r *E3006) Description() string {
	return "Checks that resource Type follows the format 'AWS::Service::Resource' or 'Custom::*'."
}

func (r *E3006) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3006"
}

func (r *E3006) Tags() []string {
	return []string{"resources", "type"}
}

// resourceTypePattern matches valid CloudFormation resource types:
// - AWS::Service::Resource (e.g., AWS::S3::Bucket)
// - AWS::Serverless::* (SAM resources)
// - Alexa::ASK::* (Alexa skill resources)
var resourceTypePattern = regexp.MustCompile(`^(AWS|Alexa)::[A-Za-z0-9]+::[A-Za-z0-9]+$`)

// customResourcePattern matches custom resources (Custom::Name)
var customResourcePattern = regexp.MustCompile(`^Custom::[A-Za-z0-9]+$`)

// modulePattern matches module-based resource types
var modulePattern = regexp.MustCompile(`^[A-Za-z0-9]+::[A-Za-z0-9]+::[A-Za-z0-9]+::MODULE$`)

func (r *E3006) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for name, res := range tmpl.Resources {
		if res.Type == "" {
			// E3001 handles missing Type
			continue
		}

		// Check if it matches valid patterns
		if !isValidResourceType(res.Type) {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Resource '%s' has invalid type '%s'. Expected format: 'AWS::Service::Resource' or 'Custom::Name'", name, res.Type),
				Line:    res.Node.Line,
				Column:  res.Node.Column,
				Path:    []string{"Resources", name, "Type"},
			})
		}
	}

	return matches
}

func isValidResourceType(t string) bool {
	// Check standard patterns (AWS::*, Alexa::*)
	if resourceTypePattern.MatchString(t) {
		return true
	}
	// Check custom resource pattern (Custom::Name)
	if customResourcePattern.MatchString(t) {
		return true
	}
	// Check module pattern
	if modulePattern.MatchString(t) {
		return true
	}
	return false
}
