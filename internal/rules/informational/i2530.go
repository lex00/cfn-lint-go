package informational

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&I2530{})
}

// I2530 suggests using SnapStart for Lambda functions with Java 11+ runtimes.
type I2530 struct{}

func (r *I2530) ID() string { return "I2530" }

func (r *I2530) ShortDesc() string {
	return "Consider SnapStart for Java11+ runtimes"
}

func (r *I2530) Description() string {
	return "Suggests enabling SnapStart for AWS Lambda functions using Java 11 or newer runtimes to reduce cold start times."
}

func (r *I2530) Source() string {
	return "https://docs.aws.amazon.com/lambda/latest/dg/snapstart.html"
}

func (r *I2530) Tags() []string {
	return []string{"lambda", "performance", "snapstart", "java"}
}

// snapStartEligibleRuntimes lists Lambda runtimes that support SnapStart
var snapStartEligibleRuntimes = map[string]bool{
	"java11":          true,
	"java17":          true,
	"java21":          true,
	"java11.al2":      true,
	"java17.al2":      true,
	"java21.al2":      true,
	"java8.al2":       false, // Java 8 doesn't support SnapStart
	"provided.al2":    false,
	"provided.al2023": false,
	"python3.9":       false,
	"python3.10":      false,
	"python3.11":      false,
	"python3.12":      false,
	"nodejs18.x":      false,
	"nodejs20.x":      false,
}

func (r *I2530) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Lambda::Function" {
			continue
		}

		// Get the Runtime property
		runtime, ok := res.Properties["Runtime"]
		if !ok {
			continue
		}

		runtimeStr, ok := runtime.(string)
		if !ok {
			// Runtime might be a Ref or other intrinsic function
			continue
		}

		// Check if runtime supports SnapStart
		if isSnapStartEligible(runtimeStr) {
			// Check if SnapStart is already configured
			snapStart, hasSnapStart := res.Properties["SnapStart"]
			if !hasSnapStart || !isSnapStartEnabled(snapStart) {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Lambda function '%s' uses runtime '%s' which supports SnapStart. Consider enabling SnapStart to reduce cold start times.", resName, runtimeStr),
					Path:    []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	return matches
}

func isSnapStartEligible(runtime string) bool {
	// Normalize runtime string (remove version suffixes if needed)
	runtime = strings.ToLower(runtime)

	// Check exact match first
	if eligible, exists := snapStartEligibleRuntimes[runtime]; exists {
		return eligible
	}

	// Check if it's a Java runtime (java11, java17, java21)
	if strings.HasPrefix(runtime, "java") {
		// Exclude java8 variants
		if strings.Contains(runtime, "java8") {
			return false
		}
		// java11, java17, java21 and newer support SnapStart
		return true
	}

	return false
}

func isSnapStartEnabled(snapStart any) bool {
	snapStartMap, ok := snapStart.(map[string]any)
	if !ok {
		return false
	}

	applyOn, ok := snapStartMap["ApplyOn"]
	if !ok {
		return false
	}

	applyOnStr, ok := applyOn.(string)
	if !ok {
		return false
	}

	return applyOnStr == "PublishedVersions"
}
