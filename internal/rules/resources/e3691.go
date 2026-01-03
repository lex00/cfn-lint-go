// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3691{})
}

// E3691 validates DB Instance engine and version.
type E3691 struct{}

func (r *E3691) ID() string { return "E3691" }

func (r *E3691) ShortDesc() string {
	return "Validate DB Instance engine and version"
}

func (r *E3691) Description() string {
	return "Validates that AWS::RDS::DBInstance resources specify valid engine types."
}

func (r *E3691) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3691"
}

func (r *E3691) Tags() []string {
	return []string{"resources", "properties", "rds", "engine"}
}

func (r *E3691) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	validEngines := map[string]bool{
		"aurora":            true,
		"aurora-mysql":      true,
		"aurora-postgresql": true,
		"mariadb":           true,
		"mysql":             true,
		"oracle-ee":         true,
		"oracle-ee-cdb":     true,
		"oracle-se2":        true,
		"oracle-se2-cdb":    true,
		"postgres":          true,
		"sqlserver-ee":      true,
		"sqlserver-se":      true,
		"sqlserver-ex":      true,
		"sqlserver-web":     true,
		"custom-oracle-ee":  true,
		"custom-oracle-ee-cdb": true,
		"custom-sqlserver-ee": true,
		"custom-sqlserver-se": true,
		"custom-sqlserver-web": true,
	}

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

		if !validEngines[engineStr] {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Invalid DB Instance engine '%s'",
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
