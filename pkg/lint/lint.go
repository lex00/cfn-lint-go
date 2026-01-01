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

// Match represents a linting issue found in a template.
type Match struct {
	Rule     string `json:"rule"`
	Message  string `json:"message"`
	Level    string `json:"level"` // error, warning, info
	Filename string `json:"filename"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
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
			Rule:     "E0000",
			Message:  err.Error(),
			Level:    "error",
			Filename: path,
			Line:     1,
			Column:   1,
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
			matches = append(matches, Match{
				Rule:     rule.ID(),
				Message:  rm.Message,
				Level:    levelFromRuleID(rule.ID()),
				Filename: filename,
				Line:     rm.Line,
				Column:   rm.Column,
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
		return "error"
	}
	switch id[0] {
	case 'E':
		return "error"
	case 'W':
		return "warning"
	case 'I':
		return "info"
	default:
		return "error"
	}
}
