// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3050{})
}

// E3050 validates IAM resource references when Path is specified.
type E3050 struct{}

func (r *E3050) ID() string { return "E3050" }

func (r *E3050) ShortDesc() string {
	return "IAM resource path in Ref"
}

func (r *E3050) Description() string {
	return "Validates that IAM resources with custom Path are not referenced by name using Ref, as CloudFormation cannot lookup resources by name when Path is set."
}

func (r *E3050) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3050"
}

func (r *E3050) Tags() []string {
	return []string{"resources", "properties", "iam", "ref"}
}

func (r *E3050) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Track IAM resources with custom paths
	resourcesWithPath := make(map[string]bool)

	// IAM resource types that support Path
	iamTypes := map[string]bool{
		"AWS::IAM::User":  true,
		"AWS::IAM::Group": true,
		"AWS::IAM::Role":  true,
	}

	// First pass: identify resources with custom Path
	for resName, res := range tmpl.Resources {
		if !iamTypes[res.Type] {
			continue
		}

		if path, hasPath := res.Properties["Path"]; hasPath {
			// Check if Path is not the default "/"
			if pathStr, ok := path.(string); ok && pathStr != "/" {
				resourcesWithPath[resName] = true
			}
		}
	}

	// Second pass: check for Ref to resources with custom paths
	for resName, res := range tmpl.Resources {
		// Check all properties for Ref
		r.checkForRef(res.Properties, resName, resourcesWithPath, &matches, res.Node.Line, res.Node.Column, []string{"Resources", resName, "Properties"})
	}

	return matches
}

func (r *E3050) checkForRef(value interface{}, currentResource string, resourcesWithPath map[string]bool, matches *[]rules.Match, line, column int, path []string) {
	switch v := value.(type) {
	case map[string]interface{}:
		// Check for Ref
		if ref, hasRef := v["Ref"]; hasRef {
			if refStr, ok := ref.(string); ok {
				if resourcesWithPath[refStr] {
					*matches = append(*matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': Cannot use Ref to IAM resource '%s' which has a custom Path. Use GetAtt to retrieve the ARN instead",
							currentResource, refStr,
						),
						Line:   line,
						Column: column,
						Path:   path,
					})
				}
			}
		} else {
			// Recursively check nested maps
			for k, val := range v {
				r.checkForRef(val, currentResource, resourcesWithPath, matches, line, column, append(path, k))
			}
		}
	case []interface{}:
		// Recursively check arrays
		for i, val := range v {
			r.checkForRef(val, currentResource, resourcesWithPath, matches, line, column, append(path, fmt.Sprintf("[%d]", i)))
		}
	}
}
