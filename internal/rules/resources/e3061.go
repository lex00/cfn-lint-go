// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3061{})
}

// E3061 validates S3 IntelligentTieringConfigurations days.
type E3061 struct{}

func (r *E3061) ID() string { return "E3061" }

func (r *E3061) ShortDesc() string {
	return "IntelligentTieringConfigurations days"
}

func (r *E3061) Description() string {
	return "Validates minimum and maximum day values for S3 Intelligent-Tiering configurations."
}

func (r *E3061) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3061"
}

func (r *E3061) Tags() []string {
	return []string{"resources", "properties", "s3", "intelligenttiering"}
}

func (r *E3061) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::S3::Bucket" {
			continue
		}

		intelligentTiering, hasTiering := res.Properties["IntelligentTieringConfigurations"]
		if !hasTiering {
			continue
		}

		tieringList, ok := intelligentTiering.([]interface{})
		if !ok {
			continue
		}

		for i, config := range tieringList {
			configMap, ok := config.(map[string]interface{})
			if !ok {
				continue
			}

			// Check Tierings
			tierings, hasTierings := configMap["Tierings"]
			if !hasTierings {
				continue
			}

			tieringsList, ok := tierings.([]interface{})
			if !ok {
				continue
			}

			for j, tiering := range tieringsList {
				tieringMap, ok := tiering.(map[string]interface{})
				if !ok {
					continue
				}

				days, hasDays := tieringMap["Days"]
				accessTier, hasAccessTier := tieringMap["AccessTier"]

				if !hasDays {
					continue
				}

				daysInt, ok := r.toInt(days)
				if !ok {
					continue
				}

				// Validate days based on access tier
				if hasAccessTier {
					accessTierStr, ok := accessTier.(string)
					if ok {
						switch accessTierStr {
						case "ARCHIVE_ACCESS":
							if daysInt < 90 {
								matches = append(matches, rules.Match{
									Message: fmt.Sprintf(
										"Resource '%s': IntelligentTieringConfiguration %d Tiering %d with AccessTier ARCHIVE_ACCESS must have Days >= 90 (got %d)",
										resName, i, j, daysInt,
									),
									Line:   res.Node.Line,
									Column: res.Node.Column,
									Path:   []string{"Resources", resName, "Properties", "IntelligentTieringConfigurations", fmt.Sprintf("[%d]", i), "Tierings", fmt.Sprintf("[%d]", j), "Days"},
								})
							}
						case "DEEP_ARCHIVE_ACCESS":
							if daysInt < 180 {
								matches = append(matches, rules.Match{
									Message: fmt.Sprintf(
										"Resource '%s': IntelligentTieringConfiguration %d Tiering %d with AccessTier DEEP_ARCHIVE_ACCESS must have Days >= 180 (got %d)",
										resName, i, j, daysInt,
									),
									Line:   res.Node.Line,
									Column: res.Node.Column,
									Path:   []string{"Resources", resName, "Properties", "IntelligentTieringConfigurations", fmt.Sprintf("[%d]", i), "Tierings", fmt.Sprintf("[%d]", j), "Days"},
								})
							}
						}
					}
				}

				// General minimum check
				if daysInt < 1 {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': IntelligentTieringConfiguration %d Tiering %d Days must be at least 1 (got %d)",
							resName, i, j, daysInt,
						),
						Line:   res.Node.Line,
						Column: res.Node.Column,
						Path:   []string{"Resources", resName, "Properties", "IntelligentTieringConfigurations", fmt.Sprintf("[%d]", i), "Tierings", fmt.Sprintf("[%d]", j), "Days"},
					})
				}
			}
		}
	}

	return matches
}

func (r *E3061) toInt(value interface{}) (int, bool) {
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
