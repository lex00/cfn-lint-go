package warnings

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W3691{})
}

// W3691 warns about deprecated DB Instance engine versions.
type W3691 struct{}

func (r *W3691) ID() string { return "W3691" }

func (r *W3691) ShortDesc() string {
	return "DB Instance deprecated engine version"
}

func (r *W3691) Description() string {
	return "Warns when RDS DBInstance resources use deprecated engine versions that may lose support."
}

func (r *W3691) Source() string {
	return "https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_PostgreSQL.html"
}

func (r *W3691) Tags() []string {
	return []string{"warnings", "rds", "dbinstance", "deprecation"}
}

// Deprecated MySQL versions and their replacements
var deprecatedMySQLVersions = map[string]string{
	"5.5": "MySQL 8.0",
	"5.6": "MySQL 8.0",
	"5.7": "MySQL 8.0",
}

// Deprecated PostgreSQL versions and their replacements
var deprecatedPostgresVersions = map[string]string{
	"9.":  "PostgreSQL 15 or later",
	"10.": "PostgreSQL 15 or later",
	"11.": "PostgreSQL 15 or later",
	"12.": "PostgreSQL 15 or later",
}

// Deprecated MariaDB versions and their replacements
var deprecatedMariaDBVersions = map[string]string{
	"10.0": "MariaDB 10.6 or later",
	"10.1": "MariaDB 10.6 or later",
	"10.2": "MariaDB 10.6 or later",
	"10.3": "MariaDB 10.6 or later",
	"10.4": "MariaDB 10.6 or later",
}

// Deprecated Oracle versions
var deprecatedOracleVersions = map[string]string{
	"11.2": "Oracle 19c or later",
	"12.1": "Oracle 19c or later",
	"12.2": "Oracle 19c or later",
}

// Deprecated SQL Server versions
var deprecatedSQLServerVersions = map[string]string{
	"11.00": "SQL Server 2019 or later",
	"12.00": "SQL Server 2019 or later",
	"13.00": "SQL Server 2019 or later",
	"14.00": "SQL Server 2019 or later",
}

func (r *W3691) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::RDS::DBInstance" {
			continue
		}

		engine, hasEngine := res.Properties["Engine"].(string)
		engineVersion, hasVersion := res.Properties["EngineVersion"].(string)

		if !hasEngine || !hasVersion {
			continue
		}

		engineLower := strings.ToLower(engine)
		versionLower := strings.ToLower(engineVersion)

		var deprecatedVersions map[string]string
		var engineName string

		switch {
		case engineLower == "mysql":
			deprecatedVersions = deprecatedMySQLVersions
			engineName = "MySQL"
		case engineLower == "postgres":
			deprecatedVersions = deprecatedPostgresVersions
			engineName = "PostgreSQL"
		case engineLower == "mariadb":
			deprecatedVersions = deprecatedMariaDBVersions
			engineName = "MariaDB"
		case strings.HasPrefix(engineLower, "oracle"):
			deprecatedVersions = deprecatedOracleVersions
			engineName = "Oracle"
		case strings.HasPrefix(engineLower, "sqlserver"):
			deprecatedVersions = deprecatedSQLServerVersions
			engineName = "SQL Server"
		default:
			continue
		}

		for deprecated, replacement := range deprecatedVersions {
			if strings.HasPrefix(versionLower, deprecated) {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("RDS DBInstance '%s' uses deprecated %s version '%s'; consider upgrading to %s", resName, engineName, engineVersion, replacement),
					Path:    []string{"Resources", resName, "Properties", "EngineVersion"},
				})
				break
			}
		}
	}

	return matches
}
