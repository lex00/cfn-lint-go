package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W1100{})
}

// W1100 warns about YAML merge key usage in CloudFormation templates.
type W1100 struct{}

func (r *W1100) ID() string { return "W1100" }

func (r *W1100) ShortDesc() string {
	return "YAML merge usage"
}

func (r *W1100) Description() string {
	return "Warns when YAML merge keys (<<) are used in CloudFormation templates. While valid YAML, merge keys can make templates harder to understand and debug."
}

func (r *W1100) Source() string {
	return "https://yaml.org/type/merge.html"
}

func (r *W1100) Tags() []string {
	return []string{"warnings", "yaml", "formatting"}
}

func (r *W1100) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check resources for merge key indicators
	for resName, res := range tmpl.Resources {
		r.checkMap(res.Properties, []string{"Resources", resName, "Properties"}, &matches)
	}

	return matches
}

func (r *W1100) checkMap(m map[string]any, path []string, matches *[]rules.Match) {
	for key, value := range m {
		if key == "<<" {
			*matches = append(*matches, rules.Match{
				Message: "YAML merge key (<<) detected; consider using explicit property definitions for clarity",
				Path:    path,
			})
		}
		if nested, ok := value.(map[string]any); ok {
			r.checkMap(nested, append(path, key), matches)
		}
		if list, ok := value.([]any); ok {
			for i, item := range list {
				if nested, ok := item.(map[string]any); ok {
					r.checkMap(nested, append(path, key, fmt.Sprintf("[%d]", i)), matches)
				}
			}
		}
	}
}
