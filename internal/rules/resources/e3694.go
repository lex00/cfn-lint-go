// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3694{})
}

// E3694 validates RDS DB Cluster instance class.
type E3694 struct{}

func (r *E3694) ID() string { return "E3694" }

func (r *E3694) ShortDesc() string {
	return "Validate RDS DB Cluster instance class"
}

func (r *E3694) Description() string {
	return "Validates that AWS::RDS::DBCluster resources with serverless v2 scaling specify valid instance classes."
}

func (r *E3694) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3694"
}

func (r *E3694) Tags() []string {
	return []string{"resources", "properties", "rds", "instanceclass"}
}

func (r *E3694) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// RDS instance class families
	validFamilies := []string{
		"db.t2.", "db.t3.", "db.t4g.",
		"db.m1.", "db.m2.", "db.m3.", "db.m4.", "db.m5.", "db.m5d.", "db.m6g.", "db.m6gd.", "db.m6i.", "db.m6id.",
		"db.r3.", "db.r4.", "db.r5.", "db.r5b.", "db.r5d.", "db.r6g.", "db.r6gd.", "db.r6i.", "db.r6id.",
		"db.x1.", "db.x1e.", "db.x2g.", "db.x2idn.", "db.x2iedn.", "db.x2iezn.",
		"db.z1d.",
	}

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::RDS::DBInstance" {
			continue
		}

		instanceClass, hasInstanceClass := res.Properties["DBInstanceClass"]
		if !hasInstanceClass || isIntrinsicFunction(instanceClass) {
			continue
		}

		instanceClassStr, ok := instanceClass.(string)
		if !ok {
			continue
		}

		// Check if instance class starts with a valid family
		isValid := false
		for _, family := range validFamilies {
			if strings.HasPrefix(instanceClassStr, family) {
				isValid = true
				break
			}
		}

		if !isValid {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Invalid RDS instance class '%s'. Must start with 'db.'",
					resName, instanceClassStr,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties", "DBInstanceClass"},
			})
		}
	}

	return matches
}
