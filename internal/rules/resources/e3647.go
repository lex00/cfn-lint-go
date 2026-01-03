// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3647{})
}

// E3647 validates ElastiCache node type.
type E3647 struct{}

func (r *E3647) ID() string { return "E3647" }

func (r *E3647) ShortDesc() string {
	return "Validate ElastiCache node type"
}

func (r *E3647) Description() string {
	return "Validates that AWS::ElastiCache::CacheCluster and ReplicationGroup resources specify valid node types."
}

func (r *E3647) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3647"
}

func (r *E3647) Tags() []string {
	return []string{"resources", "properties", "elasticache", "nodetype"}
}

func (r *E3647) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// ElastiCache node type families
	validFamilies := []string{
		"cache.t2.", "cache.t3.", "cache.t4g.",
		"cache.m3.", "cache.m4.", "cache.m5.", "cache.m6g.", "cache.m7g.",
		"cache.r3.", "cache.r4.", "cache.r5.", "cache.r6g.", "cache.r7g.",
		"cache.c1.",
	}

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::ElastiCache::CacheCluster" && res.Type != "AWS::ElastiCache::ReplicationGroup" {
			continue
		}

		nodeType, hasNodeType := res.Properties["CacheNodeType"]
		if !hasNodeType || isIntrinsicFunction(nodeType) {
			continue
		}

		nodeTypeStr, ok := nodeType.(string)
		if !ok {
			continue
		}

		// Check if node type starts with a valid family
		isValid := false
		for _, family := range validFamilies {
			if strings.HasPrefix(nodeTypeStr, family) {
				isValid = true
				break
			}
		}

		if !isValid {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf(
					"Resource '%s': Invalid ElastiCache node type '%s'. Must start with 'cache.'",
					resName, nodeTypeStr,
				),
				Line:   res.Node.Line,
				Column: res.Node.Column,
				Path:   []string{"Resources", resName, "Properties", "CacheNodeType"},
			})
		}
	}

	return matches
}
