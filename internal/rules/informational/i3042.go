package informational

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&I3042{})
}

// I3042 suggests using pseudo parameters in ARNs instead of hardcoded values.
type I3042 struct{}

func (r *I3042) ID() string { return "I3042" }

func (r *I3042) ShortDesc() string {
	return "ARN should use pseudo parameters"
}

func (r *I3042) Description() string {
	return "Suggests using AWS pseudo parameters (AWS::AccountId, AWS::Region, AWS::Partition) in ARNs instead of hardcoded values for better portability."
}

func (r *I3042) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/pseudo-parameter-reference.html"
}

func (r *I3042) Tags() []string {
	return []string{"resources", "best-practice", "arn", "portability"}
}

func (r *I3042) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check all resources
	for resName, res := range tmpl.Resources {
		checkForHardcodedARNs(res.Properties, []string{"Resources", resName, "Properties"}, &matches)
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		if out.Value != nil {
			checkForHardcodedARNs(map[string]any{"Value": out.Value}, []string{"Outputs", outName, "Value"}, &matches)
		}
	}

	return matches
}

func checkForHardcodedARNs(v any, path []string, matches *[]rules.Match) {
	switch val := v.(type) {
	case string:
		// Check if it's an ARN with hardcoded values
		if isHardcodedARN(val) {
			suggestion := suggestPseudoParams(val)
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf("ARN contains hardcoded values. Consider using pseudo parameters for better portability. Suggestion: %s", suggestion),
				Path:    path,
			})
		}
	case map[string]any:
		for key, child := range val {
			checkForHardcodedARNs(child, append(path, key), matches)
		}
	case []any:
		for i, child := range val {
			checkForHardcodedARNs(child, append(path, fmt.Sprintf("[%d]", i)), matches)
		}
	}
}

func isHardcodedARN(s string) bool {
	// Check if it's an ARN format
	if !strings.HasPrefix(s, "arn:") {
		return false
	}

	// Check for hardcoded account ID (12 digit number)
	if strings.Contains(s, ":123456789012:") || containsAccountIDPattern(s) {
		return true
	}

	// Check for hardcoded region (e.g., us-east-1, eu-west-1)
	hardcodedRegions := []string{
		":us-east-1:", ":us-east-2:", ":us-west-1:", ":us-west-2:",
		":eu-west-1:", ":eu-west-2:", ":eu-west-3:", ":eu-central-1:",
		":ap-northeast-1:", ":ap-northeast-2:", ":ap-southeast-1:", ":ap-southeast-2:",
		":ap-south-1:", ":ca-central-1:", ":sa-east-1:",
	}
	for _, region := range hardcodedRegions {
		if strings.Contains(s, region) {
			return true
		}
	}

	return false
}

func containsAccountIDPattern(s string) bool {
	// Simple check for 12-digit account ID pattern
	parts := strings.Split(s, ":")
	for _, part := range parts {
		if len(part) == 12 && isDigits(part) {
			return true
		}
	}
	return false
}

func isDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

func suggestPseudoParams(arn string) string {
	// Replace hardcoded account ID with pseudo parameter
	parts := strings.Split(arn, ":")
	if len(parts) >= 5 {
		// ARN format: arn:partition:service:region:account-id:...
		suggestion := "arn"

		// Partition (index 1)
		if parts[1] == "aws" || parts[1] == "aws-cn" || parts[1] == "aws-us-gov" {
			suggestion += ":${AWS::Partition}"
		} else {
			suggestion += ":" + parts[1]
		}

		// Service (index 2)
		suggestion += ":" + parts[2]

		// Region (index 3)
		if parts[3] != "" && !strings.HasPrefix(parts[3], "$") {
			suggestion += ":${AWS::Region}"
		} else {
			suggestion += ":" + parts[3]
		}

		// Account ID (index 4)
		if isDigits(parts[4]) {
			suggestion += ":${AWS::AccountId}"
		} else {
			suggestion += ":" + parts[4]
		}

		// Resource (index 5+)
		for i := 5; i < len(parts); i++ {
			suggestion += ":" + parts[i]
		}

		return suggestion
	}

	return "Use ${AWS::AccountId}, ${AWS::Region}, ${AWS::Partition}"
}
