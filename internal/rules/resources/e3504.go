// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3504{})
}

// E3504 validates AWS Backup plan lifecycle rules.
type E3504 struct{}

func (r *E3504) ID() string { return "E3504" }

func (r *E3504) ShortDesc() string {
	return "BackupPlan cold/delete timing"
}

func (r *E3504) Description() string {
	return "Validates that AWS Backup plan enforces a minimum 90-day gap between moving to cold storage and deletion."
}

func (r *E3504) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3504"
}

func (r *E3504) Tags() []string {
	return []string{"resources", "properties", "backup", "lifecycle"}
}

func (r *E3504) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Backup::BackupPlan" {
			continue
		}

		backupPlan, hasPlan := res.Properties["BackupPlan"]
		if !hasPlan {
			continue
		}

		planMap, ok := backupPlan.(map[string]interface{})
		if !ok {
			continue
		}

		backupPlanRules, hasRules := planMap["BackupPlanRule"]
		if !hasRules {
			continue
		}

		rulesList, ok := backupPlanRules.([]interface{})
		if !ok {
			continue
		}

		for i, rule := range rulesList {
			ruleMap, ok := rule.(map[string]interface{})
			if !ok {
				continue
			}

			lifecycle, hasLifecycle := ruleMap["Lifecycle"]
			if !hasLifecycle {
				continue
			}

			lifecycleMap, ok := lifecycle.(map[string]interface{})
			if !ok {
				continue
			}

			moveToColdStorage, hasCold := lifecycleMap["MoveToColdStorageAfterDays"]
			deleteAfter, hasDelete := lifecycleMap["DeleteAfterDays"]

			if !hasCold || !hasDelete {
				continue
			}

			coldDays, ok1 := r.toInt(moveToColdStorage)
			deleteDays, ok2 := r.toInt(deleteAfter)

			if !ok1 || !ok2 {
				continue
			}

			// Validate 90-day minimum gap
			gap := deleteDays - coldDays
			if gap < 90 {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': BackupPlanRule %d must have at least 90 days between MoveToColdStorageAfterDays (%d) and DeleteAfterDays (%d), gap is %d days",
						resName, i, coldDays, deleteDays, gap,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", "BackupPlan", "BackupPlanRule", fmt.Sprintf("[%d]", i), "Lifecycle"},
				})
			}
		}
	}

	return matches
}

func (r *E3504) toInt(value interface{}) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case float64:
		return int(v), true
	case string:
		var i int
		if _, err := fmt.Sscanf(v, "%d", &i); err == nil {
			return i, true
		}
	}
	return 0, false
}
