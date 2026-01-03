package informational

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&I3013{})
}

// I3013 suggests explicit retention periods for resources with auto-expiring content.
type I3013 struct{}

func (r *I3013) ID() string { return "I3013" }

func (r *I3013) ShortDesc() string {
	return "Auto-expiring content needs retention period"
}

func (r *I3013) Description() string {
	return "Suggests setting explicit retention periods for resources that auto-expire content (CloudWatch Logs, etc.) to prevent unexpected data loss."
}

func (r *I3013) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-logs-loggroup.html"
}

func (r *I3013) Tags() []string {
	return []string{"resources", "best-practice", "retention"}
}

// retentionPropertyMap maps resource types to their retention property names
var retentionPropertyMap = map[string]string{
	"AWS::Logs::LogGroup":            "RetentionInDays",
	"AWS::CloudWatch::LogGroup":      "RetentionInDays",
	"AWS::EC2::FlowLog":              "RetentionInDays",
	"AWS::Backup::BackupVault":       "BackupVaultTags", // Lifecycle policies
	"AWS::S3::Bucket":                "LifecycleConfiguration",
	"AWS::ECR::Repository":           "LifecyclePolicy",
	"AWS::Kinesis::Stream":           "RetentionPeriodHours",
}

func (r *I3013) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		retentionProperty, hasRetentionProperty := retentionPropertyMap[res.Type]
		if !hasRetentionProperty {
			continue
		}

		// Check if the retention property is set
		if _, hasProperty := res.Properties[retentionProperty]; !hasProperty {
			message := fmt.Sprintf("Resource '%s' of type '%s' does not have explicit '%s'. Consider setting a retention period to control data lifecycle.",
				resName, res.Type, retentionProperty)

			matches = append(matches, rules.Match{
				Message: message,
				Path:    []string{"Resources", resName, "Properties"},
			})
		}
	}

	return matches
}
