package informational

import (
	"fmt"
	"regexp"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&I2003{})
}

// I2003 validates that AllowedPattern in parameters is a valid regular expression.
type I2003 struct{}

func (r *I2003) ID() string { return "I2003" }

func (r *I2003) ShortDesc() string {
	return "Parameter AllowedPattern is valid regex"
}

func (r *I2003) Description() string {
	return "Checks that parameter AllowedPattern values are valid regular expressions. This helps catch pattern syntax errors early."
}

func (r *I2003) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/parameters-section-structure.html"
}

func (r *I2003) Tags() []string {
	return []string{"parameters", "regex", "validation"}
}

func (r *I2003) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for name, param := range tmpl.Parameters {
		if param.AllowedPattern != "" {
			// Try to compile the regex
			_, err := regexp.Compile(param.AllowedPattern)
			if err != nil {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Parameter '%s' has invalid AllowedPattern regex: %v", name, err),
					Path:    []string{"Parameters", name, "AllowedPattern"},
				})
			}
		}
	}

	return matches
}
