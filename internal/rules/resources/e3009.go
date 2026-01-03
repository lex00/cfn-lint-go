package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3009{})
}

// E3009 validates CloudFormation init configuration in resource metadata.
type E3009 struct{}

func (r *E3009) ID() string {
	return "E3009"
}

func (r *E3009) ShortDesc() string {
	return "Check CloudFormation init configuration"
}

func (r *E3009) Description() string {
	return "Ensure items in CloudFormation init adhere to standards"
}

func (r *E3009) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-init.html"
}

func (r *E3009) Tags() []string {
	return []string{"resources", "cloudformation", "init"}
}

func (r *E3009) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Metadata == nil {
			continue
		}

		// Check for AWS::CloudFormation::Init
		initRaw, hasInit := res.Metadata["AWS::CloudFormation::Init"]
		if !hasInit {
			continue
		}

		init, ok := initRaw.(map[string]any)
		if !ok {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Resource '%s' has invalid AWS::CloudFormation::Init metadata (must be an object)", resName),
				Line:    res.Node.Line,
				Column:  res.Node.Column,
				Path:    []string{"Resources", resName, "Metadata", "AWS::CloudFormation::Init"},
			})
			continue
		}

		// Validate config sets or direct configs
		for key, value := range init {
			// ConfigSets is optional and must be an object if present
			if key == "configSets" {
				if _, ok := value.(map[string]any); !ok {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Resource '%s' has invalid configSets in AWS::CloudFormation::Init (must be an object)", resName),
						Line:    res.Node.Line,
						Column:  res.Node.Column,
						Path:    []string{"Resources", resName, "Metadata", "AWS::CloudFormation::Init", "configSets"},
					})
				}
				continue
			}

			// Other keys are config names and must be objects with valid sections
			config, ok := value.(map[string]any)
			if !ok {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Resource '%s' has invalid config '%s' in AWS::CloudFormation::Init (must be an object)", resName, key),
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", resName, "Metadata", "AWS::CloudFormation::Init", key},
				})
				continue
			}

			// Validate config sections (packages, groups, users, sources, files, commands, services)
			validSections := map[string]bool{
				"packages": true,
				"groups":   true,
				"users":    true,
				"sources":  true,
				"files":    true,
				"commands": true,
				"services": true,
			}

			for section := range config {
				if !validSections[section] {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Resource '%s' has invalid section '%s' in config '%s' (must be one of: packages, groups, users, sources, files, commands, services)", resName, section, key),
						Line:    res.Node.Line,
						Column:  res.Node.Column,
						Path:    []string{"Resources", resName, "Metadata", "AWS::CloudFormation::Init", key, section},
					})
				}
			}
		}
	}

	return matches
}
