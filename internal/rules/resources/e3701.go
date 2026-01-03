// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3701{})
}

// E3701 validates CodePipeline artifact name usage.
type E3701 struct{}

func (r *E3701) ID() string { return "E3701" }

func (r *E3701) ShortDesc() string {
	return "CodePipeline InputArtifacts must reference OutputArtifacts"
}

func (r *E3701) Description() string {
	return "Validates that AWS::CodePipeline::Pipeline InputArtifacts reference OutputArtifacts from previous actions."
}

func (r *E3701) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3701"
}

func (r *E3701) Tags() []string {
	return []string{"resources", "properties", "codepipeline", "artifacts"}
}

func (r *E3701) Match(tmpl *template.Template) []rules.Match {
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

		// Collect all output artifacts
		outputArtifacts := make(map[string]bool)

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

				// Collect output artifacts
				if outputArtifactsList, hasOutput := actionMap["OutputArtifacts"]; hasOutput && !isIntrinsicFunction(outputArtifactsList) {
					if outputList, ok := outputArtifactsList.([]any); ok {
						for _, output := range outputList {
							if outputMap, ok := output.(map[string]any); ok {
								if name, hasName := outputMap["Name"]; hasName {
									if nameStr, ok := name.(string); ok {
										outputArtifacts[nameStr] = true
									}
								}
							}
						}
					}
				}

				// Validate input artifacts
				if inputArtifactsList, hasInput := actionMap["InputArtifacts"]; hasInput && !isIntrinsicFunction(inputArtifactsList) {
					if inputList, ok := inputArtifactsList.([]any); ok {
						for _, input := range inputList {
							if inputMap, ok := input.(map[string]any); ok {
								if name, hasName := inputMap["Name"]; hasName && !isIntrinsicFunction(name) {
									if nameStr, ok := name.(string); ok {
										if !outputArtifacts[nameStr] {
											matches = append(matches, rules.Match{
												Message: fmt.Sprintf(
													"Resource '%s': InputArtifact '%s' does not reference any OutputArtifact",
													resName, nameStr,
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
				}
			}
		}
	}

	return matches
}
