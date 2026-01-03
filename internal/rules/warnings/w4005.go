package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W4005{})
}

// W4005 warns about cfn-lint configuration in Metadata.
type W4005 struct{}

func (r *W4005) ID() string { return "W4005" }

func (r *W4005) ShortDesc() string {
	return "cfnlint configuration in Metadata"
}

func (r *W4005) Description() string {
	return "Warns about potential issues with cfn-lint configuration in template Metadata, such as unknown options or deprecated settings."
}

func (r *W4005) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint#metadata"
}

func (r *W4005) Tags() []string {
	return []string{"warnings", "metadata", "cfn-lint", "configuration"}
}

// Valid cfn-lint config keys
var validCfnLintConfigKeys = map[string]bool{
	"ignore_checks":        true,
	"include_checks":       true,
	"configure_rules":      true,
	"ignore_templates":     true,
	"include_experimental": true,
	"regions":              true,
	"ignore_bad_template":  true,
}

func (r *W4005) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check template-level Metadata
	if tmpl.Metadata != nil {
		r.checkMetadata(tmpl.Metadata, []string{"Metadata"}, &matches)
	}

	// Check resource-level Metadata
	for resName, res := range tmpl.Resources {
		if res.Metadata != nil {
			r.checkMetadata(res.Metadata, []string{"Resources", resName, "Metadata"}, &matches)
		}
	}

	return matches
}

func (r *W4005) checkMetadata(metadata map[string]any, path []string, matches *[]rules.Match) {
	// Check for cfn-lint configuration
	cfnLint, hasCfnLint := metadata["cfn-lint"]
	if !hasCfnLint {
		return
	}

	cfnLintMap, ok := cfnLint.(map[string]any)
	if !ok {
		*matches = append(*matches, rules.Match{
			Message: "cfn-lint metadata should be an object",
			Path:    append(path, "cfn-lint"),
		})
		return
	}

	// Check for config section
	config, hasConfig := cfnLintMap["config"]
	if !hasConfig {
		// Check if there are direct config keys (older format)
		for key := range cfnLintMap {
			if validCfnLintConfigKeys[key] {
				*matches = append(*matches, rules.Match{
					Message: fmt.Sprintf("cfn-lint configuration key '%s' should be inside a 'config' object", key),
					Path:    append(path, "cfn-lint", key),
				})
			}
		}
		return
	}

	configMap, ok := config.(map[string]any)
	if !ok {
		*matches = append(*matches, rules.Match{
			Message: "cfn-lint config should be an object",
			Path:    append(path, "cfn-lint", "config"),
		})
		return
	}

	// Check for unknown config keys
	for key := range configMap {
		if !validCfnLintConfigKeys[key] {
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf("Unknown cfn-lint config key '%s'", key),
				Path:    append(path, "cfn-lint", "config", key),
			})
		}
	}

	// Check ignore_checks format
	if ignoreChecks, ok := configMap["ignore_checks"]; ok {
		r.checkRuleList(ignoreChecks, "ignore_checks", append(path, "cfn-lint", "config", "ignore_checks"), matches)
	}

	// Check include_checks format
	if includeChecks, ok := configMap["include_checks"]; ok {
		r.checkRuleList(includeChecks, "include_checks", append(path, "cfn-lint", "config", "include_checks"), matches)
	}

	// Check regions format
	if regions, ok := configMap["regions"]; ok {
		r.checkRegions(regions, append(path, "cfn-lint", "config", "regions"), matches)
	}
}

func (r *W4005) checkRuleList(v any, configKey string, path []string, matches *[]rules.Match) {
	list, ok := v.([]any)
	if !ok {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("cfn-lint %s should be a list of rule IDs", configKey),
			Path:    path,
		})
		return
	}

	for i, item := range list {
		ruleID, ok := item.(string)
		if !ok {
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf("cfn-lint %s item %d should be a string rule ID", configKey, i),
				Path:    append(path, fmt.Sprintf("[%d]", i)),
			})
			continue
		}

		// Check rule ID format (should start with E, W, or I followed by digits)
		if len(ruleID) < 2 {
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf("Invalid rule ID '%s' in %s; rule IDs should be like 'E1001', 'W2001', etc.", ruleID, configKey),
				Path:    append(path, fmt.Sprintf("[%d]", i)),
			})
		} else if ruleID[0] != 'E' && ruleID[0] != 'W' && ruleID[0] != 'I' {
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf("Invalid rule ID '%s' in %s; rule IDs should start with E, W, or I", ruleID, configKey),
				Path:    append(path, fmt.Sprintf("[%d]", i)),
			})
		}
	}
}

func (r *W4005) checkRegions(v any, path []string, matches *[]rules.Match) {
	list, ok := v.([]any)
	if !ok {
		*matches = append(*matches, rules.Match{
			Message: "cfn-lint regions should be a list of AWS region codes",
			Path:    path,
		})
		return
	}

	validRegions := map[string]bool{
		"us-east-1": true, "us-east-2": true, "us-west-1": true, "us-west-2": true,
		"af-south-1": true, "ap-east-1": true, "ap-south-1": true, "ap-south-2": true,
		"ap-southeast-1": true, "ap-southeast-2": true, "ap-southeast-3": true, "ap-southeast-4": true,
		"ap-northeast-1": true, "ap-northeast-2": true, "ap-northeast-3": true,
		"ca-central-1": true, "ca-west-1": true,
		"eu-central-1": true, "eu-central-2": true,
		"eu-west-1": true, "eu-west-2": true, "eu-west-3": true,
		"eu-south-1": true, "eu-south-2": true, "eu-north-1": true,
		"il-central-1": true, "me-south-1": true, "me-central-1": true,
		"sa-east-1": true,
	}

	for i, item := range list {
		region, ok := item.(string)
		if !ok {
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf("cfn-lint regions item %d should be a string", i),
				Path:    append(path, fmt.Sprintf("[%d]", i)),
			})
			continue
		}

		if !validRegions[region] {
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf("Unknown AWS region '%s' in cfn-lint regions", region),
				Path:    append(path, fmt.Sprintf("[%d]", i)),
			})
		}
	}
}
