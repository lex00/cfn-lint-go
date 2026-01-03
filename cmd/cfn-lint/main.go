// Package main provides the cfn-lint CLI.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/lex00/cfn-lint-go/pkg/graph"
	"github.com/lex00/cfn-lint-go/pkg/lint"
	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"

	// Import rule packages to register them
	_ "github.com/lex00/cfn-lint-go/internal/rules/conditions"
	_ "github.com/lex00/cfn-lint-go/internal/rules/errors"
	_ "github.com/lex00/cfn-lint-go/internal/rules/functions"
	_ "github.com/lex00/cfn-lint-go/internal/rules/mappings"
	_ "github.com/lex00/cfn-lint-go/internal/rules/metadata"
	_ "github.com/lex00/cfn-lint-go/internal/rules/outputs"
	_ "github.com/lex00/cfn-lint-go/internal/rules/parameters"
	_ "github.com/lex00/cfn-lint-go/internal/rules/resources"
	_ "github.com/lex00/cfn-lint-go/internal/rules/rulessection"
	_ "github.com/lex00/cfn-lint-go/internal/rules/warnings"
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
		format      string
		regions     []string
		ignoreRules []string
	)

	cmd := &cobra.Command{
		Use:   "cfn-lint [templates...]",
		Short: "CloudFormation Linter",
		Long: `cfn-lint validates CloudFormation templates against AWS specifications and best practices.

Examples:
    cfn-lint template.yaml
    cfn-lint template.yaml --format json
    cfn-lint *.yaml --ignore-rules E1001,W3002`,
		Version: version,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLint(args, format, regions, ignoreRules)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text, json")
	cmd.Flags().StringSliceVarP(&regions, "regions", "r", nil, "AWS regions to validate against")
	cmd.Flags().StringSliceVarP(&ignoreRules, "ignore-rules", "i", nil, "Rule IDs to ignore (comma-separated)")

	cmd.AddCommand(graphCmd())
	cmd.AddCommand(listRulesCmd())

	return cmd
}

func runLint(templates []string, format string, regions []string, ignoreRules []string) error {
	linter := lint.New(lint.Options{
		Regions:     regions,
		IgnoreRules: ignoreRules,
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
				m.Location.Filename, m.Location.Start.LineNumber, m.Location.Start.ColumnNumber,
				m.Level, m.Message, m.Rule.ID)
		}
	case "json":
		// Ensure we output [] for empty slice, not null
		if matches == nil {
			matches = []lint.Match{}
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(matches); err != nil {
			return fmt.Errorf("encoding JSON: %w", err)
		}
	default:
		return fmt.Errorf("unknown format: %s (valid: text, json)", format)
	}

	if len(matches) > 0 {
		os.Exit(2)
	}
	return nil
}

func graphCmd() *cobra.Command {
	var includeParams bool

	cmd := &cobra.Command{
		Use:   "graph [template]",
		Short: "Generate DOT graph of resource dependencies",
		Long: `Generate a DOT format graph showing resource dependencies.

The output can be rendered with Graphviz:
    cfn-lint graph template.yaml | dot -Tpng -o deps.png`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tmpl, err := template.ParseFile(args[0])
			if err != nil {
				return fmt.Errorf("parsing template: %w", err)
			}

			gen := &graph.Generator{
				IncludeParameters: includeParams,
			}
			return gen.Generate(tmpl, os.Stdout)
		},
	}

	cmd.Flags().BoolVarP(&includeParams, "include-parameters", "p", false, "Include parameter nodes in the graph")

	return cmd
}

func listRulesCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "list-rules",
		Short: "List all available rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			allRules := rules.All()

			if format == "json" {
				type ruleInfo struct {
					ID          string   `json:"id"`
					ShortDesc   string   `json:"short_desc"`
					Description string   `json:"description"`
					Tags        []string `json:"tags"`
					Source      string   `json:"source,omitempty"`
				}
				var ruleList []ruleInfo
				for _, r := range allRules {
					ruleList = append(ruleList, ruleInfo{
						ID:          r.ID(),
						ShortDesc:   r.ShortDesc(),
						Description: r.Description(),
						Tags:        r.Tags(),
						Source:      r.Source(),
					})
				}
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(ruleList)
			}

			// Text format
			if len(allRules) == 0 {
				fmt.Println("No rules registered.")
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "RULE\tDESCRIPTION\tTAGS")
			fmt.Fprintln(w, "----\t-----------\t----")
			for _, r := range allRules {
				tags := ""
				if len(r.Tags()) > 0 {
					tags = fmt.Sprintf("%v", r.Tags())
				}
				fmt.Fprintf(w, "%s\t%s\t%s\n", r.ID(), r.ShortDesc(), tags)
			}
			return w.Flush()
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text, json")

	return cmd
}
