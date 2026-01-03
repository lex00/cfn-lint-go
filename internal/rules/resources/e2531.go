package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E2531{})
}

// E2531 validates that Lambda functions do not use deprecated runtimes.
type E2531 struct{}

func (r *E2531) ID() string { return "E2531" }

func (r *E2531) ShortDesc() string {
	return "Validate Lambda runtime is not deprecated"
}

func (r *E2531) Description() string {
	return "AWS Lambda deprecates runtimes when their underlying language version reaches end of life. This rule validates that Lambda functions do not use deprecated runtimes."
}

func (r *E2531) Source() string {
	return "https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html"
}

func (r *E2531) Tags() []string {
	return []string{"lambda", "runtime", "deprecation"}
}

// deprecatedLambdaRuntimes maps deprecated runtimes to their deprecation info
var deprecatedLambdaRuntimes = map[string]string{
	// Python
	"python2.7": "deprecated since July 2021",
	"python3.6": "deprecated since July 2022",

	// Node.js
	"nodejs":     "deprecated since October 2016",
	"nodejs4.3":  "deprecated since April 2018",
	"nodejs6.10": "deprecated since August 2019",
	"nodejs8.10": "deprecated since December 2019",
	"nodejs10.x": "deprecated since July 2021",
	"nodejs12.x": "deprecated since March 2023",
	"nodejs14.x": "deprecated since November 2023",
	"nodejs16.x": "deprecated since March 2024",

	// Ruby
	"ruby2.5": "deprecated since July 2021",
	"ruby2.7": "deprecated since December 2023",

	// Java
	"java8": "deprecated since December 2023 (use java8.al2)",

	// .NET
	"dotnetcore1.0": "deprecated since July 2019",
	"dotnetcore2.0": "deprecated since May 2019",
	"dotnetcore2.1": "deprecated since January 2022",
	"dotnetcore3.1": "deprecated since April 2023",
	"dotnet5.0":     "deprecated since May 2022",

	// Go
	"go1.x": "deprecated since December 2023 (use provided.al2)",
}

func (r *E2531) Match(tmpl *template.Template) []rules.Match {
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
			// Runtime might be a Ref or other intrinsic function, skip validation
			continue
		}

		// Check if runtime is deprecated
		if deprecationInfo, deprecated := isRuntimeDeprecated(runtimeStr); deprecated {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Lambda function '%s' uses deprecated runtime '%s' (%s). Please migrate to a supported runtime.", resName, runtimeStr, deprecationInfo),
				Line:    res.Node.Line,
				Column:  res.Node.Column,
				Path:    []string{"Resources", resName, "Properties", "Runtime"},
			})
		}
	}

	return matches
}

func isRuntimeDeprecated(runtime string) (string, bool) {
	// Normalize runtime string
	runtime = strings.ToLower(runtime)

	// Check exact match
	if deprecationInfo, exists := deprecatedLambdaRuntimes[runtime]; exists {
		return deprecationInfo, true
	}

	return "", false
}
