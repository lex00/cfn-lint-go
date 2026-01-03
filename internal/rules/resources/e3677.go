// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3677{})
}

// E3677 validates Lambda ZipFile runtime.
type E3677 struct{}

func (r *E3677) ID() string { return "E3677" }

func (r *E3677) ShortDesc() string {
	return "Lambda ZipFile supports JavaScript and Python only"
}

func (r *E3677) Description() string {
	return "Validates that AWS::Lambda::Function resources using ZipFile specify a JavaScript or Python runtime."
}

func (r *E3677) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3677"
}

func (r *E3677) Tags() []string {
	return []string{"resources", "properties", "lambda", "zipfile"}
}

func (r *E3677) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Lambda::Function" && res.Type != "AWS::Serverless::Function" {
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

		runtime, hasRuntime := res.Properties["Runtime"]
		if !hasRuntime || isIntrinsicFunction(runtime) {
			continue
		}

		runtimeStr, ok := runtime.(string)
		if !ok {
			continue
		}

		// ZipFile only supports nodejs* and python* runtimes
		if !strings.HasPrefix(runtimeStr, "nodejs") && !strings.HasPrefix(runtimeStr, "python") {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Lambda ZipFile only supports nodejs* and python* runtimes. Got: %s",
					resName, runtimeStr,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties", "Runtime"},
			})
		}
	}

	return matches
}
