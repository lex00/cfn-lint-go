// Package functions contains intrinsic function validation rules (E1xxx).
package functions

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1040{})
}

// E1040 checks that GetAtt has valid format.
type E1040 struct{}

func (r *E1040) ID() string { return "E1040" }

func (r *E1040) ShortDesc() string {
	return "GetAtt format error"
}

func (r *E1040) Description() string {
	return "Checks that Fn::GetAtt has valid format: either 'Resource.Attribute' string or [Resource, Attribute] array."
}

func (r *E1040) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getatt.html"
}

func (r *E1040) Tags() []string {
	return []string{"functions", "getatt"}
}

func (r *E1040) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check all resources
	for resName, res := range tmpl.Resources {
		getAtts := findAllGetAttFormats(res.Properties)
		for _, ga := range getAtts {
			if err := r.validateGetAttFormat(ga); err != "" {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("%s in resource '%s'", err, resName),
					Line:    ga.line,
					Column:  ga.column,
					Path:    []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		getAtts := findAllGetAttFormats(out.Value)
		for _, ga := range getAtts {
			if err := r.validateGetAttFormat(ga); err != "" {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("%s in output '%s'", err, outName),
					Line:    ga.line,
					Column:  ga.column,
					Path:    []string{"Outputs", outName, "Value"},
				})
			}
		}
	}

	return matches
}

func (r *E1040) validateGetAttFormat(ga getAttFormatInfo) string {
	switch ga.valueType {
	case "string":
		// Must contain at least one dot to separate resource and attribute
		if !strings.Contains(ga.stringValue, ".") {
			return fmt.Sprintf("Fn::GetAtt string format must be 'Resource.Attribute', got '%s'", ga.stringValue)
		}
		parts := strings.SplitN(ga.stringValue, ".", 2)
		if parts[0] == "" {
			return "Fn::GetAtt resource name cannot be empty"
		}
		if len(parts) < 2 || parts[1] == "" {
			return "Fn::GetAtt attribute name cannot be empty"
		}
	case "array":
		if ga.arrayLen < 2 {
			return fmt.Sprintf("Fn::GetAtt array format must have at least 2 elements [Resource, Attribute], got %d", ga.arrayLen)
		}
		if !ga.firstIsString {
			return "Fn::GetAtt first element (resource name) must be a string"
		}
		if !ga.secondIsString {
			return "Fn::GetAtt second element (attribute name) must be a string"
		}
		if ga.firstValue == "" {
			return "Fn::GetAtt resource name cannot be empty"
		}
		if ga.secondValue == "" {
			return "Fn::GetAtt attribute name cannot be empty"
		}
	default:
		return fmt.Sprintf("Fn::GetAtt must be a string 'Resource.Attribute' or array [Resource, Attribute], got %s", ga.valueType)
	}
	return ""
}

type getAttFormatInfo struct {
	valueType      string
	stringValue    string
	arrayLen       int
	firstIsString  bool
	secondIsString bool
	firstValue     string
	secondValue    string
	line           int
	column         int
}

func findAllGetAttFormats(v any) []getAttFormatInfo {
	var results []getAttFormatInfo
	findGetAttFormatsRecursive(v, &results)
	return results
}

func findGetAttFormatsRecursive(v any, results *[]getAttFormatInfo) {
	switch val := v.(type) {
	case map[string]any:
		if ga, ok := val["Fn::GetAtt"]; ok {
			info := parseGetAttFormat(ga)
			*results = append(*results, info)
		}
		for _, child := range val {
			findGetAttFormatsRecursive(child, results)
		}
	case []any:
		for _, child := range val {
			findGetAttFormatsRecursive(child, results)
		}
	}
}

func parseGetAttFormat(v any) getAttFormatInfo {
	info := getAttFormatInfo{}

	switch val := v.(type) {
	case string:
		info.valueType = "string"
		info.stringValue = val
	case []any:
		info.valueType = "array"
		info.arrayLen = len(val)
		if len(val) >= 1 {
			if s, ok := val[0].(string); ok {
				info.firstIsString = true
				info.firstValue = s
			}
		}
		if len(val) >= 2 {
			if s, ok := val[1].(string); ok {
				info.secondIsString = true
				info.secondValue = s
			}
		}
	default:
		info.valueType = fmt.Sprintf("%T", v)
	}

	return info
}
