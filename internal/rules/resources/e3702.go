// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3702{})
}

// E3702 validates CodePipeline artifact count validation.
type E3702 struct{}

func (r *E3702) ID() string { return "E3702" }

func (r *E3702) ShortDesc() string {
	return "Validate CodePipeline artifact count"
}

func (r *E3702) Description() string {
	return "Validates that AWS::CodePipeline::Pipeline actions have appropriate artifact counts based on action type."
}

func (r *E3702) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3702"
}

func (r *E3702) Tags() []string {
	return []string{"resources", "properties", "codepipeline", "artifacts"}
}

func (r *E3702) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::CodePipeline::Pipeline" {
			continue
		}

		stages, hasStages := res.Properties["Stages"]
		if !hasStages || isIntrinsicFunction(stages) {
			continue
		}

		stageList, ok := stages.([]any)
		if !ok {
			continue
		}

		for _, stage := range stageList {
			stageMap, ok := stage.(map[string]any)
			if !ok {
				continue
			}

			actions, hasActions := stageMap["Actions"]
			if !hasActions || isIntrinsicFunction(actions) {
				continue
			}

			actionList, ok := actions.([]any)
			if !ok {
				continue
			}

			for _, action := range actionList {
				actionMap, ok := action.(map[string]any)
				if !ok {
					continue
				}

				actionTypeId, hasActionTypeId := actionMap["ActionTypeId"]
				if !hasActionTypeId || isIntrinsicFunction(actionTypeId) {
					continue
				}

				actionTypeMap, ok := actionTypeId.(map[string]any)
				if !ok {
					continue
				}

				category, hasCategory := actionTypeMap["Category"]
				if !hasCategory || isIntrinsicFunction(category) {
					continue
				}

				categoryStr, ok := category.(string)
				if !ok {
					continue
				}

				// Source actions must have at least one output artifact
				if categoryStr == "Source" {
					outputArtifacts, hasOutput := actionMap["OutputArtifacts"]
					if !hasOutput || isIntrinsicFunction(outputArtifacts) {
						continue
					}

					outputList, ok := outputArtifacts.([]any)
					if !ok {
						continue
					}

					if len(outputList) == 0 {
						matches = append(matches, rules.Match{
							Message: fmt.Sprintf(
								"Resource '%s': Source action must have at least one OutputArtifact",
								resName,
							),
							Line:   res.Node.Line,
							Column: res.Node.Column,
							Path:   []string{"Resources", resName, "Properties", "Stages"},
						})
					}
				}

				// Deploy actions cannot have output artifacts (in most cases)
				if categoryStr == "Deploy" {
					outputArtifacts, hasOutput := actionMap["OutputArtifacts"]
					if hasOutput && !isIntrinsicFunction(outputArtifacts) {
						outputList, ok := outputArtifacts.([]any)
						if ok && len(outputList) > 0 {
							matches = append(matches, rules.Match{
								Message: fmt.Sprintf(
									"Resource '%s': Deploy action typically should not have OutputArtifacts",
									resName,
								),
								Line:   res.Node.Line,
								Column: res.Node.Column,
								Path:   []string{"Resources", resName, "Properties", "Stages"},
							})
						}
					}
				}
			}
		}
	}

	return matches
}
