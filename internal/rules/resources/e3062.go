// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3062{})
}

// E3062 validates RDS DB instance class compatibility with engine.
type E3062 struct{}

func (r *E3062) ID() string { return "E3062" }

func (r *E3062) ShortDesc() string {
	return "RDS instance class by engine"
}

func (r *E3062) Description() string {
	return "Validates that RDS DB instance class is compatible with the specified database engine and version."
}

func (r *E3062) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3062"
}

func (r *E3062) Tags() []string {
	return []string{"resources", "properties", "rds", "instance"}
}

// Simplified engine to instance class family compatibility
var engineInstanceCompatibility = map[string][]string{
	"mysql":      {"db.t2", "db.t3", "db.t4g", "db.m5", "db.m6g", "db.r5", "db.r6g"},
	"postgres":   {"db.t2", "db.t3", "db.t4g", "db.m5", "db.m6g", "db.r5", "db.r6g"},
	"mariadb":    {"db.t2", "db.t3", "db.m5", "db.r5"},
	"oracle-ee":  {"db.t3", "db.m5", "db.r5"},
	"oracle-se2": {"db.t3", "db.m5", "db.r5"},
	"sqlserver-ee": {"db.t3", "db.m5", "db.r5"},
	"sqlserver-se": {"db.t3", "db.m5", "db.r5"},
	"sqlserver-ex": {"db.t2", "db.t3"},
	"sqlserver-web": {"db.t3", "db.m5"},
	"aurora":        {"db.t3", "db.t4g", "db.r5", "db.r6g"},
	"aurora-mysql":  {"db.t3", "db.t4g", "db.r5", "db.r6g"},
	"aurora-postgresql": {"db.t3", "db.t4g", "db.r5", "db.r6g"},
}

func (r *E3062) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::RDS::DBInstance" {
			continue
		}

		engine, hasEngine := res.Properties["Engine"]
		instanceClass, hasInstanceClass := res.Properties["DBInstanceClass"]

		if !hasEngine || !hasInstanceClass {
			continue
		}

		engineStr, ok1 := engine.(string)
		instanceClassStr, ok2 := instanceClass.(string)

		if !ok1 || !ok2 {
			continue
		}

		// Normalize engine name
		engineStr = strings.ToLower(engineStr)

		// Get compatible instance families
		compatibleFamilies, engineKnown := engineInstanceCompatibility[engineStr]
		if !engineKnown {
			// Unknown engine, skip validation
			continue
		}

		// Extract instance family (e.g., db.t3 from db.t3.micro)
		instanceFamily := r.getInstanceFamily(instanceClassStr)

		// Check compatibility
		isCompatible := false
		for _, family := range compatibleFamilies {
			if strings.HasPrefix(instanceFamily, family) {
				isCompatible = true
				break
			}
		}

		if !isCompatible {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': DB instance class '%s' may not be compatible with engine '%s'. Compatible instance families: %v",
					resName, instanceClassStr, engineStr, compatibleFamilies,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties", "DBInstanceClass"},
			})
		}
	}

	return matches
}

// getInstanceFamily extracts the instance family from instance class
// e.g., "db.t3.micro" -> "db.t3"
func (r *E3062) getInstanceFamily(instanceClass string) string {
	parts := strings.Split(instanceClass, ".")
	if len(parts) >= 2 {
		return parts[0] + "." + parts[1]
	}
	return instanceClass
}
