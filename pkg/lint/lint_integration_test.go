package lint_test

import (
	"path/filepath"
	"testing"

	"github.com/lex00/cfn-lint-go/internal/testutil"
	"github.com/lex00/cfn-lint-go/pkg/lint"

	// Import rule packages to register them
	_ "github.com/lex00/cfn-lint-go/internal/rules/conditions"
	_ "github.com/lex00/cfn-lint-go/internal/rules/errors"
	_ "github.com/lex00/cfn-lint-go/internal/rules/formats"
	_ "github.com/lex00/cfn-lint-go/internal/rules/functions"
	_ "github.com/lex00/cfn-lint-go/internal/rules/informational"
	_ "github.com/lex00/cfn-lint-go/internal/rules/mappings"
	_ "github.com/lex00/cfn-lint-go/internal/rules/metadata"
	_ "github.com/lex00/cfn-lint-go/internal/rules/modules"
	_ "github.com/lex00/cfn-lint-go/internal/rules/outputs"
	_ "github.com/lex00/cfn-lint-go/internal/rules/parameters"
	_ "github.com/lex00/cfn-lint-go/internal/rules/resources"
	_ "github.com/lex00/cfn-lint-go/internal/rules/rulessection"
	_ "github.com/lex00/cfn-lint-go/internal/rules/warnings"
)

// TestGoodTemplates ensures all templates in testdata/templates/good pass linting without errors.
func TestGoodTemplates(t *testing.T) {
	templates := testutil.ListTemplates(t, testutil.GoodTemplatesDir())
	if len(templates) == 0 {
		t.Skip("No good templates found")
	}

	for _, tmplPath := range templates {
		name := filepath.Base(tmplPath)
		t.Run(name, func(t *testing.T) {
			matches := testutil.LintFile(t, tmplPath, lint.Options{})
			errors := testutil.FilterByLevel(matches, "Error")
			testutil.AssertMatchCount(t, errors, 0)
		})
	}
}

// TestBadTemplates ensures all templates in testdata/templates/bad produce at least one error.
func TestBadTemplates(t *testing.T) {
	templates := testutil.ListTemplates(t, testutil.BadTemplatesDir())
	if len(templates) == 0 {
		t.Skip("No bad templates found")
	}

	for _, tmplPath := range templates {
		name := filepath.Base(tmplPath)
		t.Run(name, func(t *testing.T) {
			matches := testutil.LintFile(t, tmplPath, lint.Options{})
			errors := testutil.FilterByLevel(matches, "Error")
			if len(errors) == 0 {
				t.Errorf("Expected at least one error for bad template %s, got none", name)
			}
		})
	}
}

// TestSAMTemplates ensures SAM templates can be parsed and linted.
func TestSAMTemplates(t *testing.T) {
	templates := testutil.ListTemplates(t, testutil.SAMTemplatesDir())
	if len(templates) == 0 {
		t.Skip("No SAM templates found")
	}

	for _, tmplPath := range templates {
		name := filepath.Base(tmplPath)
		t.Run(name, func(t *testing.T) {
			// SAM templates should parse without errors
			_ = testutil.LoadTemplate(t, tmplPath)

			// Lint the template
			matches := testutil.LintFile(t, tmplPath, lint.Options{})

			// Log any issues found (for visibility, not failures)
			for _, m := range matches {
				t.Logf("  [%s] %s: %s", m.Level, m.Rule.ID, m.Message)
			}
		})
	}
}

// TestIssueRegressionTemplates runs regression tests for previously fixed issues.
func TestIssueRegressionTemplates(t *testing.T) {
	templates := testutil.ListTemplates(t, testutil.IssuesTemplatesDir())
	if len(templates) == 0 {
		t.Skip("No issue templates found")
	}

	for _, tmplPath := range templates {
		name := filepath.Base(tmplPath)
		t.Run(name, func(t *testing.T) {
			// These templates should parse without crashing
			tmpl := testutil.LoadTemplate(t, tmplPath)
			if tmpl == nil {
				t.Fatal("Template loaded as nil")
			}

			// Lint should complete without panicking
			matches := testutil.LintFile(t, tmplPath, lint.Options{})

			// Log any issues
			for _, m := range matches {
				t.Logf("  [%s] %s: %s", m.Level, m.Rule.ID, m.Message)
			}
		})
	}
}

// TestSpecificBadTemplateErrors verifies specific templates produce expected errors.
func TestSpecificBadTemplateErrors(t *testing.T) {
	tests := []struct {
		name          string
		templateFile  string
		expectedRules []string // Expected rule IDs to be triggered
	}{
		{
			name:          "invalid_ref",
			templateFile:  "invalid_ref.yaml",
			expectedRules: []string{"E1001", "E1010"},
		},
		{
			name:          "missing_type",
			templateFile:  "missing_type.yaml",
			expectedRules: []string{"E3001"},
		},
		{
			name:          "circular_dependency",
			templateFile:  "circular_dependency.yaml",
			expectedRules: []string{"E3004"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := testutil.TemplatePath("bad", tc.templateFile)
			matches := testutil.LintFile(t, path, lint.Options{})

			for _, expectedRule := range tc.expectedRules {
				testutil.AssertHasError(t, matches, expectedRule)
			}
		})
	}
}

// TestIgnoreRules verifies that ignored rules don't produce matches.
func TestIgnoreRules(t *testing.T) {
	path := testutil.TemplatePath("bad", "invalid_ref.yaml")

	// First, lint without ignoring rules
	matchesWithRules := testutil.LintFile(t, path, lint.Options{})
	if len(matchesWithRules) == 0 {
		t.Skip("Template produces no matches, can't test ignore")
	}

	// Get the first rule ID
	firstRule := matchesWithRules[0].Rule.ID

	// Now lint with that rule ignored
	matchesIgnored := testutil.LintFile(t, path, lint.Options{
		IgnoreRules: []string{firstRule},
	})

	// Count matches for the ignored rule
	for _, m := range matchesIgnored {
		if m.Rule.ID == firstRule {
			t.Errorf("Rule %s was supposed to be ignored but produced a match", firstRule)
		}
	}
}

// TestLevelDistribution verifies match levels are set correctly.
func TestLevelDistribution(t *testing.T) {
	templates := testutil.ListTemplates(t, testutil.GoodTemplatesDir())
	if len(templates) == 0 {
		t.Skip("No templates found")
	}

	// Use the first template
	path := templates[0]
	matches := testutil.LintFile(t, path, lint.Options{})

	// Check that all matches have valid levels
	for _, m := range matches {
		switch m.Level {
		case "Error", "Warning", "Informational":
			// Valid
		default:
			t.Errorf("Invalid level %q for rule %s", m.Level, m.Rule.ID)
		}

		// Verify level matches rule ID prefix
		if len(m.Rule.ID) > 0 {
			prefix := string(m.Rule.ID[0])
			expectedLevel := "Error"
			switch prefix {
			case "E":
				expectedLevel = "Error"
			case "W":
				expectedLevel = "Warning"
			case "I":
				expectedLevel = "Informational"
			}
			if m.Level != expectedLevel {
				t.Errorf("Rule %s has level %s but expected %s", m.Rule.ID, m.Level, expectedLevel)
			}
		}
	}
}
