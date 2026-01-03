// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3673{})
}

// E3673 validates ImageId requirement.
type E3673 struct{}

func (r *E3673) ID() string { return "E3673" }

func (r *E3673) ShortDesc() string {
	return "Validate ImageId requirement"
}

func (r *E3673) Description() string {
	return "Validates that AWS::EC2::Instance resources specify ImageId unless using LaunchTemplate."
}

func (r *E3673) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3673"
}

func (r *E3673) Tags() []string {
	return []string{"resources", "properties", "ec2", "imageid"}
}

func (r *E3673) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::EC2::Instance" {
			continue
		}

		_, hasImageId := res.Properties["ImageId"]
		_, hasLaunchTemplate := res.Properties["LaunchTemplate"]

		// ImageId is required unless LaunchTemplate is specified
		if !hasImageId && !hasLaunchTemplate {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': EC2 Instance must specify ImageId or LaunchTemplate",
					resName,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties"},
			})
		}
	}

	return matches
}
