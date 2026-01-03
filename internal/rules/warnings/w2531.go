package warnings

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W2531{})
}

// W2531 warns about Lambda functions using deprecated or EOL runtimes.
type W2531 struct{}

func (r *W2531) ID() string { return "W2531" }

func (r *W2531) ShortDesc() string {
	return "Lambda EOL runtime warning"
}

func (r *W2531) Description() string {
	return "Warns when Lambda functions use deprecated or end-of-life runtimes that may lose support."
}

func (r *W2531) Source() string {
	return "https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html"
}

func (r *W2531) Tags() []string {
	return []string{"warnings", "lambda", "runtime", "deprecation"}
}

// Deprecated/EOL runtimes and their suggested replacements
var deprecatedRuntimes = map[string]string{
	"python2.7":     "python3.12",
	"python3.6":     "python3.12",
	"python3.7":     "python3.12",
	"python3.8":     "python3.12",
	"nodejs":        "nodejs20.x",
	"nodejs4.3":     "nodejs20.x",
	"nodejs6.10":    "nodejs20.x",
	"nodejs8.10":    "nodejs20.x",
	"nodejs10.x":    "nodejs20.x",
	"nodejs12.x":    "nodejs20.x",
	"nodejs14.x":    "nodejs20.x",
	"nodejs16.x":    "nodejs20.x",
	"dotnetcore1.0": "dotnet8",
	"dotnetcore2.0": "dotnet8",
	"dotnetcore2.1": "dotnet8",
	"dotnetcore3.1": "dotnet8",
	"dotnet5.0":     "dotnet8",
	"dotnet6":       "dotnet8",
	"ruby2.5":       "ruby3.3",
	"ruby2.7":       "ruby3.3",
	"java8":         "java21",
	"go1.x":         "provided.al2023",
}

// Runtimes nearing EOL (warning level)
var nearingEOLRuntimes = map[string]string{
	"python3.9":  "python3.12",
	"nodejs18.x": "nodejs20.x",
	"java11":     "java21",
	"ruby3.2":    "ruby3.3",
}

func (r *W2531) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Lambda::Function" {
			continue
		}

		runtime, ok := res.Properties["Runtime"].(string)
		if !ok {
			continue
		}

		runtimeLower := strings.ToLower(runtime)

		// Check for deprecated/EOL runtimes
		if replacement, deprecated := deprecatedRuntimes[runtimeLower]; deprecated {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Lambda function '%s' uses deprecated runtime '%s'; consider upgrading to '%s'", resName, runtime, replacement),
				Path:    []string{"Resources", resName, "Properties", "Runtime"},
			})
			continue
		}

		// Check for runtimes nearing EOL
		if replacement, nearingEOL := nearingEOLRuntimes[runtimeLower]; nearingEOL {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Lambda function '%s' uses runtime '%s' which is approaching end-of-life; consider upgrading to '%s'", resName, runtime, replacement),
				Path:    []string{"Resources", resName, "Properties", "Runtime"},
			})
		}
	}

	return matches
}
