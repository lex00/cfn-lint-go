package output

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/lint"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
)

// WritePretty writes matches in a pretty, colorized format with context.
func WritePretty(w io.Writer, matches []lint.Match, noColor bool) error {
	if len(matches) == 0 {
		if !noColor {
			fmt.Fprintf(w, "%s✓ No issues found%s\n", colorBlue, colorReset)
		} else {
			fmt.Fprintln(w, "✓ No issues found")
		}
		return nil
	}

	// Group matches by file
	fileMatches := make(map[string][]lint.Match)
	for _, m := range matches {
		filename := m.Location.Filename
		if filename == "" {
			filename = "unknown"
		}
		fileMatches[filename] = append(fileMatches[filename], m)
	}

	// Sort filenames for consistent output
	var filenames []string
	for filename := range fileMatches {
		filenames = append(filenames, filename)
	}
	sort.Strings(filenames)

	// Track counts
	errorCount := 0
	warningCount := 0
	infoCount := 0

	// Process each file
	for i, filename := range filenames {
		matches := fileMatches[filename]

		// Sort matches by line number
		sort.Slice(matches, func(i, j int) bool {
			return matches[i].Location.Start.LineNumber < matches[j].Location.Start.LineNumber
		})

		// Print file header
		if i > 0 {
			fmt.Fprintln(w)
		}
		if !noColor {
			fmt.Fprintf(w, "%s%s%s%s\n", colorBold, colorBlue, filename, colorReset)
		} else {
			fmt.Fprintf(w, "%s\n", filename)
		}
		fmt.Fprintln(w, strings.Repeat("─", len(filename)))

		// Read file for context (best effort)
		fileLines := readFileLines(filename)

		// Print each match
		for _, m := range matches {
			levelColor := colorRed
			levelSymbol := "✖"

			switch m.Level {
			case "Error":
				errorCount++
				levelColor = colorRed
				levelSymbol = "✖"
			case "Warning":
				warningCount++
				levelColor = colorYellow
				levelSymbol = "⚠"
			case "Informational":
				infoCount++
				levelColor = colorBlue
				levelSymbol = "ℹ"
			}

			// Print location and rule
			if !noColor {
				fmt.Fprintf(w, "\n  %s%s%s Line %d:%d - %s[%s]%s\n",
					levelColor, levelSymbol, colorReset,
					m.Location.Start.LineNumber,
					m.Location.Start.ColumnNumber,
					colorGray, m.Rule.ID, colorReset)
			} else {
				fmt.Fprintf(w, "\n  %s Line %d:%d - [%s]\n",
					levelSymbol,
					m.Location.Start.LineNumber,
					m.Location.Start.ColumnNumber,
					m.Rule.ID)
			}

			// Print message
			if !noColor {
				fmt.Fprintf(w, "  %s\n", m.Message)
			} else {
				fmt.Fprintf(w, "  %s\n", m.Message)
			}

			// Print code context if available
			if len(fileLines) > 0 && m.Location.Start.LineNumber > 0 {
				printContext(w, fileLines, m.Location.Start.LineNumber, noColor)
			}
		}
	}

	// Print summary
	fmt.Fprintln(w)
	fmt.Fprintln(w, strings.Repeat("═", 60))
	if !noColor {
		fmt.Fprintf(w, "%sSummary:%s ", colorBold, colorReset)
		if errorCount > 0 {
			fmt.Fprintf(w, "%s%d errors%s ", colorRed, errorCount, colorReset)
		}
		if warningCount > 0 {
			fmt.Fprintf(w, "%s%d warnings%s ", colorYellow, warningCount, colorReset)
		}
		if infoCount > 0 {
			fmt.Fprintf(w, "%s%d info%s", colorBlue, infoCount, colorReset)
		}
		fmt.Fprintln(w)
	} else {
		fmt.Fprintf(w, "Summary: %d errors, %d warnings, %d info\n", errorCount, warningCount, infoCount)
	}

	return nil
}

// readFileLines reads a file and returns its lines.
func readFileLines(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

// printContext prints lines around the error location.
func printContext(w io.Writer, lines []string, lineNum int, noColor bool) {
	contextBefore := 2
	contextAfter := 2

	start := lineNum - contextBefore - 1
	if start < 0 {
		start = 0
	}

	end := lineNum + contextAfter
	if end > len(lines) {
		end = len(lines)
	}

	// Find max line number width for padding
	maxLineNum := end
	lineNumWidth := len(fmt.Sprintf("%d", maxLineNum))

	fmt.Fprintln(w)
	for i := start; i < end; i++ {
		lineNumStr := fmt.Sprintf("%*d", lineNumWidth, i+1)

		if i+1 == lineNum {
			// Highlight the error line
			if !noColor {
				fmt.Fprintf(w, "    %s%s │ %s%s\n", colorRed, lineNumStr, lines[i], colorReset)
			} else {
				fmt.Fprintf(w, "  > %s │ %s\n", lineNumStr, lines[i])
			}
		} else {
			// Context lines
			if !noColor {
				fmt.Fprintf(w, "    %s%s │%s %s\n", colorGray, lineNumStr, colorReset, lines[i])
			} else {
				fmt.Fprintf(w, "    %s │ %s\n", lineNumStr, lines[i])
			}
		}
	}
}
