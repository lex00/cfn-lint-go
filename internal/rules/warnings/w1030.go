package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W1030{})
}

// W1030 warns about Ref function value issues that may cause problems.
type W1030 struct{}

func (r *W1030) ID() string { return "W1030" }

func (r *W1030) ShortDesc() string {
	return "Ref function value validation"
}

func (r *W1030) Description() string {
	return "Warns about potential issues with Ref function values, such as empty strings or whitespace."
}

func (r *W1030) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-ref.html"
}

func (r *W1030) Tags() []string {
	return []string{"warnings", "functions", "ref"}
}

func (r *W1030) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check resources
	for resName, res := range tmpl.Resources {
		r.checkValue(res.Properties, []string{"Resources", resName, "Properties"}, &matches)
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		r.checkValue(out.Value, []string{"Outputs", outName, "Value"}, &matches)
	}

	return matches
}

func (r *W1030) checkValue(v any, path []string, matches *[]rules.Match) {
	switch val := v.(type) {
	case map[string]any:
		if ref, ok := val["Ref"]; ok {
			r.checkRef(ref, path, matches)
		}
		for key, child := range val {
			r.checkValue(child, append(path, key), matches)
		}
	case []any:
		for i, child := range val {
			r.checkValue(child, append(path, fmt.Sprintf("[%d]", i)), matches)
		}
	}
}

func (r *W1030) checkRef(ref any, path []string, matches *[]rules.Match) {
	refStr, ok := ref.(string)
	if !ok {
		return
	}

	// Check for whitespace in reference name
	for _, c := range refStr {
		if c == ' ' || c == '\t' || c == '\n' {
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf("Ref value '%s' contains whitespace which may cause issues", refStr),
				Path:    path,
			})
			break
		}
	}
}
