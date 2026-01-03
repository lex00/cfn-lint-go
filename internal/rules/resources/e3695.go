// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3695{})
}

// E3695 validates ElastiCache engine and version compatibility.
type E3695 struct{}

func (r *E3695) ID() string { return "E3695" }

func (r *E3695) ShortDesc() string {
	return "Validate ElastiCache engine and version compatibility"
}

func (r *E3695) Description() string {
	return "Validates that AWS::ElastiCache::CacheCluster and ReplicationGroup resources have compatible Engine and EngineVersion values."
}

func (r *E3695) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3695"
}

func (r *E3695) Tags() []string {
	return []string{"resources", "properties", "elasticache", "engine"}
}

func (r *E3695) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	validEngines := map[string]bool{
		"memcached": true,
		"redis":     true,
	}

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::ElastiCache::CacheCluster" && res.Type != "AWS::ElastiCache::ReplicationGroup" {
			continue
		}

		engine, hasEngine := res.Properties["Engine"]
		if !hasEngine || isIntrinsicFunction(engine) {
			continue
		}

		engineStr, ok := engine.(string)
		if !ok {
			continue
		}

		if !validEngines[engineStr] {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Invalid ElastiCache engine '%s'. Must be 'redis' or 'memcached'",
					resName, engineStr,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties", "Engine"},
			})
		}
	}

	return matches
}
