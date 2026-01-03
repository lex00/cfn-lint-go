// Package main provides the cfn-lint CLI.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/lex00/cfn-lint-go/pkg/config"
	"github.com/lex00/cfn-lint-go/pkg/graph"
	"github.com/lex00/cfn-lint-go/pkg/lint"
	"github.com/lex00/cfn-lint-go/pkg/output"
	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"

	// Import rule packages to register them
	_ "github.com/lex00/cfn-lint-go/internal/rules/conditions"
	_ "github.com/lex00/cfn-lint-go/internal/rules/errors"
	_ "github.com/lex00/cfn-lint-go/internal/rules/formats"
	_ "github.com/lex00/cfn-lint-go/internal/rules/functions"
	_ "github.com/lex00/cfn-lint-go/internal/rules/informational"
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
		format              string
		outputFile          string
		noColor             bool
		configFile          string
		regions             []string
		ignoreRules         []string
		includeRules        []string
		includeExperimental bool
	)

	cmd := &cobra.Command{
		Use:   "cfn-lint [templates...]",
		Short: "CloudFormation Linter",
		Long: `cfn-lint validates CloudFormation templates against AWS specifications and best practices.

Examples:
    cfn-lint template.yaml
    cfn-lint template.yaml --format json
    cfn-lint template.yaml --format sarif
    cfn-lint template.yaml --format junit --output results.xml
    cfn-lint template.yaml --format pretty
    cfn-lint *.yaml --ignore-rules E1001,W3002
    cfn-lint template.yaml --config .cfnlintrc.yaml`,
		Version: version,
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLint(args, format, outputFile, configFile, noColor, regions, ignoreRules, includeRules, includeExperimental)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "", "Output format: text, json, sarif, junit, pretty")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write output to file instead of stdout")
	cmd.Flags().BoolVar(&noColor, "no-color", false, "Disable colored output (for pretty format)")
	cmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to config file (.cfnlintrc)")
	cmd.Flags().StringSliceVarP(&regions, "regions", "r", nil, "AWS regions to validate against")
	cmd.Flags().StringSliceVarP(&ignoreRules, "ignore-rules", "i", nil, "Rule IDs to ignore (comma-separated)")
	cmd.Flags().StringSliceVar(&includeRules, "include-checks", nil, "Rule IDs to include (even if ignored)")
	cmd.Flags().BoolVar(&includeExperimental, "include-experimental", false, "Include experimental rules")

	cmd.AddCommand(graphCmd())
	cmd.AddCommand(listRulesCmd())

	return cmd
}

func runLint(templates []string, format, outputFile, configFile string, noColor bool, regions []string, ignoreRules, includeRules []string, includeExperimental bool) error {
	// Load config file if specified or found
	var cfg *config.Config
	if configFile != "" {
		// Explicit config file
		loadedCfg, err := config.Load(configFile)
		if err != nil {
			return fmt.Errorf("loading config file: %w", err)
		}
		cfg = loadedCfg
	} else {
		// Try to find config file
		foundPath, err := config.Find()
		if err == nil {
			loadedCfg, err := config.Load(foundPath)
			if err != nil {
				return fmt.Errorf("loading config file %s: %w", foundPath, err)
			}
			cfg = loadedCfg
		} else {
			// No config file found, use defaults
			cfg = &config.Config{}
		}
	}

	// Merge CLI flags with config (CLI takes precedence)
	cliCfg := &config.Config{
		Templates:           templates,
		Regions:             regions,
		IgnoreChecks:        ignoreRules,
		IncludeChecks:       includeRules,
		IncludeExperimental: includeExperimental,
		Format:              format,
		OutputFile:          outputFile,
	}
	finalCfg := config.Merge(cfg, cliCfg)

	// Determine templates to lint
	templatesToLint := finalCfg.Templates
	if len(templatesToLint) == 0 {
		return fmt.Errorf("no templates specified")
	}

	// Determine effective ignore rules (ignoreChecks - includeChecks)
	effectiveIgnoreRules := make([]string, 0)
	includeSet := make(map[string]bool)
	for _, rule := range finalCfg.IncludeChecks {
		includeSet[rule] = true
	}
	for _, rule := range finalCfg.IgnoreChecks {
		if !includeSet[rule] {
			effectiveIgnoreRules = append(effectiveIgnoreRules, rule)
		}
	}

	linter := lint.New(lint.Options{
		Regions:             finalCfg.Regions,
		IgnoreRules:         effectiveIgnoreRules,
		IncludeExperimental: finalCfg.IncludeExperimental,
	})

	var allMatches []lint.Match
	for _, path := range templatesToLint {
		matches, err := linter.LintFile(path)
		if err != nil {
			return fmt.Errorf("linting %s: %w", path, err)
		}
		allMatches = append(allMatches, matches...)
	}

	// Determine output format
	outFormat := finalCfg.Format
	if outFormat == "" {
		outFormat = "text"
	}

	// Determine output file
	outFile := finalCfg.OutputFile

	// Determine output writer
	writer := os.Stdout
	if outFile != "" {
		f, err := os.Create(outFile)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()
		writer = f
	}

	return outputMatches(writer, allMatches, outFormat, noColor)
}

func outputMatches(w *os.File, matches []lint.Match, format string, noColor bool) error {
	switch format {
	case "text":
		for _, m := range matches {
			fmt.Fprintf(w, "%s:%d:%d: %s %s [%s]\n",
				m.Location.Filename, m.Location.Start.LineNumber, m.Location.Start.ColumnNumber,
				m.Level, m.Message, m.Rule.ID)
		}
	case "json":
		// Ensure we output [] for empty slice, not null
		if matches == nil {
			matches = []lint.Match{}
		}
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		if err := enc.Encode(matches); err != nil {
			return fmt.Errorf("encoding JSON: %w", err)
		}
	case "sarif":
		if err := output.WriteSARIF(w, matches, version); err != nil {
			return fmt.Errorf("encoding SARIF: %w", err)
		}
	case "junit":
		if err := output.WriteJUnit(w, matches); err != nil {
			return fmt.Errorf("encoding JUnit: %w", err)
		}
	case "pretty":
		if err := output.WritePretty(w, matches, noColor); err != nil {
			return fmt.Errorf("writing pretty output: %w", err)
		}
	default:
		return fmt.Errorf("unknown format: %s (valid: text, json, sarif, junit, pretty)", format)
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
			_, _ = fmt.Fprintln(w, "RULE\tDESCRIPTION\tTAGS")
			_, _ = fmt.Fprintln(w, "----\t-----------\t----")
			for _, r := range allRules {
				tags := ""
				if len(r.Tags()) > 0 {
					tags = fmt.Sprintf("%v", r.Tags())
				}
				_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", r.ID(), r.ShortDesc(), tags)
			}
			return w.Flush()
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text, json")

	return cmd
}
