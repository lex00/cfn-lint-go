// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3692{})
}

// E3692 validates Multi-AZ DB cluster config.
type E3692 struct{}

func (r *E3692) ID() string { return "E3692" }

func (r *E3692) ShortDesc() string {
	return "Validate Multi-AZ DB cluster config"
}

func (r *E3692) Description() string {
	return "Validates that Multi-AZ RDS DB clusters are configured correctly."
}

func (r *E3692) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3692"
}

func (r *E3692) Tags() []string {
	return []string{"resources", "properties", "rds", "multiaz"}
}

func (r *E3692) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::RDS::DBCluster" {
			continue
		}

		dbClusterInstanceClass, hasInstanceClass := res.Properties["DBClusterInstanceClass"]
		if !hasInstanceClass {
			continue
		}

		// If DBClusterInstanceClass is specified, this is a Multi-AZ cluster
		// Multi-AZ clusters require AllocatedStorage and Iops
		if !isIntrinsicFunction(dbClusterInstanceClass) {
			_, hasAllocatedStorage := res.Properties["AllocatedStorage"]
			if !hasAllocatedStorage {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': Multi-AZ DB Cluster (with DBClusterInstanceClass) must specify AllocatedStorage",
						resName,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties"},
				})
			}

			storageType, hasStorageType := res.Properties["StorageType"]
			if hasStorageType && !isIntrinsicFunction(storageType) {
				storageTypeStr, ok := storageType.(string)
				if ok && storageTypeStr == "io1" {
					_, hasIops := res.Properties["Iops"]
					if !hasIops {
						matches = append(matches, rules.Match{
							Message: fmt.Sprintf(
								"Resource '%s': Multi-AZ DB Cluster with StorageType 'io1' must specify Iops",
								resName,
							),
							Line:   res.Node.Line,
							Column: res.Node.Column,
							Path:   []string{"Resources", resName, "Properties"},
						})
					}
				}
			}
		}
	}

	return matches
}
