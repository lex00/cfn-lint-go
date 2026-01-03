// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"encoding/json"
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
	"gopkg.in/yaml.v3"
)

func init() {
	rules.Register(&E3051{})
}

// E3051 validates SSM document structure.
type E3051 struct{}

func (r *E3051) ID() string { return "E3051" }

func (r *E3051) ShortDesc() string {
	return "SSM document structure"
}

func (r *E3051) Description() string {
	return "Validates that AWS::SSM::Document Content property contains valid JSON or YAML structure."
}

func (r *E3051) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3051"
}

func (r *E3051) Tags() []string {
	return []string{"resources", "properties", "ssm", "document"}
}

func (r *E3051) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::SSM::Document" {
			continue
		}

		content, hasContent := res.Properties["Content"]
		if !hasContent {
			continue
		}

		// Content can be a string (JSON/YAML) or a map/object
		switch v := content.(type) {
		case string:
			// Try to parse as JSON first
			var jsonData interface{}
			if err := json.Unmarshal([]byte(v), &jsonData); err != nil {
				// Try to parse as YAML
				var yamlData interface{}
				if err := yaml.Unmarshal([]byte(v), &yamlData); err != nil {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf(
							"Resource '%s': SSM Document Content must be valid JSON or YAML (JSON error: %v, YAML error: %v)",
							resName, err, err,
						),
						Line:   res.Node.Line,
						Column: res.Node.Column,
						Path:   []string{"Resources", resName, "Properties", "Content"},
					})
				}
			}
		case map[string]interface{}:
			// Already a valid structure, validate it has required fields
			if _, hasSchemaVersion := v["schemaVersion"]; !hasSchemaVersion {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': SSM Document Content should include 'schemaVersion'",
						resName,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", "Content"},
				})
			}
		}
	}

	return matches
}
