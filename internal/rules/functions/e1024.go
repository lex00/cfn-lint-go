// Package functions contains intrinsic function validation rules (E1xxx).
package functions

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1024{})
}

// E1024 checks that Fn::Cidr is properly configured.
type E1024 struct{}

func (r *E1024) ID() string { return "E1024" }

func (r *E1024) ShortDesc() string {
	return "Fn::Cidr function error"
}

func (r *E1024) Description() string {
	return "Checks that Fn::Cidr is properly configured with [ipBlock, count, cidrBits]."
}

func (r *E1024) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-cidr.html"
}

func (r *E1024) Tags() []string {
	return []string{"functions", "cidr"}
}

func (r *E1024) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check all resources
	for resName, res := range tmpl.Resources {
		findCidrErrors(res.Properties, []string{"Resources", resName, "Properties"}, &matches)
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		if out.Value != nil {
			findCidrErrors(out.Value, []string{"Outputs", outName, "Value"}, &matches)
		}
	}

	return matches
}

func findCidrErrors(v any, path []string, matches *[]rules.Match) {
	switch val := v.(type) {
	case map[string]any:
		if cidrVal, ok := val["Fn::Cidr"]; ok {
			validateCidrArgs(cidrVal, path, matches)
		}
		for key, child := range val {
			findCidrErrors(child, append(path, key), matches)
		}
	case []any:
		for i, child := range val {
			findCidrErrors(child, append(path, fmt.Sprintf("[%d]", i)), matches)
		}
	}
}

func validateCidrArgs(v any, path []string, matches *[]rules.Match) {
	arr, ok := v.([]any)
	if !ok {
		*matches = append(*matches, rules.Match{
			Message: "Fn::Cidr requires an array of [ipBlock, count, cidrBits]",
			Path:    path,
		})
		return
	}

	if len(arr) != 3 {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("Fn::Cidr requires exactly 3 elements [ipBlock, count, cidrBits], got %d", len(arr)),
			Path:    path,
		})
		return
	}

	// First element must be ipBlock (string or intrinsic)
	if !isValidCidrIPBlock(arr[0]) {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("Fn::Cidr ipBlock must be a string or intrinsic function, got %T", arr[0]),
			Path:    path,
		})
	}

	// Second element must be count (integer or intrinsic)
	if !isValidCidrCount(arr[1]) {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("Fn::Cidr count must be an integer (1-256) or intrinsic function, got %T", arr[1]),
			Path:    path,
		})
	}

	// Third element must be cidrBits (integer or intrinsic)
	if !isValidCidrBits(arr[2]) {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("Fn::Cidr cidrBits must be an integer (1-128) or intrinsic function, got %T", arr[2]),
			Path:    path,
		})
	}
}

func isValidCidrIPBlock(v any) bool {
	switch val := v.(type) {
	case string:
		return true
	case map[string]any:
		return len(val) > 0
	default:
		return false
	}
}

func isValidCidrCount(v any) bool {
	switch val := v.(type) {
	case int, int64, float64:
		return true
	case string:
		return true // String representation of number
	case map[string]any:
		return len(val) > 0
	default:
		return false
	}
}

func isValidCidrBits(v any) bool {
	switch val := v.(type) {
	case int, int64, float64:
		return true
	case string:
		return true // String representation of number
	case map[string]any:
		return len(val) > 0
	default:
		return false
	}
}
