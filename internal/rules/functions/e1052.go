// Package functions contains intrinsic function validation rules (E1xxx).
package functions

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1052{})
}

// E1052 validates SSM Parameter Store dynamic references.
type E1052 struct{}

func (r *E1052) ID() string { return "E1052" }

func (r *E1052) ShortDesc() string {
	return "SSM Parameter dynamic reference validation"
}

func (r *E1052) Description() string {
	return "Validates that dynamic references to SSM parameters are used in valid locations and have proper format."
}

func (r *E1052) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/dynamic-references-ssm.html"
}

func (r *E1052) Tags() []string {
	return []string{"functions", "dynamic-references", "ssm"}
}

// SSM parameter name pattern - alphanumeric and ._-/
var ssmParameterNamePattern = regexp.MustCompile(`^[a-zA-Z0-9/_.\-]+$`)

func (r *E1052) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check resources
	for resName, res := range tmpl.Resources {
		dynRefs := findAllDynamicRefs(res.Properties)
		for _, dr := range dynRefs {
			if strings.HasPrefix(dr.content, "ssm:") || strings.HasPrefix(dr.content, "ssm-secure:") {
				if err := r.validateSSMRef(dr); err != "" {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("%s in resource '%s'", err, resName),
						Line:    dr.line,
						Column:  dr.column,
						Path:    []string{"Resources", resName, "Properties"},
					})
				}
			}
		}
	}

	// SSM dynamic refs can be used in outputs (unlike secretsmanager)
	for outName, out := range tmpl.Outputs {
		dynRefs := findAllDynamicRefs(out.Value)
		for _, dr := range dynRefs {
			if strings.HasPrefix(dr.content, "ssm:") || strings.HasPrefix(dr.content, "ssm-secure:") {
				if err := r.validateSSMRef(dr); err != "" {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("%s in output '%s'", err, outName),
						Line:    dr.line,
						Column:  dr.column,
						Path:    []string{"Outputs", outName, "Value"},
					})
				}
			}
		}
	}

	return matches
}

func (r *E1052) validateSSMRef(dr dynamicRefInfo) string {
	// Format: ssm:parameter-name:version or ssm-secure:parameter-name:version
	// Minimum: ssm:parameter-name or ssm-secure:parameter-name

	var rest string

	if strings.HasPrefix(dr.content, "ssm-secure:") {
		rest = strings.TrimPrefix(dr.content, "ssm-secure:")
	} else {
		rest = strings.TrimPrefix(dr.content, "ssm:")
	}

	parts := strings.Split(rest, ":")
	if len(parts) < 1 || parts[0] == "" {
		return fmt.Sprintf("SSM dynamic reference '%s' is missing parameter name", dr.fullRef)
	}

	parameterName := parts[0]

	// Check for trailing backslash (not supported)
	if strings.HasSuffix(dr.content, "\\") {
		return fmt.Sprintf("SSM dynamic reference '%s' cannot end with a backslash (\\)", dr.fullRef)
	}

	// Validate parameter name format (alphanumeric and ._-/)
	if !ssmParameterNamePattern.MatchString(parameterName) {
		return fmt.Sprintf("SSM dynamic reference '%s' has invalid parameter name '%s', must contain only alphanumeric characters and ._-/", dr.fullRef, parameterName)
	}

	// If version is specified, validate it's an integer
	if len(parts) > 1 && parts[1] != "" {
		version := parts[1]
		// Check if version is numeric
		if !regexp.MustCompile(`^\d+$`).MatchString(version) {
			return fmt.Sprintf("SSM dynamic reference '%s' has invalid version '%s', must be an integer", dr.fullRef, version)
		}
	}

	return ""
}
