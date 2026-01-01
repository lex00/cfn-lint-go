// Package functions contains intrinsic function validation rules (E1xxx).
package functions

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1010{})
}

// E1010 checks that GetAtt references point to existing resources.
type E1010 struct{}

func (r *E1010) ID() string { return "E1010" }

func (r *E1010) ShortDesc() string {
	return "GetAtt to undefined resource"
}

func (r *E1010) Description() string {
	return "Checks that all Fn::GetAtt intrinsic functions reference valid resources."
}

func (r *E1010) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getatt.html"
}

func (r *E1010) Tags() []string {
	return []string{"functions", "getatt"}
}

func (r *E1010) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check all resources for invalid GetAtt references
	for resName, res := range tmpl.Resources {
		getAtts := findAllGetAtts(res.Properties)
		for _, ga := range getAtts {
			if !tmpl.HasResource(ga.resource) {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("GetAtt references undefined resource '%s' in resource '%s'", ga.resource, resName),
					Line:    ga.line,
					Column:  ga.column,
					Path:    []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		getAtts := findAllGetAtts(out.Value)
		for _, ga := range getAtts {
			if !tmpl.HasResource(ga.resource) {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("GetAtt references undefined resource '%s' in output '%s'", ga.resource, outName),
					Line:    ga.line,
					Column:  ga.column,
					Path:    []string{"Outputs", outName, "Value"},
				})
			}
		}
	}

	return matches
}

type getAttInfo struct {
	resource  string
	attribute string
	line      int
	column    int
}

func findAllGetAtts(v any) []getAttInfo {
	var results []getAttInfo
	findGetAttsRecursive(v, &results)
	return results
}

func findGetAttsRecursive(v any, results *[]getAttInfo) {
	switch val := v.(type) {
	case map[string]any:
		if ga, ok := val["Fn::GetAtt"]; ok {
			info := parseGetAttValue(ga)
			if info.resource != "" {
				*results = append(*results, info)
			}
		}
		for _, child := range val {
			findGetAttsRecursive(child, results)
		}
	case []any:
		for _, child := range val {
			findGetAttsRecursive(child, results)
		}
	}
}

func parseGetAttValue(v any) getAttInfo {
	switch val := v.(type) {
	case string:
		// "Resource.Attribute" format
		parts := strings.SplitN(val, ".", 2)
		if len(parts) >= 1 {
			info := getAttInfo{resource: parts[0]}
			if len(parts) == 2 {
				info.attribute = parts[1]
			}
			return info
		}
	case []any:
		// [Resource, Attribute] format
		if len(val) >= 1 {
			if res, ok := val[0].(string); ok {
				info := getAttInfo{resource: res}
				if len(val) >= 2 {
					if attr, ok := val[1].(string); ok {
						info.attribute = attr
					}
				}
				return info
			}
		}
	}
	return getAttInfo{}
}
