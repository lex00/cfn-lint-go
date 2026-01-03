package informational

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
	"gopkg.in/yaml.v3"
)

func init() {
	rules.Register(&I3011{})
}

// I3011 suggests explicit UpdateReplacePolicy/DeletionPolicy for stateful resources.
type I3011 struct{}

func (r *I3011) ID() string { return "I3011" }

func (r *I3011) ShortDesc() string {
	return "Stateful resources need explicit policies"
}

func (r *I3011) Description() string {
	return "Suggests setting explicit UpdateReplacePolicy and DeletionPolicy for stateful resources (databases, storage, etc.) to prevent accidental data loss."
}

func (r *I3011) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-attribute-deletionpolicy.html"
}

func (r *I3011) Tags() []string {
	return []string{"resources", "best-practice", "data-protection"}
}

// statefulResourceTypes lists resource types that typically contain important data
var statefulResourceTypes = map[string]bool{
	"AWS::RDS::DBInstance":                 true,
	"AWS::RDS::DBCluster":                  true,
	"AWS::DynamoDB::Table":                 true,
	"AWS::S3::Bucket":                      true,
	"AWS::EFS::FileSystem":                 true,
	"AWS::FSx::FileSystem":                 true,
	"AWS::Redshift::Cluster":               true,
	"AWS::ElastiCache::CacheCluster":       true,
	"AWS::ElastiCache::ReplicationGroup":   true,
	"AWS::Neptune::DBCluster":              true,
	"AWS::DocDB::DBCluster":                true,
	"AWS::Backup::BackupVault":             true,
	"AWS::EBS::Volume":                     true,
	"AWS::EC2::Volume":                     true,
	"AWS::Timestream::Database":            true,
	"AWS::Timestream::Table":               true,
	"AWS::QLDB::Ledger":                    true,
	"AWS::MQ::Broker":                      true,
	"AWS::MemoryDB::Cluster":               true,
	"AWS::Cassandra::Keyspace":             true,
	"AWS::Cassandra::Table":                true,
	"AWS::OpenSearchService::Domain":       true,
	"AWS::Elasticsearch::Domain":           true,
	"AWS::Kinesis::Stream":                 true,
	"AWS::KinesisFirehose::DeliveryStream": true,
	"AWS::KinesisAnalyticsV2::Application": true,
	"AWS::CloudWatch::LogGroup":            true,
	"AWS::ECR::Repository":                 true,
	"AWS::SecretsManager::Secret":          true,
	"AWS::SSM::Parameter":                  true,
	"AWS::CloudFormation::Stack":           true,
}

func (r *I3011) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if !statefulResourceTypes[res.Type] {
			continue
		}

		// Check the resource node for DeletionPolicy and UpdateReplacePolicy
		hasDeletionPolicy := false
		hasUpdateReplacePolicy := false

		if res.Node != nil && res.Node.Kind == yaml.MappingNode {
			for i := 0; i < len(res.Node.Content); i += 2 {
				key := res.Node.Content[i]
				switch key.Value {
				case "DeletionPolicy":
					hasDeletionPolicy = true
				case "UpdateReplacePolicy":
					hasUpdateReplacePolicy = true
				}
			}
		}

		if !hasDeletionPolicy || !hasUpdateReplacePolicy {
			missingPolicies := []string{}
			if !hasDeletionPolicy {
				missingPolicies = append(missingPolicies, "DeletionPolicy")
			}
			if !hasUpdateReplacePolicy {
				missingPolicies = append(missingPolicies, "UpdateReplacePolicy")
			}

			message := fmt.Sprintf("Stateful resource '%s' of type '%s' should have explicit %s to prevent accidental data loss. Consider setting to 'Retain' or 'Snapshot'.",
				resName, res.Type, joinPolicies(missingPolicies))

			matches = append(matches, rules.Match{
				Message: message,
				Path:    []string{"Resources", resName},
			})
		}
	}

	return matches
}

func joinPolicies(policies []string) string {
	if len(policies) == 1 {
		return policies[0]
	}
	if len(policies) == 2 {
		return policies[0] + " and " + policies[1]
	}
	return ""
}
