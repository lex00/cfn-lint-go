// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3700{})
}

// E3700 validates CodePipeline source actions in first stage only.
type E3700 struct{}

func (r *E3700) ID() string { return "E3700" }

func (r *E3700) ShortDesc() string {
	return "CodePipeline Source actions must be in first stage"
}

func (r *E3700) Description() string {
	return "Validates that AWS::CodePipeline::Pipeline Source actions are only in the first stage."
}

func (r *E3700) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3700"
}

func (r *E3700) Tags() []string {
	return []string{"resources", "properties", "codepipeline", "stages"}
}

func (r *E3700) Match(tmpl *template.Template) []rules.Match {
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

		for stageIdx, stage := range stageList {
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

				actionCategory, hasCategory := actionMap["ActionTypeId"]
				if !hasCategory || isIntrinsicFunction(actionCategory) {
					continue
				}

				categoryMap, ok := actionCategory.(map[string]any)
				if !ok {
					continue
				}

				category, hasCategory := categoryMap["Category"]
				if !hasCategory || isIntrinsicFunction(category) {
					continue
				}

				categoryStr, ok := category.(string)
				if !ok {
					continue
				}

				// Source actions must be in the first stage (index 0)
				if categoryStr == "Source" && stageIdx != 0 {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': Source actions must be in the first stage only. Found in stage %d",
							resName, stageIdx,
						),
						Line:   res.Node.Line,
						Column: res.Node.Column,
						Path:   []string{"Resources", resName, "Properties", "Stages"},
					})
				}
			}
		}
	}

	return matches
}
