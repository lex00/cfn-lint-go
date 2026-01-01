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
	rules.Register(&E1019{})
}

// E1019 checks that Sub function has valid syntax and references.
type E1019 struct{}

func (r *E1019) ID() string { return "E1019" }

func (r *E1019) ShortDesc() string {
	return "Sub function validation"
}

func (r *E1019) Description() string {
	return "Checks that Fn::Sub has valid syntax and all variable references are defined."
}

func (r *E1019) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-sub.html"
}

func (r *E1019) Tags() []string {
	return []string{"functions", "sub"}
}

// Regex to find ${VarName} patterns in Sub strings
var subVarPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// Pseudo-parameters that are always valid
var pseudoParams = map[string]bool{
	"AWS::AccountId":        true,
	"AWS::NotificationARNs": true,
	"AWS::NoValue":          true,
	"AWS::Partition":        true,
	"AWS::Region":           true,
	"AWS::StackId":          true,
	"AWS::StackName":        true,
	"AWS::URLSuffix":        true,
}

func (r *E1019) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Build set of valid references
	validRefs := make(map[string]bool)
	for name := range tmpl.Resources {
		validRefs[name] = true
	}
	for name := range tmpl.Parameters {
		validRefs[name] = true
	}
	for pseudo := range pseudoParams {
		validRefs[pseudo] = true
	}

	// Check all resources
	for resName, res := range tmpl.Resources {
		subs := findAllSubs(res.Properties)
		for _, sub := range subs {
			errs := r.validateSub(sub, validRefs)
			for _, err := range errs {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("%s in resource '%s'", err, resName),
					Line:    sub.line,
					Column:  sub.column,
					Path:    []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		subs := findAllSubs(out.Value)
		for _, sub := range subs {
			errs := r.validateSub(sub, validRefs)
			for _, err := range errs {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("%s in output '%s'", err, outName),
					Line:    sub.line,
					Column:  sub.column,
					Path:    []string{"Outputs", outName, "Value"},
				})
			}
		}
	}

	return matches
}

func (r *E1019) validateSub(sub subInfo, validRefs map[string]bool) []string {
	var errors []string

	// Extract variable names from the template string
	varsInTemplate := extractSubVariables(sub.templateString)

	// Build set of variables defined in the optional map
	definedVars := make(map[string]bool)
	for k := range sub.variableMap {
		definedVars[k] = true
	}

	// Check each variable is either in the map, a valid ref, or a GetAtt
	for _, v := range varsInTemplate {
		// Skip if defined in the variable map
		if definedVars[v] {
			continue
		}

		// Handle GetAtt syntax: Resource.Attribute
		if strings.Contains(v, ".") {
			parts := strings.SplitN(v, ".", 2)
			if !validRefs[parts[0]] && !pseudoParams[parts[0]] {
				errors = append(errors, fmt.Sprintf("Fn::Sub references undefined resource '%s' via GetAtt syntax", parts[0]))
			}
			continue
		}

		// Check if it's a valid reference
		if !validRefs[v] {
			errors = append(errors, fmt.Sprintf("Fn::Sub references undefined variable '%s'", v))
		}
	}

	return errors
}

func extractSubVariables(s string) []string {
	matches := subVarPattern.FindAllStringSubmatch(s, -1)
	var vars []string
	for _, m := range matches {
		if len(m) >= 2 {
			// Skip literal $$ escapes and !Literal markers
			if !strings.HasPrefix(m[1], "!") {
				vars = append(vars, m[1])
			}
		}
	}
	return vars
}

type subInfo struct {
	templateString string
	variableMap    map[string]any
	line           int
	column         int
}

func findAllSubs(v any) []subInfo {
	var results []subInfo
	findSubsRecursive(v, &results)
	return results
}

func findSubsRecursive(v any, results *[]subInfo) {
	switch val := v.(type) {
	case map[string]any:
		if sub, ok := val["Fn::Sub"]; ok {
			info := parseSubValue(sub)
			*results = append(*results, info)
		}
		for _, child := range val {
			findSubsRecursive(child, results)
		}
	case []any:
		for _, child := range val {
			findSubsRecursive(child, results)
		}
	}
}

func parseSubValue(v any) subInfo {
	switch val := v.(type) {
	case string:
		return subInfo{templateString: val}
	case []any:
		info := subInfo{}
		if len(val) >= 1 {
			if s, ok := val[0].(string); ok {
				info.templateString = s
			}
		}
		if len(val) >= 2 {
			if m, ok := val[1].(map[string]any); ok {
				info.variableMap = m
			}
		}
		return info
	}
	return subInfo{}
}
