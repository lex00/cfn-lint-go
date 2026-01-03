// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3636{})
}

// E3636 validates CodeBuild projects using S3 have Location.
type E3636 struct{}

func (r *E3636) ID() string { return "E3636" }

func (r *E3636) ShortDesc() string {
	return "CodeBuild S3 source requires Location"
}

func (r *E3636) Description() string {
	return "Validates that CodeBuild projects with Source Type 'S3' also specify the Location property."
}

func (r *E3636) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3636"
}

func (r *E3636) Tags() []string {
	return []string{"resources", "properties", "codebuild", "source"}
}

func (r *E3636) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::CodeBuild::Project" {
			continue
		}

		source, hasSource := res.Properties["Source"]
		if !hasSource || isIntrinsicFunction(source) {
			continue
		}

		sourceMap, ok := source.(map[string]any)
		if !ok {
			continue
		}

		sourceType, hasType := sourceMap["Type"]
		if !hasType || isIntrinsicFunction(sourceType) {
			continue
		}

		sourceTypeStr, ok := sourceType.(string)
		if !ok {
			continue
		}

		if sourceTypeStr == "S3" {
			if _, hasLocation := sourceMap["Location"]; !hasLocation {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': CodeBuild Project with Source Type 'S3' must specify Location property",
						resName,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", "Source"},
				})
			}
		}
	}

	return matches
}
