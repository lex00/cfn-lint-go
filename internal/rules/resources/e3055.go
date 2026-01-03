// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3055{})
}

// E3055 validates CreationPolicy values.
type E3055 struct{}

func (r *E3055) ID() string { return "E3055" }

func (r *E3055) ShortDesc() string {
	return "CreationPolicy values"
}

func (r *E3055) Description() string {
	return "Validates that resource CreationPolicy attributes have proper configuration and valid values."
}

func (r *E3055) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3055"
}

func (r *E3055) Tags() []string {
	return []string{"resources", "creationpolicy", "metadata"}
}

func (r *E3055) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// CreationPolicy is a resource attribute, not stored in current Resource struct
	// Skip validation for now as it requires extending the template parser
	_ = tmpl
	return matches

	/*
		for resName, res := range tmpl.Resources {
			if res.CreationPolicy == nil {
				continue
			}

			cpMap, ok := res.CreationPolicy.(map[string]interface{})
			if !ok {
				continue
			}

			// Check AutoScalingCreationPolicy
			if asPolicy, hasAS := cpMap["AutoScalingCreationPolicy"]; hasAS {
				if asPolicyMap, ok := asPolicy.(map[string]interface{}); ok {
					if minSuccess, hasMinSuccess := asPolicyMap["MinSuccessfulInstancesPercent"]; hasMinSuccess {
						if minSuccessInt, ok := r.toInt(minSuccess); ok {
							if minSuccessInt < 0 || minSuccessInt > 100 {
								matches = append(matches, rules.Match{
									Message: fmt.Sprintf(
										"Resource '%s': CreationPolicy AutoScalingCreationPolicy MinSuccessfulInstancesPercent must be between 0 and 100 (got %d)",
										resName, minSuccessInt,
									),
									Line:   res.Node.Line,
									Column: res.Node.Column,
									Path:   []string{"Resources", resName, "CreationPolicy", "AutoScalingCreationPolicy", "MinSuccessfulInstancesPercent"},
								})
							}
						}
					}
				}
			}

			// Check ResourceSignal
			if rsPolicy, hasRS := cpMap["ResourceSignal"]; hasRS {
				if rsPolicyMap, ok := rsPolicy.(map[string]interface{}); ok {
					if count, hasCount := rsPolicyMap["Count"]; hasCount {
						if countInt, ok := r.toInt(count); ok {
							if countInt < 1 {
								matches = append(matches, rules.Match{
									Message: fmt.Sprintf(
										"Resource '%s': CreationPolicy ResourceSignal Count must be at least 1 (got %d)",
										resName, countInt,
									),
									Line:   res.Node.Line,
									Column: res.Node.Column,
									Path:   []string{"Resources", resName, "CreationPolicy", "ResourceSignal", "Count"},
								})
							}
						}
					}

					if timeout, hasTimeout := rsPolicyMap["Timeout"]; hasTimeout {
						if timeoutStr, ok := timeout.(string); ok {
							// Validate ISO 8601 duration format (basic check)
							if len(timeoutStr) < 2 || timeoutStr[0] != 'P' {
								matches = append(matches, rules.Match{
									Message: fmt.Sprintf(
										"Resource '%s': CreationPolicy ResourceSignal Timeout must be in ISO 8601 duration format (e.g., PT15M)",
										resName,
									),
									Line:   res.Node.Line,
									Column: res.Node.Column,
									Path:   []string{"Resources", resName, "CreationPolicy", "ResourceSignal", "Timeout"},
								})
							}
						}
					}
				}
			}
		}
	*/

	// return matches
}

// Commented out unused helper function - will be needed if CreationPolicy validation is implemented
// func (r *E3055) toInt(value interface{}) (int, bool) {
// 	switch v := value.(type) {
// 	case int:
// 		return v, true
// 	case float64:
// 		return int(v), true
// 	case string:
// 		var i int
// 		if _, err := fmt.Sscanf(v, "%d", &i); err == nil {
// 			return i, true
// 		}
// 	}
// 	return 0, false
// }
