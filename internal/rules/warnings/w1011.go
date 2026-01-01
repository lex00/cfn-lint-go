// Package warnings contains warning-level rules (Wxxx).
package warnings

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W1011{})
}

// W1011 warns when secrets appear to be hardcoded instead of using dynamic references.
type W1011 struct{}

func (r *W1011) ID() string { return "W1011" }

func (r *W1011) ShortDesc() string {
	return "Use dynamic references for secrets"
}

func (r *W1011) Description() string {
	return "Warns when secret-like properties contain static values instead of dynamic references ({{resolve:secretsmanager:...}} or {{resolve:ssm-secure:...}})."
}

func (r *W1011) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/dynamic-references.html"
}

func (r *W1011) Tags() []string {
	return []string{"warnings", "security", "secrets", "dynamic-references"}
}

// Properties that typically contain secrets
var secretPropertyNames = map[string]bool{
	"password":           true,
	"secret":             true,
	"secretkey":          true,
	"apikey":             true,
	"accesskey":          true,
	"privatekey":         true,
	"masterpassword":     true,
	"masteruserpassword": true,
	"adminpassword":      true,
	"dbpassword":         true,
	"token":              true,
	"authtoken":          true,
	"credentials":        true,
}

// Pattern to detect dynamic references
var dynamicRefDetectPattern = regexp.MustCompile(`\{\{resolve:(secretsmanager|ssm-secure):`)

func (r *W1011) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		secretProps := findSecretProperties(res.Properties, nil)
		for _, sp := range secretProps {
			// Check if the value is a static string without dynamic reference
			if strVal, ok := sp.value.(string); ok {
				if !dynamicRefDetectPattern.MatchString(strVal) && len(strVal) > 0 {
					// Check it's not a Ref or other intrinsic
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Property '%s' in resource '%s' appears to contain a secret. Consider using dynamic references ({{resolve:secretsmanager:...}} or {{resolve:ssm-secure:...}})", sp.path, resName),
						Path:    append([]string{"Resources", resName, "Properties"}, strings.Split(sp.path, ".")...),
					})
				}
			}
		}
	}

	return matches
}

type secretPropInfo struct {
	path  string
	value any
}

func findSecretProperties(v any, path []string) []secretPropInfo {
	var results []secretPropInfo

	switch val := v.(type) {
	case map[string]any:
		// Skip if this is an intrinsic function
		for key := range val {
			if strings.HasPrefix(key, "Fn::") || key == "Ref" {
				return results
			}
		}

		for key, child := range val {
			newPath := append(path, key)
			lowerKey := strings.ToLower(key)
			if secretPropertyNames[lowerKey] {
				results = append(results, secretPropInfo{
					path:  strings.Join(newPath, "."),
					value: child,
				})
			}
			results = append(results, findSecretProperties(child, newPath)...)
		}
	case []any:
		for i, child := range val {
			newPath := append(path, fmt.Sprintf("[%d]", i))
			results = append(results, findSecretProperties(child, newPath)...)
		}
	}

	return results
}
