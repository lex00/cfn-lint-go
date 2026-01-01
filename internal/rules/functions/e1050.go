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
	rules.Register(&E1050{})
}

// E1050 checks that dynamic references have valid syntax.
type E1050 struct{}

func (r *E1050) ID() string { return "E1050" }

func (r *E1050) ShortDesc() string {
	return "Dynamic reference syntax error"
}

func (r *E1050) Description() string {
	return "Checks that dynamic references ({{resolve:...}}) have valid syntax."
}

func (r *E1050) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/dynamic-references.html"
}

func (r *E1050) Tags() []string {
	return []string{"functions", "dynamic-references"}
}

// Pattern to find dynamic references
var dynamicRefSyntaxPattern = regexp.MustCompile(`\{\{resolve:([^}]+)\}\}`)

// Valid dynamic reference services
var validDynamicRefServices = map[string]bool{
	"ssm":            true,
	"ssm-secure":     true,
	"secretsmanager": true,
}

func (r *E1050) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check all resources
	for resName, res := range tmpl.Resources {
		dynRefs := findAllDynamicRefs(res.Properties)
		for _, dr := range dynRefs {
			if err := r.validateDynamicRef(dr); err != "" {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("%s in resource '%s'", err, resName),
					Line:    dr.line,
					Column:  dr.column,
					Path:    []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		dynRefs := findAllDynamicRefs(out.Value)
		for _, dr := range dynRefs {
			if err := r.validateDynamicRef(dr); err != "" {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("%s in output '%s'", err, outName),
					Line:    dr.line,
					Column:  dr.column,
					Path:    []string{"Outputs", outName, "Value"},
				})
			}
		}
	}

	return matches
}

func (r *E1050) validateDynamicRef(dr dynamicRefInfo) string {
	// Parse the dynamic reference
	// Format: {{resolve:service-name:reference-key}}
	// or: {{resolve:ssm:parameter-name:version}}
	// or: {{resolve:secretsmanager:secret-id:secret-string:json-key:version-stage:version-id}}

	parts := strings.Split(dr.content, ":")
	if len(parts) < 2 {
		return fmt.Sprintf("Dynamic reference '%s' is malformed, expected {{resolve:service:reference}}", dr.fullRef)
	}

	service := parts[0]
	if !validDynamicRefServices[service] {
		return fmt.Sprintf("Dynamic reference uses unknown service '%s', valid services are: ssm, ssm-secure, secretsmanager", service)
	}

	// Validate based on service type
	switch service {
	case "ssm", "ssm-secure":
		// Format: ssm:parameter-name or ssm:parameter-name:version
		if len(parts) < 2 || parts[1] == "" {
			return fmt.Sprintf("Dynamic reference '%s' is missing parameter name", dr.fullRef)
		}
	case "secretsmanager":
		// Format: secretsmanager:secret-id:secret-string:json-key:version-stage:version-id
		if len(parts) < 2 || parts[1] == "" {
			return fmt.Sprintf("Dynamic reference '%s' is missing secret ID", dr.fullRef)
		}
	}

	return ""
}

type dynamicRefInfo struct {
	fullRef string
	content string // Content after "resolve:"
	line    int
	column  int
}

func findAllDynamicRefs(v any) []dynamicRefInfo {
	var results []dynamicRefInfo
	findDynamicRefsRecursive(v, &results)
	return results
}

func findDynamicRefsRecursive(v any, results *[]dynamicRefInfo) {
	switch val := v.(type) {
	case string:
		// Find all dynamic references in this string
		matches := dynamicRefSyntaxPattern.FindAllStringSubmatch(val, -1)
		for _, m := range matches {
			if len(m) >= 2 {
				*results = append(*results, dynamicRefInfo{
					fullRef: m[0],
					content: m[1],
				})
			}
		}
	case map[string]any:
		for _, child := range val {
			findDynamicRefsRecursive(child, results)
		}
	case []any:
		for _, child := range val {
			findDynamicRefsRecursive(child, results)
		}
	}
}
