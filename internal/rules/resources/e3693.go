// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3693{})
}

// E3693 validates Aurora DB cluster config.
type E3693 struct{}

func (r *E3693) ID() string { return "E3693" }

func (r *E3693) ShortDesc() string {
	return "Validate Aurora DB cluster config"
}

func (r *E3693) Description() string {
	return "Validates that Aurora DB clusters are configured correctly and do not specify incompatible properties."
}

func (r *E3693) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3693"
}

func (r *E3693) Tags() []string {
	return []string{"resources", "properties", "rds", "aurora"}
}

func (r *E3693) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::RDS::DBCluster" {
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
			// Aurora clusters should not specify AllocatedStorage (unless Multi-AZ)
			_, hasInstanceClass := res.Properties["DBClusterInstanceClass"]
			if !hasInstanceClass {
				_, hasAllocatedStorage := res.Properties["AllocatedStorage"]
				if hasAllocatedStorage {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': Aurora DB Cluster should not specify AllocatedStorage unless using Multi-AZ (DBClusterInstanceClass)",
							resName,
						),
						Line:   res.Node.Line,
						Column: res.Node.Column,
						Path:   []string{"Resources", resName, "Properties", "AllocatedStorage"},
					})
				}
			}
		}
	}

	return matches
}
