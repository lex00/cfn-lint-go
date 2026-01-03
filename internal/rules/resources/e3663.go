// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3663{})
}

// E3663 validates Lambda environment variable names.
type E3663 struct{}

func (r *E3663) ID() string { return "E3663" }

func (r *E3663) ShortDesc() string {
	return "Validate Lambda environment variable names"
}

func (r *E3663) Description() string {
	return "Validates that AWS::Lambda::Function environment variables do not use reserved names."
}

func (r *E3663) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3663"
}

func (r *E3663) Tags() []string {
	return []string{"resources", "properties", "lambda", "environment"}
}

func (r *E3663) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Reserved Lambda environment variable names
	reservedNames := map[string]bool{
		"AWS_ACCESS_KEY":                true,
		"AWS_ACCESS_KEY_ID":             true,
		"AWS_DEFAULT_REGION":            true,
		"AWS_EXECUTION_ENV":             true,
		"AWS_LAMBDA_FUNCTION_MEMORY_SIZE": true,
		"AWS_LAMBDA_FUNCTION_NAME":      true,
		"AWS_LAMBDA_FUNCTION_VERSION":   true,
		"AWS_LAMBDA_LOG_GROUP_NAME":     true,
		"AWS_LAMBDA_LOG_STREAM_NAME":    true,
		"AWS_REGION":                    true,
		"AWS_SECRET_ACCESS_KEY":         true,
		"AWS_SECRET_KEY":                true,
		"AWS_SECURITY_TOKEN":            true,
		"AWS_SESSION_TOKEN":             true,
		"LAMBDA_RUNTIME_DIR":            true,
		"LAMBDA_TASK_ROOT":              true,
		"TZ":                            true,
		"_HANDLER":                      true,
		"_X_AMZN_TRACE_ID":              true,
	}

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Lambda::Function" && res.Type != "AWS::Serverless::Function" {
			continue
		}

		environment, hasEnvironment := res.Properties["Environment"]
		if !hasEnvironment || isIntrinsicFunction(environment) {
			continue
		}

		envMap, ok := environment.(map[string]any)
		if !ok {
			continue
		}

		variables, hasVariables := envMap["Variables"]
		if !hasVariables || isIntrinsicFunction(variables) {
			continue
		}

		varsMap, ok := variables.(map[string]any)
		if !ok {
			continue
		}

		for varName := range varsMap {
			if reservedNames[varName] {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': Lambda environment variable '%s' is a reserved name and cannot be used",
						resName, varName,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", "Environment", "Variables", varName},
				})
			}
		}
	}

	return matches
}
