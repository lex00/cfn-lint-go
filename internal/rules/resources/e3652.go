// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3652{})
}

// E3652 validates Elasticsearch domain instance types.
type E3652 struct{}

func (r *E3652) ID() string { return "E3652" }

func (r *E3652) ShortDesc() string {
	return "Validate Elasticsearch domain instance types"
}

func (r *E3652) Description() string {
	return "Validates that AWS::Elasticsearch::Domain resources specify valid instance types."
}

func (r *E3652) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3652"
}

func (r *E3652) Tags() []string {
	return []string{"resources", "properties", "elasticsearch", "instancetype"}
}

func (r *E3652) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Elasticsearch::Domain" && res.Type != "AWS::OpenSearchService::Domain" {
			continue
		}

		clusterConfig, hasClusterConfig := res.Properties["ElasticsearchClusterConfig"]
		if !hasClusterConfig || isIntrinsicFunction(clusterConfig) {
			continue
		}

		clusterConfigMap, ok := clusterConfig.(map[string]any)
		if !ok {
			continue
		}

		instanceType, hasInstanceType := clusterConfigMap["InstanceType"]
		if !hasInstanceType || isIntrinsicFunction(instanceType) {
			continue
		}

		instanceTypeStr, ok := instanceType.(string)
		if !ok {
			continue
		}

		// Elasticsearch/OpenSearch instance types must end with .elasticsearch
		if !strings.HasSuffix(instanceTypeStr, ".elasticsearch") {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Invalid Elasticsearch instance type '%s'. Must end with '.elasticsearch'",
					resName, instanceTypeStr,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties", "ElasticsearchClusterConfig", "InstanceType"},
			})
		}
	}

	return matches
}
