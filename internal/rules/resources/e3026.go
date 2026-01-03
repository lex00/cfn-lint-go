package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3026{})
}

// E3026 validates ElastiCache Redis cluster settings.
type E3026 struct{}

func (r *E3026) ID() string {
	return "E3026"
}

func (r *E3026) ShortDesc() string {
	return "Check Elastic Cache Redis Cluster settings"
}

func (r *E3026) Description() string {
	return "Evaluates automatic failover enablement when cluster mode is active"
}

func (r *E3026) Source() string {
	return "https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/AutoFailover.html"
}

func (r *E3026) Tags() []string {
	return []string{"resources", "elasticache", "redis"}
}

func (r *E3026) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::ElastiCache::ReplicationGroup" {
			continue
		}

		// Check if cluster mode is enabled
		clusterModeEnabled := false
		if cmRaw, ok := res.Properties["ClusterMode"]; ok {
			if !isIntrinsicFunction(cmRaw) {
				if cm, ok := cmRaw.(string); ok && cm == "enabled" {
					clusterModeEnabled = true
				}
			}
		}

		// Check NumNodeGroups (if > 1, cluster mode is implied)
		if numNodeGroupsRaw, ok := res.Properties["NumNodeGroups"]; ok {
			if !isIntrinsicFunction(numNodeGroupsRaw) {
				if numNodeGroups, ok := numNodeGroupsRaw.(int); ok && numNodeGroups > 1 {
					clusterModeEnabled = true
				}
			}
		}

		if clusterModeEnabled {
			// AutomaticFailoverEnabled must be true
			afEnabled := false
			if afRaw, ok := res.Properties["AutomaticFailoverEnabled"]; ok {
				if !isIntrinsicFunction(afRaw) {
					if af, ok := afRaw.(bool); ok && af {
						afEnabled = true
					}
				}
			}

			if !afEnabled {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("ElastiCache ReplicationGroup '%s' with cluster mode enabled must have AutomaticFailoverEnabled set to true", resName),
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	return matches
}
