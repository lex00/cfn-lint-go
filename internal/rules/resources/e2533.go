package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E2533{})
}

// E2533 validates Lambda runtime update compatibility.
// Changing the Runtime property requires replacement of the function in some cases.
type E2533 struct{}

func (r *E2533) ID() string { return "E2533" }

func (r *E2533) ShortDesc() string {
	return "Validate Lambda runtime updatability"
}

func (r *E2533) Description() string {
	return "AWS Lambda allows updating the Runtime property, but changes between incompatible runtime families (e.g., Python to Node.js) or major version changes may require function replacement. This rule validates runtime update compatibility."
}

func (r *E2533) Source() string {
	return "https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html"
}

func (r *E2533) Tags() []string {
	return []string{"lambda", "runtime", "updates"}
}

func (r *E2533) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Lambda::Function" {
			continue
		}

		// Get the Runtime property
		runtime, hasRuntime := res.Properties["Runtime"]
		if !hasRuntime {
			continue
		}

		runtimeStr, ok := runtime.(string)
		if !ok {
			// Runtime might be a Ref or other intrinsic function, skip validation
			continue
		}

		// Check for PackageType property
		packageType, hasPackageType := res.Properties["PackageType"]
		if hasPackageType {
			packageTypeStr, ok := packageType.(string)
			if ok && packageTypeStr == "Image" && runtimeStr != "" {
				// Runtime should not be specified for Image-based functions
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Lambda function '%s' has PackageType 'Image' but also specifies Runtime '%s'. Runtime should not be specified for container image functions.", resName, runtimeStr),
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", resName, "Properties", "Runtime"},
				})
			}
		}

		// Validate runtime is recognized
		if !isValidRuntime(runtimeStr) {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Lambda function '%s' specifies unrecognized runtime '%s'. This may cause deployment issues.", resName, runtimeStr),
				Line:    res.Node.Line,
				Column:  res.Node.Column,
				Path:    []string{"Resources", resName, "Properties", "Runtime"},
			})
		}
	}

	return matches
}

func isCustomRuntime(runtime string) bool {
	runtime = strings.ToLower(runtime)
	return strings.HasPrefix(runtime, "provided")
}

func isValidRuntime(runtime string) bool {
	runtime = strings.ToLower(runtime)

	// List of known valid runtimes (current and deprecated)
	validRuntimes := map[string]bool{
		// Python
		"python2.7":  true,
		"python3.6":  true,
		"python3.7":  true,
		"python3.8":  true,
		"python3.9":  true,
		"python3.10": true,
		"python3.11": true,
		"python3.12": true,
		"python3.13": true,

		// Node.js
		"nodejs":     true,
		"nodejs4.3":  true,
		"nodejs6.10": true,
		"nodejs8.10": true,
		"nodejs10.x": true,
		"nodejs12.x": true,
		"nodejs14.x": true,
		"nodejs16.x": true,
		"nodejs18.x": true,
		"nodejs20.x": true,
		"nodejs22.x": true,

		// Java
		"java8":      true,
		"java8.al2":  true,
		"java11":     true,
		"java17":     true,
		"java21":     true,
		"java11.al2": true,
		"java17.al2": true,
		"java21.al2": true,

		// .NET
		"dotnet6":       true,
		"dotnet8":       true,
		"dotnetcore1.0": true,
		"dotnetcore2.0": true,
		"dotnetcore2.1": true,
		"dotnetcore3.1": true,
		"dotnet5.0":     true,

		// Ruby
		"ruby2.5": true,
		"ruby2.7": true,
		"ruby3.2": true,
		"ruby3.3": true,

		// Go
		"go1.x": true,

		// Custom
		"provided":        true,
		"provided.al2":    true,
		"provided.al2023": true,
	}

	return validRuntimes[runtime]
}
