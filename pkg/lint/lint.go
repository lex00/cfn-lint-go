// Package lint provides the public API for CloudFormation template linting.
package lint

import (
	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

// Linter validates CloudFormation templates.
type Linter struct {
	options Options
	rules   []rules.Rule
}

// Options configures the linter.
type Options struct {
	// Regions to validate against. Empty means all regions.
	Regions []string

	// IgnoreRules is a list of rule IDs to skip.
	IgnoreRules []string

	// IncludeExperimental enables experimental rules.
	IncludeExperimental bool
}

// Match represents a linting issue found in a template (Python cfn-lint compatible format).
type Match struct {
	Rule     MatchRule     `json:"Rule"`
	Location MatchLocation `json:"Location"`
	Level    string        `json:"Level"` // "Error", "Warning", "Informational"
	Message  string        `json:"Message"`
}

// MatchRule contains rule metadata.
type MatchRule struct {
	ID               string `json:"Id"`
	Description      string `json:"Description"`
	ShortDescription string `json:"ShortDescription"`
	Source           string `json:"Source"`
}

// MatchLocation contains the location of the issue.
type MatchLocation struct {
	Start    MatchPosition `json:"Start"`
	End      MatchPosition `json:"End"`
	Path     []any         `json:"Path"`
	Filename string        `json:"Filename"`
}

// MatchPosition represents a line/column position.
type MatchPosition struct {
	LineNumber   int `json:"LineNumber"`
	ColumnNumber int `json:"ColumnNumber"`
}

// New creates a new Linter with the given options.
func New(opts Options) *Linter {
	return &Linter{
		options: opts,
		rules:   rules.All(),
	}
}

// LintFile lints a CloudFormation template file.
func (l *Linter) LintFile(path string) ([]Match, error) {
	tmpl, err := template.ParseFile(path)
	if err != nil {
		return []Match{{
			Rule: MatchRule{
				ID:               "E0000",
				Description:      "Checks that the template can be parsed",
				ShortDescription: "Template parse error",
				Source:           "https://github.com/lex00/cfn-lint-go",
			},
			Location: MatchLocation{
				Start:    MatchPosition{LineNumber: 1, ColumnNumber: 1},
				End:      MatchPosition{LineNumber: 1, ColumnNumber: 1},
				Path:     []any{},
				Filename: path,
			},
			Level:   "Error",
			Message: err.Error(),
		}}, nil
	}

	return l.Lint(tmpl, path)
}

// Lint lints a parsed CloudFormation template.
func (l *Linter) Lint(tmpl *template.Template, filename string) ([]Match, error) {
	var matches []Match

	for _, rule := range l.rules {
		if l.isIgnored(rule.ID()) {
			continue
		}

		ruleMatches := rule.Match(tmpl)
		for _, rm := range ruleMatches {
			// Convert []string path to []any for JSON compatibility
			path := make([]any, len(rm.Path))
			for i, p := range rm.Path {
				path[i] = p
			}

			matches = append(matches, Match{
				Rule: MatchRule{
					ID:               rule.ID(),
					Description:      rule.Description(),
					ShortDescription: rule.ShortDesc(),
					Source:           rule.Source(),
				},
				Location: MatchLocation{
					Start:    MatchPosition{LineNumber: rm.Line, ColumnNumber: rm.Column},
					End:      MatchPosition{LineNumber: rm.Line, ColumnNumber: rm.Column},
					Path:     path,
					Filename: filename,
				},
				Level:   levelFromRuleID(rule.ID()),
				Message: rm.Message,
			})
		}
	}

	return matches, nil
}

func (l *Linter) isIgnored(ruleID string) bool {
	for _, ignored := range l.options.IgnoreRules {
		if ignored == ruleID {
			return true
		}
	}
	return false
}

func levelFromRuleID(id string) string {
	if len(id) == 0 {
		return "Error"
	}
	switch id[0] {
	case 'E':
		return "Error"
	case 'W':
		return "Warning"
	case 'I':
		return "Informational"
	default:
		return "Error"
	}
}
