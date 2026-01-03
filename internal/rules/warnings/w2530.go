package warnings

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W2530{})
}

// W2530 warns about SnapStart configuration issues for Lambda functions.
type W2530 struct{}

func (r *W2530) ID() string { return "W2530" }

func (r *W2530) ShortDesc() string {
	return "SnapStart configuration"
}

func (r *W2530) Description() string {
	return "Warns when SnapStart is configured for Lambda functions with unsupported runtimes or configurations."
}

func (r *W2530) Source() string {
	return "https://docs.aws.amazon.com/lambda/latest/dg/snapstart.html"
}

func (r *W2530) Tags() []string {
	return []string{"warnings", "lambda", "snapstart", "performance"}
}

// Runtimes that support SnapStart
var snapStartSupportedRuntimes = map[string]bool{
	"java11":     true,
	"java17":     true,
	"java21":     true,
	"java11.al2": true,
	"java17.al2": true,
	"java21.al2": true,
}

func (r *W2530) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Lambda::Function" {
			continue
		}

		// Check if SnapStart is configured
		snapStart, hasSnapStart := res.Properties["SnapStart"]
		if !hasSnapStart {
			// Check if Java runtime without SnapStart - suggest enabling it
			if runtime, ok := res.Properties["Runtime"].(string); ok {
				if snapStartSupportedRuntimes[strings.ToLower(runtime)] {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Lambda function '%s' uses %s runtime; consider enabling SnapStart for improved cold start performance", resName, runtime),
						Path:    []string{"Resources", resName, "Properties"},
					})
				}
			}
			continue
		}

		// SnapStart is configured - validate it
		snapStartMap, ok := snapStart.(map[string]any)
		if !ok {
			continue
		}

		applyOn, _ := snapStartMap["ApplyOn"].(string)
		if applyOn != "PublishedVersions" && applyOn != "None" && applyOn != "" {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Lambda function '%s' has invalid SnapStart ApplyOn value '%s'; must be 'PublishedVersions' or 'None'", resName, applyOn),
				Path:    []string{"Resources", resName, "Properties", "SnapStart", "ApplyOn"},
			})
		}

		// Check runtime compatibility
		runtime, hasRuntime := res.Properties["Runtime"].(string)
		if hasRuntime && applyOn == "PublishedVersions" {
			if !snapStartSupportedRuntimes[strings.ToLower(runtime)] {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Lambda function '%s' has SnapStart enabled but runtime '%s' is not supported; SnapStart requires Java 11 or later", resName, runtime),
					Path:    []string{"Resources", resName, "Properties", "SnapStart"},
				})
			}
		}

		// Check for ephemeral storage warning
		if ephemeralStorage, ok := res.Properties["EphemeralStorage"].(map[string]any); ok {
			if size, ok := ephemeralStorage["Size"].(float64); ok && size > 512 {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Lambda function '%s' has SnapStart enabled with custom ephemeral storage size; cached snapshots may be larger", resName),
					Path:    []string{"Resources", resName, "Properties", "EphemeralStorage"},
				})
			}
		}

		// Check for VPC configuration warning
		if _, hasVpc := res.Properties["VpcConfig"]; hasVpc && applyOn == "PublishedVersions" {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Lambda function '%s' has SnapStart enabled with VPC configuration; first invocation after snapshot may still have network initialization latency", resName),
				Path:    []string{"Resources", resName, "Properties", "VpcConfig"},
			})
		}
	}

	return matches
}
