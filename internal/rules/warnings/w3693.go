package warnings

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W3693{})
}

// W3693 warns about Aurora DB cluster properties that are ignored.
type W3693 struct{}

func (r *W3693) ID() string { return "W3693" }

func (r *W3693) ShortDesc() string {
	return "Aurora DB cluster ignored properties"
}

func (r *W3693) Description() string {
	return "Warns when Aurora DBCluster resources specify properties that are ignored or not applicable to Aurora."
}

func (r *W3693) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-rds-dbcluster.html"
}

func (r *W3693) Tags() []string {
	return []string{"warnings", "rds", "aurora", "dbcluster"}
}

// Properties that are ignored for Aurora Serverless
var auroraServerlessIgnoredProps = []string{
	"DBClusterInstanceClass",
	"AllocatedStorage",
	"Iops",
	"StorageType",
}

// Properties that are ignored for Aurora Provisioned
var auroraProvisionedIgnoredProps = []string{
	"ServerlessV2ScalingConfiguration",
}

// Properties specific to Aurora that don't apply to non-Aurora
var aurorOnlyProps = []string{
	"ServerlessV2ScalingConfiguration",
	"ScalingConfiguration",
	"EnableHttpEndpoint",
	"GlobalClusterIdentifier",
}

func (r *W3693) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::RDS::DBCluster" {
			continue
		}

		engine, hasEngine := res.Properties["Engine"].(string)
		if !hasEngine {
			continue
		}

		engineLower := strings.ToLower(engine)
		isAurora := strings.Contains(engineLower, "aurora")

		// Check for Aurora-only properties on non-Aurora clusters
		if !isAurora {
			for _, propName := range aurorOnlyProps {
				if _, hasProp := res.Properties[propName]; hasProp {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("RDS DBCluster '%s' uses '%s' but engine '%s' is not Aurora; this property is Aurora-specific", resName, propName, engine),
						Path:    []string{"Resources", resName, "Properties", propName},
					})
				}
			}
			continue
		}

		// Check for Aurora Serverless vs Provisioned
		engineMode, _ := res.Properties["EngineMode"].(string)
		_, hasServerlessV2Config := res.Properties["ServerlessV2ScalingConfiguration"]
		_, hasScalingConfig := res.Properties["ScalingConfiguration"]

		isServerless := strings.ToLower(engineMode) == "serverless" || hasServerlessV2Config || hasScalingConfig

		if isServerless {
			// Check for properties ignored in Serverless mode
			for _, propName := range auroraServerlessIgnoredProps {
				if _, hasProp := res.Properties[propName]; hasProp {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Aurora DBCluster '%s' specifies '%s' but this property is ignored for Aurora Serverless", resName, propName),
						Path:    []string{"Resources", resName, "Properties", propName},
					})
				}
			}
		} else {
			// Check for properties only for Serverless
			for _, propName := range auroraProvisionedIgnoredProps {
				if _, hasProp := res.Properties[propName]; hasProp && engineMode != "" && strings.ToLower(engineMode) != "serverless" {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Aurora DBCluster '%s' specifies '%s' but EngineMode is '%s'; this property is for Serverless mode", resName, propName, engineMode),
						Path:    []string{"Resources", resName, "Properties", propName},
					})
				}
			}
		}

		// Check for conflicting engine mode and scaling configuration
		if strings.ToLower(engineMode) == "serverless" && hasServerlessV2Config {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Aurora DBCluster '%s' has EngineMode 'serverless' (v1) with ServerlessV2ScalingConfiguration; use ScalingConfiguration for Serverless v1 or remove EngineMode for v2", resName),
				Path:    []string{"Resources", resName, "Properties"},
			})
		}

		// Check for both ScalingConfiguration and ServerlessV2ScalingConfiguration
		if hasScalingConfig && hasServerlessV2Config {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Aurora DBCluster '%s' has both ScalingConfiguration (v1) and ServerlessV2ScalingConfiguration (v2); use only one", resName),
				Path:    []string{"Resources", resName, "Properties"},
			})
		}
	}

	return matches
}
