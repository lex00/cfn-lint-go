// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3678{})
}

// E3678 validates ZipFile requires runtime.
type E3678 struct{}

func (r *E3678) ID() string { return "E3678" }

func (r *E3678) ShortDesc() string {
	return "ZipFile requires runtime"
}

func (r *E3678) Description() string {
	return "Validates that AWS::Lambda::Function resources using ZipFile specify a Runtime."
}

func (r *E3678) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3678"
}

func (r *E3678) Tags() []string {
	return []string{"resources", "properties", "lambda", "zipfile"}
}

func (r *E3678) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Lambda::Function" {
			continue
		}

		code, hasCode := res.Properties["Code"]
		if !hasCode || isIntrinsicFunction(code) {
			continue
		}

		codeMap, ok := code.(map[string]any)
		if !ok {
			continue
		}

		_, hasZipFile := codeMap["ZipFile"]
		if !hasZipFile {
			continue
		}

		_, hasRuntime := res.Properties["Runtime"]
		if !hasRuntime {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Lambda Function using ZipFile must specify Runtime",
					resName,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties"},
			})
		}
	}

	return matches
}
