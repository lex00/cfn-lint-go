// Package outputs contains output validation rules (E6xxx).
package outputs

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E6102{})
}

// E6102 checks that output Export Name values resolve to strings.
type E6102 struct{}

func (r *E6102) ID() string { return "E6102" }

func (r *E6102) ShortDesc() string {
	return "Export Name must be a string"
}

func (r *E6102) Description() string {
	return "Checks that output Export Name values resolve to string types."
}

func (r *E6102) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/outputs-section-structure.html"
}

func (r *E6102) Tags() []string {
	return []string{"outputs", "exports", "types"}
}

func (r *E6102) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for name, out := range tmpl.Outputs {
		if out.Export == nil {
			continue // No Export to validate
		}

		exportName, hasName := out.Export["Name"]
		if !hasName {
			continue // E6003 handles missing Export.Name
		}

		if !isStringOrIntrinsic(exportName) {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Output '%s' Export Name must be a string, got %T", name, exportName),
				Line:    out.Node.Line,
				Column:  out.Node.Column,
				Path:    []string{"Outputs", name, "Export", "Name"},
			})
		}
	}

	return matches
}
