// Package lint provides the public API for CloudFormation template linting.
//
// This is the main entry point for using cfn-lint-go as a library. It provides
// a simple interface to lint CloudFormation templates and get structured results
// compatible with the Python cfn-lint tool.
//
// # Basic Usage
//
// Create a linter and lint a file:
//
//	linter := lint.New(lint.Options{})
//	matches, err := linter.LintFile("template.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, m := range matches {
//	    fmt.Printf("[%s] %s\n", m.Rule.ID, m.Message)
//	}
//
// # Configuration
//
// The linter can be configured with Options:
//
//	linter := lint.New(lint.Options{
//	    Regions:     []string{"us-east-1", "us-west-2"},
//	    IgnoreRules: []string{"W1001", "W2001"},
//	    IncludeExperimental: true,
//	})
//
// # Match Results
//
// Each Match contains:
//   - Rule: metadata about the rule that matched (ID, description)
//   - Location: file, line, column, and JSON path to the issue
//   - Level: "Error", "Warning", or "Informational"
//   - Message: human-readable description of the issue
//
// Rule IDs follow the Python cfn-lint convention:
//   - E#### - Errors (template won't work)
//   - W#### - Warnings (best practice violations)
//   - I#### - Informational (suggestions)
package lint
