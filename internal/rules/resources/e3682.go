// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3682{})
}

// E3682 validates Aurora property exclusions.
type E3682 struct{}

func (r *E3682) ID() string { return "E3682" }

func (r *E3682) ShortDesc() string {
	return "Validate Aurora property exclusions"
}

func (r *E3682) Description() string {
	return "Validates that Aurora DB instances do not specify certain properties that are not applicable."
}

func (r *E3682) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3682"
}

func (r *E3682) Tags() []string {
	return []string{"resources", "properties", "rds", "aurora"}
}

func (r *E3682) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::RDS::DBInstance" {
			continue
		}

		engine, hasEngine := res.Properties["Engine"]
		if !hasEngine || isIntrinsicFunction(engine) {
			continue
		}

		engineStr, ok := engine.(string)
		if !ok {
			continue
		}

		// Check if it's an Aurora engine
		if strings.HasPrefix(engineStr, "aurora") {
			// Aurora instances cannot have these properties
			invalidProps := []string{"BackupRetentionPeriod", "MasterUsername", "MasterUserPassword"}
			for _, prop := range invalidProps {
				if _, hasProp := res.Properties[prop]; hasProp {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': Aurora DB Instance must not specify %s (define at cluster level)",
							resName, prop,
						),
						Line:   res.Node.Line,
						Column: res.Node.Column,
						Path:   []string{"Resources", resName, "Properties", prop},
					})
				}
			}
		}
	}

	return matches
}
