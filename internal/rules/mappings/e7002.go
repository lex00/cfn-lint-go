// Package mappings contains mapping validation rules (E7xxx).
package mappings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E7002{})
}

// E7002 checks that mapping names don't exceed the length limit.
type E7002 struct{}

func (r *E7002) ID() string { return "E7002" }

func (r *E7002) ShortDesc() string {
	return "Mapping name length error"
}

func (r *E7002) Description() string {
	return "Checks that mapping names don't exceed 255 characters."
}

func (r *E7002) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cloudformation-limits.html"
}

func (r *E7002) Tags() []string {
	return []string{"mappings", "limits"}
}

const maxMappingNameLength = 255

func (r *E7002) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for mapName, mapping := range tmpl.Mappings {
		if len(mapName) > maxMappingNameLength {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Mapping name '%s' exceeds maximum length of %d characters (got %d)", mapName, maxMappingNameLength, len(mapName)),
				Line:    mapping.Node.Line,
				Column:  mapping.Node.Column,
				Path:    []string{"Mappings", mapName},
			})
		}
	}

	return matches
}
