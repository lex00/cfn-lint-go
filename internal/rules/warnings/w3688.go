package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W3688{})
}

// W3688 warns about DBCluster properties that are ignored during restore operations.
type W3688 struct{}

func (r *W3688) ID() string { return "W3688" }

func (r *W3688) ShortDesc() string {
	return "DBCluster restore ignored properties"
}

func (r *W3688) Description() string {
	return "Warns when RDS DBCluster resources specify properties that are ignored during restore from snapshot or point-in-time recovery."
}

func (r *W3688) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-rds-dbcluster.html"
}

func (r *W3688) Tags() []string {
	return []string{"warnings", "rds", "dbcluster", "restore"}
}

// Properties ignored during restore
var dbClusterRestoreIgnoredProps = []string{
	"DatabaseName",
	"MasterUsername",
	"MasterUserPassword",
	"StorageEncrypted",
	"KmsKeyId",
	"Engine",
	"EngineVersion",
}

func (r *W3688) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::RDS::DBCluster" {
			continue
		}

		// Check if this is a restore operation
		_, hasSnapshotId := res.Properties["SnapshotIdentifier"]
		_, hasSourceDBClusterIdentifier := res.Properties["SourceDBClusterIdentifier"]
		_, hasRestoreToTime := res.Properties["RestoreToTime"]
		_, hasRestoreType := res.Properties["RestoreType"]
		_, hasUseLatestRestorableTime := res.Properties["UseLatestRestorableTime"]

		isRestore := hasSnapshotId || hasSourceDBClusterIdentifier || hasRestoreToTime || hasRestoreType || hasUseLatestRestorableTime

		if !isRestore {
			continue
		}

		// Check for ignored properties
		for _, propName := range dbClusterRestoreIgnoredProps {
			if _, hasProp := res.Properties[propName]; hasProp {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("RDS DBCluster '%s' specifies '%s' but this property is ignored during restore operations", resName, propName),
					Path:    []string{"Resources", resName, "Properties", propName},
				})
			}
		}
	}

	return matches
}
