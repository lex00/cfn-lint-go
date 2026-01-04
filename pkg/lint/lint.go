// Package lint provides the public API for CloudFormation template linting.
package lint

import (
	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/sam"
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

	// DisableSAMTransform disables automatic SAM transformation.
	// When true, SAM templates are linted as-is without transformation.
	DisableSAMTransform bool

	// SAMTransformOptions configures SAM transformation behavior.
	SAMTransformOptions *sam.TransformOptions
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
	// Check if SAM transformation is needed
	if sam.IsSAMTemplate(tmpl) && !l.options.DisableSAMTransform {
		return l.lintSAM(tmpl, filename)
	}

	return l.lintCloudFormation(tmpl, filename, nil)
}

// lintSAM handles SAM template linting with transformation.
func (l *Linter) lintSAM(tmpl *template.Template, filename string) ([]Match, error) {
	// Transform SAM to CloudFormation
	result, err := sam.Transform(tmpl, l.options.SAMTransformOptions)
	if err != nil {
		// Return SAM transform error as a lint error
		return []Match{{
			Rule: MatchRule{
				ID:               "E0010",
				Description:      "Checks that SAM templates can be transformed",
				ShortDescription: "SAM transform failed",
				Source:           "https://github.com/lex00/cfn-lint-go",
			},
			Location: MatchLocation{
				Start:    MatchPosition{LineNumber: 1, ColumnNumber: 1},
				End:      MatchPosition{LineNumber: 1, ColumnNumber: 1},
				Path:     []any{},
				Filename: filename,
			},
			Level:   "Error",
			Message: err.Error(),
		}}, nil
	}

	// Lint the transformed template
	return l.lintCloudFormation(result.Template, filename, result.SourceMap)
}

// lintCloudFormation lints a CloudFormation template with optional source mapping.
func (l *Linter) lintCloudFormation(tmpl *template.Template, filename string, sourceMap *sam.SourceMap) ([]Match, error) {
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

			// Get line/column, potentially mapping back to SAM source
			line, column := rm.Line, rm.Column
			if sourceMap != nil && len(rm.Path) >= 2 {
				// Try to get resource name from path
				if rm.Path[0] == "Resources" {
					resourceName := rm.Path[1]
					mappedLoc := sourceMap.MapError(resourceName, "", rm.Line)
					if mappedLoc.Line > 0 {
						line = mappedLoc.Line
						column = mappedLoc.Column
					}
				}
			}

			matches = append(matches, Match{
				Rule: MatchRule{
					ID:               rule.ID(),
					Description:      rule.Description(),
					ShortDescription: rule.ShortDesc(),
					Source:           rule.Source(),
				},
				Location: MatchLocation{
					Start:    MatchPosition{LineNumber: line, ColumnNumber: column},
					End:      MatchPosition{LineNumber: line, ColumnNumber: column},
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
