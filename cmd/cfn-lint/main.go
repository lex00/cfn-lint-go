// Package main provides the cfn-lint CLI.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/lex00/cfn-lint-go/pkg/lint"
)

var version = "dev"

func main() {
	if err := rootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	var (
		format  string
		regions []string
	)

	cmd := &cobra.Command{
		Use:   "cfn-lint [templates...]",
		Short: "CloudFormation Linter",
		Long: `cfn-lint validates CloudFormation templates against AWS specifications and best practices.

Examples:
    cfn-lint template.yaml
    cfn-lint template.yaml --format json
    cfn-lint *.yaml --format sarif
    cfn-lint template.yaml --graph > deps.dot`,
		Version: version,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLint(args, format, regions)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text, json, sarif, junit")
	cmd.Flags().StringSliceVarP(&regions, "regions", "r", nil, "AWS regions to validate against")

	cmd.AddCommand(graphCmd())
	cmd.AddCommand(listRulesCmd())

	return cmd
}

func runLint(templates []string, format string, regions []string) error {
	linter := lint.New(lint.Options{
		Regions: regions,
	})

	var allMatches []lint.Match
	for _, path := range templates {
		matches, err := linter.LintFile(path)
		if err != nil {
			return fmt.Errorf("linting %s: %w", path, err)
		}
		allMatches = append(allMatches, matches...)
	}

	return outputMatches(allMatches, format)
}

func outputMatches(matches []lint.Match, format string) error {
	switch format {
	case "text":
		for _, m := range matches {
			fmt.Printf("%s:%d:%d: %s %s [%s]\n",
				m.Filename, m.Line, m.Column,
				m.Level, m.Message, m.Rule)
		}
	case "json":
		// TODO: Implement JSON output
		fmt.Println("[]")
	case "sarif":
		// TODO: Implement SARIF output
		fmt.Println("{}")
	case "junit":
		// TODO: Implement JUnit output
		fmt.Println("<testsuites/>")
	default:
		return fmt.Errorf("unknown format: %s", format)
	}

	if len(matches) > 0 {
		os.Exit(2)
	}
	return nil
}

func graphCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "graph [template]",
		Short: "Generate DOT graph of resource dependencies",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement graph generation
			fmt.Println("digraph G {}")
			return nil
		},
	}
}

func listRulesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-rules",
		Short: "List all available rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: List registered rules
			fmt.Println("No rules registered yet")
			return nil
		},
	}
}
