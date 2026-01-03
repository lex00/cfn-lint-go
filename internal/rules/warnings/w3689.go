package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W3689{})
}

// W3689 warns about DBInstance properties that are ignored when restoring from source.
type W3689 struct{}

func (r *W3689) ID() string { return "W3689" }

func (r *W3689) ShortDesc() string {
	return "Source DB ignored properties"
}

func (r *W3689) Description() string {
	return "Warns when RDS DBInstance resources specify properties that are ignored when restoring from a snapshot or read replica."
}

func (r *W3689) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-rds-dbinstance.html"
}

func (r *W3689) Tags() []string {
	return []string{"warnings", "rds", "dbinstance", "restore"}
}

// Properties ignored during restore from snapshot
var dbInstanceSnapshotIgnoredProps = []string{
	"DBName",
	"MasterUsername",
	"MasterUserPassword",
	"StorageEncrypted",
	"KmsKeyId",
	"Engine",
	"CharacterSetName",
}

// Properties ignored for read replicas
var dbInstanceReplicaIgnoredProps = []string{
	"DBName",
	"MasterUsername",
	"MasterUserPassword",
	"Engine",
	"EngineVersion",
	"AllocatedStorage",
	"StorageEncrypted",
	"KmsKeyId",
	"CharacterSetName",
	"Port",
}

func (r *W3689) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::RDS::DBInstance" {
			continue
		}

		// Check if this is a snapshot restore
		_, hasSnapshotId := res.Properties["DBSnapshotIdentifier"]
		if hasSnapshotId {
			for _, propName := range dbInstanceSnapshotIgnoredProps {
				if _, hasProp := res.Properties[propName]; hasProp {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("RDS DBInstance '%s' specifies '%s' but this property is ignored when restoring from snapshot", resName, propName),
						Path:    []string{"Resources", resName, "Properties", propName},
					})
				}
			}
			continue
		}

		// Check if this is a read replica
		_, hasSourceDBInstanceIdentifier := res.Properties["SourceDBInstanceIdentifier"]
		if hasSourceDBInstanceIdentifier {
			for _, propName := range dbInstanceReplicaIgnoredProps {
				if _, hasProp := res.Properties[propName]; hasProp {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("RDS DBInstance '%s' specifies '%s' but this property is ignored for read replicas", resName, propName),
						Path:    []string{"Resources", resName, "Properties", propName},
					})
				}
			}
		}
	}

	return matches
}
