package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W2030{})
}

// W2030 warns when parameter default values might not match allowed values.
type W2030 struct{}

func (r *W2030) ID() string { return "W2030" }

func (r *W2030) ShortDesc() string {
	return "Parameter valid value check"
}

func (r *W2030) Description() string {
	return "Warns when a parameter's constraints (MinValue, MaxValue, MinLength, MaxLength) may be too restrictive or allow invalid values."
}

func (r *W2030) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/parameters-section-structure.html"
}

func (r *W2030) Tags() []string {
	return []string{"warnings", "parameters", "validation"}
}

func (r *W2030) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for paramName, param := range tmpl.Parameters {
		// Check numeric constraints
		if param.Type == "Number" {
			r.checkNumberConstraints(paramName, param, &matches)
		}

		// Check string length constraints
		if param.Type == "String" {
			r.checkStringConstraints(paramName, param, &matches)
		}

		// Check AllowedValues
		if len(param.AllowedValues) > 0 {
			r.checkAllowedValues(paramName, param, &matches)
		}
	}

	return matches
}

func (r *W2030) checkNumberConstraints(paramName string, param *template.Parameter, matches *[]rules.Match) {
	var minVal, maxVal float64
	hasMin, hasMax := false, false

	if param.MinValue != nil {
		minVal = *param.MinValue
		hasMin = true
	}

	if param.MaxValue != nil {
		maxVal = *param.MaxValue
		hasMax = true
	}

	// Check if min > max
	if hasMin && hasMax && minVal > maxVal {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("Parameter '%s' has MinValue (%v) greater than MaxValue (%v)", paramName, minVal, maxVal),
			Path:    []string{"Parameters", paramName},
		})
	}

	// Check if range is very large
	if hasMin && hasMax && maxVal-minVal > 1e9 {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("Parameter '%s' has a very large range (%v to %v); consider if this is intentional", paramName, minVal, maxVal),
			Path:    []string{"Parameters", paramName},
		})
	}
}

func (r *W2030) checkStringConstraints(paramName string, param *template.Parameter, matches *[]rules.Match) {
	var minLen, maxLen int
	hasMin, hasMax := false, false

	if param.MinLength != nil {
		minLen = *param.MinLength
		hasMin = true
	}

	if param.MaxLength != nil {
		maxLen = *param.MaxLength
		hasMax = true
	}

	// Check if min > max
	if hasMin && hasMax && minLen > maxLen {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("Parameter '%s' has MinLength (%d) greater than MaxLength (%d)", paramName, minLen, maxLen),
			Path:    []string{"Parameters", paramName},
		})
	}

	// Check for zero max length
	if hasMax && maxLen == 0 {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("Parameter '%s' has MaxLength of 0; only empty strings are allowed", paramName),
			Path:    []string{"Parameters", paramName},
		})
	}
}

func (r *W2030) checkAllowedValues(paramName string, param *template.Parameter, matches *[]rules.Match) {
	// Check for duplicate allowed values
	seen := make(map[string]bool)
	for _, val := range param.AllowedValues {
		valStr := fmt.Sprintf("%v", val)
		if seen[valStr] {
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf("Parameter '%s' has duplicate value '%s' in AllowedValues", paramName, valStr),
				Path:    []string{"Parameters", paramName, "AllowedValues"},
			})
		}
		seen[valStr] = true
	}

	// Check if only one allowed value (might as well be a default)
	if len(param.AllowedValues) == 1 {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("Parameter '%s' has only one AllowedValue; consider using a default value instead", paramName),
			Path:    []string{"Parameters", paramName, "AllowedValues"},
		})
	}
}
