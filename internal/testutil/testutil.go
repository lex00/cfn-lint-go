// Package testutil provides test utilities for cfn-lint-go.
package testutil

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/lint"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

// ProjectRoot returns the root directory of the project.
// This is useful for locating test fixtures.
func ProjectRoot() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to get project root")
	}
	// testutil.go is at internal/testutil/testutil.go
	// so project root is 2 levels up
	return filepath.Dir(filepath.Dir(filepath.Dir(filename)))
}

// TestDataDir returns the path to the testdata directory.
func TestDataDir() string {
	return filepath.Join(ProjectRoot(), "testdata")
}

// TemplatesDir returns the path to the testdata/templates directory.
func TemplatesDir() string {
	return filepath.Join(TestDataDir(), "templates")
}

// GoodTemplatesDir returns the path to testdata/templates/good.
func GoodTemplatesDir() string {
	return filepath.Join(TemplatesDir(), "good")
}

// BadTemplatesDir returns the path to testdata/templates/bad.
func BadTemplatesDir() string {
	return filepath.Join(TemplatesDir(), "bad")
}

// SAMTemplatesDir returns the path to testdata/templates/sam.
func SAMTemplatesDir() string {
	return filepath.Join(TemplatesDir(), "sam")
}

// IssuesTemplatesDir returns the path to testdata/templates/issues.
func IssuesTemplatesDir() string {
	return filepath.Join(TemplatesDir(), "issues")
}

// LoadTemplate loads and parses a template from a file.
func LoadTemplate(t *testing.T, path string) *template.Template {
	t.Helper()
	tmpl, err := template.ParseFile(path)
	if err != nil {
		t.Fatalf("Failed to parse template %s: %v", path, err)
	}
	return tmpl
}

// LoadTemplateBytes loads and parses a template from bytes.
func LoadTemplateBytes(t *testing.T, content []byte) *template.Template {
	t.Helper()
	tmpl, err := template.Parse(content)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}
	return tmpl
}

// LintFile lints a template file and returns the matches.
func LintFile(t *testing.T, path string, opts lint.Options) []lint.Match {
	t.Helper()
	linter := lint.New(opts)
	matches, err := linter.LintFile(path)
	if err != nil {
		t.Fatalf("Failed to lint file %s: %v", path, err)
	}
	return matches
}

// LintTemplate lints a parsed template and returns the matches.
func LintTemplate(t *testing.T, tmpl *template.Template, filename string, opts lint.Options) []lint.Match {
	t.Helper()
	linter := lint.New(opts)
	matches, err := linter.Lint(tmpl, filename)
	if err != nil {
		t.Fatalf("Failed to lint template: %v", err)
	}
	return matches
}

// ListTemplates returns all template files in a directory.
func ListTemplates(t *testing.T, dir string) []string {
	t.Helper()
	var templates []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext == ".yaml" || ext == ".yml" || ext == ".json" {
			templates = append(templates, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to list templates in %s: %v", dir, err)
	}
	return templates
}

// AssertNoErrors asserts that there are no Error-level matches.
func AssertNoErrors(t *testing.T, matches []lint.Match) {
	t.Helper()
	for _, m := range matches {
		if m.Level == "Error" {
			t.Errorf("Unexpected error: %s - %s at %s:%d",
				m.Rule.ID, m.Message, m.Location.Filename, m.Location.Start.LineNumber)
		}
	}
}

// AssertHasError asserts that there is at least one match with the given rule ID.
func AssertHasError(t *testing.T, matches []lint.Match, ruleID string) {
	t.Helper()
	for _, m := range matches {
		if m.Rule.ID == ruleID {
			return
		}
	}
	t.Errorf("Expected error with rule %s, but none found", ruleID)
}

// AssertMatchCount asserts that there are exactly n matches.
func AssertMatchCount(t *testing.T, matches []lint.Match, expected int) {
	t.Helper()
	if len(matches) != expected {
		t.Errorf("Expected %d matches, got %d", expected, len(matches))
		for _, m := range matches {
			t.Logf("  - %s: %s", m.Rule.ID, m.Message)
		}
	}
}

// FilterByLevel returns matches filtered by level (Error, Warning, Informational).
func FilterByLevel(matches []lint.Match, level string) []lint.Match {
	var filtered []lint.Match
	for _, m := range matches {
		if m.Level == level {
			filtered = append(filtered, m)
		}
	}
	return filtered
}

// FilterByRuleID returns matches filtered by rule ID prefix.
func FilterByRuleID(matches []lint.Match, prefix string) []lint.Match {
	var filtered []lint.Match
	for _, m := range matches {
		if len(m.Rule.ID) >= len(prefix) && m.Rule.ID[:len(prefix)] == prefix {
			filtered = append(filtered, m)
		}
	}
	return filtered
}

// TemplatePath returns the full path to a template in testdata/templates.
func TemplatePath(subdir, filename string) string {
	return filepath.Join(TemplatesDir(), subdir, filename)
}
