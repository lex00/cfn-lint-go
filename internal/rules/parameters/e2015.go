// Package parameters contains parameter validation rules (E2xxx).
package parameters

import (
	"fmt"
	"regexp"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E2015{})
}

// E2015 checks that parameter default values satisfy all constraints.
type E2015 struct{}

func (r *E2015) ID() string { return "E2015" }

func (r *E2015) ShortDesc() string {
	return "Default value is within constraints"
}

func (r *E2015) Description() string {
	return "Validates that parameter Default values satisfy AllowedValues, AllowedPattern, MinValue/MaxValue, and MinLength/MaxLength constraints."
}

func (r *E2015) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E2015"
}

func (r *E2015) Tags() []string {
	return []string{"parameters", "default", "constraints"}
}

func (r *E2015) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for name, param := range tmpl.Parameters {
		if param.Default == nil {
			continue
		}

		// Check AllowedValues
		if len(param.AllowedValues) > 0 {
			found := false
			for _, allowed := range param.AllowedValues {
				if fmt.Sprintf("%v", param.Default) == fmt.Sprintf("%v", allowed) {
					found = true
					break
				}
			}
			if !found {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Parameter '%s' default value '%v' is not in AllowedValues", name, param.Default),
					Line:    param.Node.Line,
					Column:  param.Node.Column,
					Path:    []string{"Parameters", name, "Default"},
				})
			}
		}

		// Check AllowedPattern
		if param.AllowedPattern != "" {
			defaultStr := fmt.Sprintf("%v", param.Default)
			re, err := regexp.Compile(param.AllowedPattern)
			if err == nil && !re.MatchString(defaultStr) {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Parameter '%s' default value '%s' does not match AllowedPattern '%s'", name, defaultStr, param.AllowedPattern),
					Line:    param.Node.Line,
					Column:  param.Node.Column,
					Path:    []string{"Parameters", name, "Default"},
				})
			}
		}

		// Check MinValue/MaxValue for numeric types
		if param.Type == "Number" {
			if numVal, ok := toFloat64(param.Default); ok {
				if param.MinValue != nil && numVal < *param.MinValue {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Parameter '%s' default value %v is less than MinValue %v", name, numVal, *param.MinValue),
						Line:    param.Node.Line,
						Column:  param.Node.Column,
						Path:    []string{"Parameters", name, "Default"},
					})
				}
				if param.MaxValue != nil && numVal > *param.MaxValue {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Parameter '%s' default value %v is greater than MaxValue %v", name, numVal, *param.MaxValue),
						Line:    param.Node.Line,
						Column:  param.Node.Column,
						Path:    []string{"Parameters", name, "Default"},
					})
				}
			}
		}

		// Check MinLength/MaxLength for string types
		if param.Type == "String" {
			defaultStr := fmt.Sprintf("%v", param.Default)
			strLen := len(defaultStr)
			if param.MinLength != nil && strLen < *param.MinLength {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Parameter '%s' default value length %d is less than MinLength %d", name, strLen, *param.MinLength),
					Line:    param.Node.Line,
					Column:  param.Node.Column,
					Path:    []string{"Parameters", name, "Default"},
				})
			}
			if param.MaxLength != nil && strLen > *param.MaxLength {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Parameter '%s' default value length %d is greater than MaxLength %d", name, strLen, *param.MaxLength),
					Line:    param.Node.Line,
					Column:  param.Node.Column,
					Path:    []string{"Parameters", name, "Default"},
				})
			}
		}
	}

	return matches
}

func toFloat64(v any) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case float64:
		return n, true
	case float32:
		return float64(n), true
	}
	return 0, false
}
