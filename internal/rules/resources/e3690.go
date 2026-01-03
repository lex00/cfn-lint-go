// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3690{})
}

// E3690 validates DB Cluster engine and version.
type E3690 struct{}

func (r *E3690) ID() string { return "E3690" }

func (r *E3690) ShortDesc() string {
	return "Validate DB Cluster engine and version"
}

func (r *E3690) Description() string {
	return "Validates that AWS::RDS::DBCluster resources specify valid engine types."
}

func (r *E3690) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3690"
}

func (r *E3690) Tags() []string {
	return []string{"resources", "properties", "rds", "engine"}
}

func (r *E3690) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	validEngines := map[string]bool{
		"aurora":            true,
		"aurora-mysql":      true,
		"aurora-postgresql": true,
		"mysql":             true,
		"postgres":          true,
	}

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

		if !validEngines[engineStr] {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Invalid DB Cluster engine '%s'. Must be aurora, aurora-mysql, aurora-postgresql, mysql, or postgres",
					resName, engineStr,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties", "Engine"},
			})
		}
	}

	return matches
}
