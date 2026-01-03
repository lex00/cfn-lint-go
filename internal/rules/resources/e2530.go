package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E2530{})
}

// E2530 validates that SnapStart is only used with compatible Lambda runtimes.
type E2530 struct{}

func (r *E2530) ID() string { return "E2530" }

func (r *E2530) ShortDesc() string {
	return "Validate SnapStart runtime support"
}

func (r *E2530) Description() string {
	return "AWS Lambda SnapStart is only supported for Java 11 and newer runtimes. This rule validates that SnapStart is not configured for incompatible runtimes."
}

func (r *E2530) Source() string {
	return "https://docs.aws.amazon.com/lambda/latest/dg/snapstart.html"
}

func (r *E2530) Tags() []string {
	return []string{"lambda", "snapstart", "runtime"}
}

// snapStartSupportedRuntimes lists Lambda runtimes that support SnapStart
var snapStartSupportedRuntimes = map[string]bool{
	"java11":     true,
	"java17":     true,
	"java21":     true,
	"java11.al2": true,
	"java17.al2": true,
	"java21.al2": true,
}

func (r *E2530) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Lambda::Function" {
			continue
		}

		// Check if SnapStart is configured
		snapStart, hasSnapStart := res.Properties["SnapStart"]
		if !hasSnapStart {
			continue
		}

		// Validate SnapStart is a map
		snapStartMap, ok := snapStart.(map[string]any)
		if !ok {
			// If it's an intrinsic function, we can't validate
			continue
		}

		// Check if ApplyOn is set to PublishedVersions (enabled)
		applyOn, hasApplyOn := snapStartMap["ApplyOn"]
		if !hasApplyOn {
			continue
		}

		applyOnStr, ok := applyOn.(string)
		if !ok {
			// If it's an intrinsic function, skip validation
			continue
		}

		// Only validate if SnapStart is actually enabled
		if applyOnStr != "PublishedVersions" {
			continue
		}

		// Get the Runtime property
		runtime, ok := res.Properties["Runtime"]
		if !ok {
			// If no runtime specified, Lambda will use a default, but SnapStart won't work
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Lambda function '%s' has SnapStart enabled but no Runtime specified. SnapStart requires Java 11 or newer.", resName),
				Line:    res.Node.Line,
				Column:  res.Node.Column,
				Path:    []string{"Resources", resName, "Properties", "SnapStart"},
			})
			continue
		}

		runtimeStr, ok := runtime.(string)
		if !ok {
			// Runtime might be a Ref or other intrinsic function, skip validation
			continue
		}

		// Validate runtime supports SnapStart
		if !isSnapStartSupported(runtimeStr) {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Lambda function '%s' has SnapStart enabled with runtime '%s', but SnapStart is only supported for Java 11 and newer runtimes (java11, java17, java21)", resName, runtimeStr),
				Line:    res.Node.Line,
				Column:  res.Node.Column,
				Path:    []string{"Resources", resName, "Properties", "SnapStart"},
			})
		}
	}

	return matches
}

func isSnapStartSupported(runtime string) bool {
	// Normalize runtime string
	runtime = strings.ToLower(runtime)

	// Check exact match first
	if supported, exists := snapStartSupportedRuntimes[runtime]; exists {
		return supported
	}

	// Check if it's a newer Java runtime (e.g., java22, java23, etc.)
	if strings.HasPrefix(runtime, "java") {
		// Exclude java8 variants
		if strings.Contains(runtime, "java8") || strings.Contains(runtime, "8") {
			return false
		}
		// Assume newer Java versions support SnapStart
		return true
	}

	return false
}
