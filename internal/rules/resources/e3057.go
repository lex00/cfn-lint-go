// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3057{})
}

// E3057 validates CloudFront TargetOriginId references.
type E3057 struct{}

func (r *E3057) ID() string { return "E3057" }

func (r *E3057) ShortDesc() string {
	return "CloudFront TargetOriginId"
}

func (r *E3057) Description() string {
	return "Validates that CloudFront distribution TargetOriginId references a defined Origin within the DistributionConfig."
}

func (r *E3057) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3057"
}

func (r *E3057) Tags() []string {
	return []string{"resources", "properties", "cloudfront"}
}

func (r *E3057) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::CloudFront::Distribution" {
			continue
		}

		distConfig, hasDistConfig := res.Properties["DistributionConfig"]
		if !hasDistConfig {
			continue
		}

		distConfigMap, ok := distConfig.(map[string]interface{})
		if !ok {
			continue
		}

		// Collect all Origin IDs
		originIDs := make(map[string]bool)
		if origins, hasOrigins := distConfigMap["Origins"]; hasOrigins {
			if originsList, ok := origins.([]interface{}); ok {
				for _, origin := range originsList {
					if originMap, ok := origin.(map[string]interface{}); ok {
						if id, hasID := originMap["Id"]; hasID {
							if idStr, ok := id.(string); ok {
								originIDs[idStr] = true
							}
						}
					}
				}
			}
		}

		// Check DefaultCacheBehavior TargetOriginId
		if defaultBehavior, hasDefault := distConfigMap["DefaultCacheBehavior"]; hasDefault {
			if behaviorMap, ok := defaultBehavior.(map[string]interface{}); ok {
				if targetOriginId, hasTarget := behaviorMap["TargetOriginId"]; hasTarget {
					if targetStr, ok := targetOriginId.(string); ok {
						if !originIDs[targetStr] {
							matches = append(matches, rules.Match{
								Message: fmt.Sprintf(
									"Resource '%s': DefaultCacheBehavior TargetOriginId '%s' does not reference a defined Origin",
									resName, targetStr,
								),
								Line:   res.Node.Line,
								Column: res.Node.Column,
								Path:   []string{"Resources", resName, "Properties", "DistributionConfig", "DefaultCacheBehavior", "TargetOriginId"},
							})
						}
					}
				}
			}
		}

		// Check CacheBehaviors TargetOriginId
		if cacheBehaviors, hasBehaviors := distConfigMap["CacheBehaviors"]; hasBehaviors {
			if behaviorsList, ok := cacheBehaviors.([]interface{}); ok {
				for i, behavior := range behaviorsList {
					if behaviorMap, ok := behavior.(map[string]interface{}); ok {
						if targetOriginId, hasTarget := behaviorMap["TargetOriginId"]; hasTarget {
							if targetStr, ok := targetOriginId.(string); ok {
								if !originIDs[targetStr] {
									matches = append(matches, rules.Match{
										Message: fmt.Sprintf(
											"Resource '%s': CacheBehavior %d TargetOriginId '%s' does not reference a defined Origin",
											resName, i, targetStr,
										),
										Line:   res.Node.Line,
										Column: res.Node.Column,
										Path:   []string{"Resources", resName, "Properties", "DistributionConfig", "CacheBehaviors", fmt.Sprintf("[%d]", i), "TargetOriginId"},
									})
								}
							}
						}
					}
				}
			}
		}
	}

	return matches
}
