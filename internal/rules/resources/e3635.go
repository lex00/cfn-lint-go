// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3635{})
}

// E3635 validates Neptune DB instance class.
type E3635 struct{}

func (r *E3635) ID() string { return "E3635" }

func (r *E3635) ShortDesc() string {
	return "Validate Neptune DB instance class"
}

func (r *E3635) Description() string {
	return "Validates that AWS::Neptune::DBInstance resources specify valid instance classes."
}

func (r *E3635) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3635"
}

func (r *E3635) Tags() []string {
	return []string{"resources", "properties", "neptune", "instanceclass"}
}

func (r *E3635) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Neptune::DBInstance" {
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

		// Neptune instance classes start with db.r4., db.r5., db.r6g., db.t3., db.t4g.
		if !strings.HasPrefix(instanceClassStr, "db.r4.") &&
			!strings.HasPrefix(instanceClassStr, "db.r5.") &&
			!strings.HasPrefix(instanceClassStr, "db.r6g.") &&
			!strings.HasPrefix(instanceClassStr, "db.t3.") &&
			!strings.HasPrefix(instanceClassStr, "db.t4g.") {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Invalid Neptune instance class '%s'. Must start with db.r4., db.r5., db.r6g., db.t3., or db.t4g.",
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
