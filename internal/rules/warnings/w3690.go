package warnings

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W3690{})
}

// W3690 warns about deprecated DB Cluster engine versions.
type W3690 struct{}

func (r *W3690) ID() string { return "W3690" }

func (r *W3690) ShortDesc() string {
	return "DB Cluster deprecated engine version"
}

func (r *W3690) Description() string {
	return "Warns when RDS DBCluster resources use deprecated engine versions that may lose support."
}

func (r *W3690) Source() string {
	return "https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/Aurora.VersionPolicy.html"
}

func (r *W3690) Tags() []string {
	return []string{"warnings", "rds", "dbcluster", "deprecation"}
}

// Deprecated Aurora MySQL versions and their replacements
var deprecatedAuroraMySQLVersions = map[string]string{
	"5.6.10a":           "8.0 (Aurora MySQL 3.x)",
	"5.6.mysql_aurora.": "8.0 (Aurora MySQL 3.x)",
	"5.7.12":            "8.0 (Aurora MySQL 3.x)",
	"5.7.mysql_aurora.": "8.0 (Aurora MySQL 3.x) or upgrade to 2.x",
	"2.07":              "Aurora MySQL 3.x (MySQL 8.0)",
	"2.08":              "Aurora MySQL 3.x (MySQL 8.0)",
	"2.09":              "Aurora MySQL 3.x (MySQL 8.0)",
	"2.10":              "Aurora MySQL 3.x (MySQL 8.0)",
}

// Deprecated Aurora PostgreSQL versions and their replacements
var deprecatedAuroraPostgresVersions = map[string]string{
	"9.6": "PostgreSQL 15 or later",
	"10.": "PostgreSQL 15 or later",
	"11.": "PostgreSQL 15 or later",
	"12.": "PostgreSQL 15 or later",
	"13.": "PostgreSQL 15 or later",
}

func (r *W3690) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::RDS::DBCluster" {
			continue
		}

		engine, hasEngine := res.Properties["Engine"].(string)
		engineVersion, hasVersion := res.Properties["EngineVersion"].(string)

		if !hasEngine || !hasVersion {
			continue
		}

		engineLower := strings.ToLower(engine)
		versionLower := strings.ToLower(engineVersion)

		// Check Aurora MySQL
		if strings.Contains(engineLower, "aurora-mysql") || engineLower == "aurora" {
			for deprecated, replacement := range deprecatedAuroraMySQLVersions {
				if strings.HasPrefix(versionLower, strings.ToLower(deprecated)) {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("RDS DBCluster '%s' uses deprecated Aurora MySQL version '%s'; consider upgrading to %s", resName, engineVersion, replacement),
						Path:    []string{"Resources", resName, "Properties", "EngineVersion"},
					})
					break
				}
			}
		}

		// Check Aurora PostgreSQL
		if strings.Contains(engineLower, "aurora-postgresql") {
			for deprecated, replacement := range deprecatedAuroraPostgresVersions {
				if strings.HasPrefix(versionLower, deprecated) {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("RDS DBCluster '%s' uses Aurora PostgreSQL version '%s' which is approaching end-of-life; consider upgrading to %s", resName, engineVersion, replacement),
						Path:    []string{"Resources", resName, "Properties", "EngineVersion"},
					})
					break
				}
			}
		}
	}

	return matches
}
