// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3703{})
}

// E3703 validates CodePipeline action configuration.
type E3703 struct{}

func (r *E3703) ID() string { return "E3703" }

func (r *E3703) ShortDesc() string {
	return "Validate CodePipeline action configuration"
}

func (r *E3703) Description() string {
	return "Validates that AWS::CodePipeline::Pipeline actions have valid configuration based on action type."
}

func (r *E3703) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3703"
}

func (r *E3703) Tags() []string {
	return []string{"resources", "properties", "codepipeline", "configuration"}
}

func (r *E3703) Match(tmpl *template.Template) []rules.Match {
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

				// Validate that action has required fields
				_, hasName := actionMap["Name"]
				if !hasName {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': Pipeline action must have a Name",
							resName,
						),
						Line:   res.Node.Line,
						Column: res.Node.Column,
						Path:   []string{"Resources", resName, "Properties", "Stages"},
					})
				}

				_, hasActionTypeId := actionMap["ActionTypeId"]
				if !hasActionTypeId {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': Pipeline action must have an ActionTypeId",
							resName,
						),
						Line:   res.Node.Line,
						Column: res.Node.Column,
						Path:   []string{"Resources", resName, "Properties", "Stages"},
					})
					continue
				}

				actionTypeId, ok := actionMap["ActionTypeId"].(map[string]any)
				if !ok || isIntrinsicFunction(actionMap["ActionTypeId"]) {
					continue
				}

				// Validate ActionTypeId has required fields
				_, hasCategory := actionTypeId["Category"]
				if !hasCategory {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': Pipeline action ActionTypeId must have Category",
							resName,
						),
						Line:   res.Node.Line,
						Column: res.Node.Column,
						Path:   []string{"Resources", resName, "Properties", "Stages"},
					})
				}

				_, hasOwner := actionTypeId["Owner"]
				if !hasOwner {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': Pipeline action ActionTypeId must have Owner",
							resName,
						),
						Line:   res.Node.Line,
						Column: res.Node.Column,
						Path:   []string{"Resources", resName, "Properties", "Stages"},
					})
				}

				_, hasProvider := actionTypeId["Provider"]
				if !hasProvider {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': Pipeline action ActionTypeId must have Provider",
							resName,
						),
						Line:   res.Node.Line,
						Column: res.Node.Column,
						Path:   []string{"Resources", resName, "Properties", "Stages"},
					})
				}

				_, hasVersion := actionTypeId["Version"]
				if !hasVersion {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': Pipeline action ActionTypeId must have Version",
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

	return matches
}
