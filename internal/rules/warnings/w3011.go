package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
	"gopkg.in/yaml.v3"
)

func init() {
	rules.Register(&W3011{})
}

// W3011 warns when UpdateReplacePolicy and DeletionPolicy have potentially conflicting configurations.
type W3011 struct{}

func (r *W3011) ID() string { return "W3011" }

func (r *W3011) ShortDesc() string {
	return "UpdateReplacePolicy/DeletionPolicy both set"
}

func (r *W3011) Description() string {
	return "Warns when UpdateReplacePolicy and DeletionPolicy are both set with potentially conflicting values."
}

func (r *W3011) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-attribute-deletionpolicy.html"
}

func (r *W3011) Tags() []string {
	return []string{"warnings", "resources", "deletion-policy", "update-replace-policy"}
}

func (r *W3011) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		deletionPolicy, updateReplacePolicy := r.getPolicies(res)

		// Skip if neither is set
		if deletionPolicy == "" && updateReplacePolicy == "" {
			continue
		}

		// Check for potentially conflicting combinations
		if deletionPolicy == "Retain" && updateReplacePolicy == "Delete" {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Resource '%s' has DeletionPolicy 'Retain' but UpdateReplacePolicy 'Delete'; replacement updates will delete the resource", resName),
				Path:    []string{"Resources", resName},
			})
		}

		if deletionPolicy == "Delete" && updateReplacePolicy == "Retain" {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Resource '%s' has DeletionPolicy 'Delete' but UpdateReplacePolicy 'Retain'; consider if this asymmetry is intentional", resName),
				Path:    []string{"Resources", resName},
			})
		}

		// Check for resources that should have protection
		criticalResourceTypes := map[string]bool{
			"AWS::RDS::DBInstance":           true,
			"AWS::RDS::DBCluster":            true,
			"AWS::DynamoDB::Table":           true,
			"AWS::S3::Bucket":                true,
			"AWS::EFS::FileSystem":           true,
			"AWS::Elasticsearch::Domain":     true,
			"AWS::OpenSearchService::Domain": true,
		}

		if criticalResourceTypes[res.Type] {
			if deletionPolicy == "" && updateReplacePolicy == "" {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Resource '%s' of type '%s' has no DeletionPolicy or UpdateReplacePolicy; consider adding protection for data resources", resName, res.Type),
					Path:    []string{"Resources", resName},
				})
			}
		}

		// Check for Snapshot policy on resources that don't support it
		snapshotUnsupportedTypes := map[string]bool{
			"AWS::S3::Bucket":       true,
			"AWS::Lambda::Function": true,
			"AWS::IAM::Role":        true,
		}

		if deletionPolicy == "Snapshot" && snapshotUnsupportedTypes[res.Type] {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Resource '%s' of type '%s' has DeletionPolicy 'Snapshot' but this resource type does not support snapshots", resName, res.Type),
				Path:    []string{"Resources", resName, "DeletionPolicy"},
			})
		}

		if updateReplacePolicy == "Snapshot" && snapshotUnsupportedTypes[res.Type] {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Resource '%s' of type '%s' has UpdateReplacePolicy 'Snapshot' but this resource type does not support snapshots", resName, res.Type),
				Path:    []string{"Resources", resName, "UpdateReplacePolicy"},
			})
		}
	}

	return matches
}

func (r *W3011) getPolicies(res *template.Resource) (string, string) {
	var deletionPolicy, updateReplacePolicy string

	if res.Node == nil || res.Node.Kind != yaml.MappingNode {
		return "", ""
	}

	for i := 0; i < len(res.Node.Content); i += 2 {
		key := res.Node.Content[i]
		value := res.Node.Content[i+1]

		switch key.Value {
		case "DeletionPolicy":
			deletionPolicy = value.Value
		case "UpdateReplacePolicy":
			updateReplacePolicy = value.Value
		}
	}

	return deletionPolicy, updateReplacePolicy
}
