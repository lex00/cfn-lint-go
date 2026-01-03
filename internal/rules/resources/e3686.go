// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3686{})
}

// E3686 validates Serverless RDS properties.
type E3686 struct{}

func (r *E3686) ID() string { return "E3686" }

func (r *E3686) ShortDesc() string {
	return "Validate Serverless RDS properties"
}

func (r *E3686) Description() string {
	return "Validates that Aurora Serverless DB clusters do not specify incompatible properties."
}

func (r *E3686) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3686"
}

func (r *E3686) Tags() []string {
	return []string{"resources", "properties", "rds", "serverless"}
}

func (r *E3686) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::RDS::DBCluster" {
			continue
		}

		engineMode, hasEngineMode := res.Properties["EngineMode"]
		if !hasEngineMode || isIntrinsicFunction(engineMode) {
			continue
		}

		engineModeStr, ok := engineMode.(string)
		if !ok {
			continue
		}

		if strings.ToLower(engineModeStr) == "serverless" {
			// Serverless clusters cannot have these properties
			invalidProps := []string{"MasterUsername", "MasterUserPassword"}
			for _, prop := range invalidProps {
				if _, hasProp := res.Properties[prop]; hasProp {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': Aurora Serverless DB Cluster should not specify %s in serverless mode",
							resName, prop,
						),
						Line:   res.Node.Line,
						Column: res.Node.Column,
						Path:   []string{"Resources", resName, "Properties", prop},
					})
				}
			}
		}
	}

	return matches
}
