package warnings

import (
	"fmt"
	"strconv"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W3034{})
}

// W3034 warns when parameter values used in resources may be outside valid ranges.
type W3034 struct{}

func (r *W3034) ID() string { return "W3034" }

func (r *W3034) ShortDesc() string {
	return "Parameter value range check"
}

func (r *W3034) Description() string {
	return "Warns when parameter values used in resources may not satisfy the resource property's expected range."
}

func (r *W3034) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/parameters-section-structure.html"
}

func (r *W3034) Tags() []string {
	return []string{"warnings", "parameters", "resources", "validation"}
}

// Known property ranges for common resources
var knownPropertyRanges = map[string]map[string]struct{ min, max float64 }{
	"AWS::Lambda::Function": {
		"Timeout":    {1, 900},
		"MemorySize": {128, 10240},
	},
	"AWS::SQS::Queue": {
		"VisibilityTimeout":             {0, 43200},
		"MessageRetentionPeriod":        {60, 1209600},
		"MaximumMessageSize":            {1024, 262144},
		"DelaySeconds":                  {0, 900},
		"ReceiveMessageWaitTimeSeconds": {0, 20},
	},
	"AWS::AutoScaling::AutoScalingGroup": {
		"MinSize":                {0, 10000},
		"MaxSize":                {0, 10000},
		"DesiredCapacity":        {0, 10000},
		"HealthCheckGracePeriod": {0, 7200},
		"DefaultCooldown":        {0, 86400},
	},
	"AWS::RDS::DBInstance": {
		"AllocatedStorage":      {20, 65536},
		"BackupRetentionPeriod": {0, 35},
		"MonitoringInterval":    {0, 60},
	},
	"AWS::ElasticLoadBalancingV2::TargetGroup": {
		"HealthCheckIntervalSeconds": {5, 300},
		"HealthCheckTimeoutSeconds":  {2, 120},
		"HealthyThresholdCount":      {2, 10},
		"UnhealthyThresholdCount":    {2, 10},
	},
}

func (r *W3034) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		propRanges, hasRanges := knownPropertyRanges[res.Type]
		if !hasRanges {
			continue
		}

		for propName, ranges := range propRanges {
			propValue, hasProp := res.Properties[propName]
			if !hasProp {
				continue
			}

			// Check if the value is a Ref to a parameter
			if refMap, ok := propValue.(map[string]any); ok {
				if paramName, ok := refMap["Ref"].(string); ok {
					r.checkParameterRange(paramName, propName, ranges.min, ranges.max, resName, res.Type, tmpl, &matches)
				}
			}

			// Check direct numeric values
			var numValue float64
			var isNumber bool
			switch v := propValue.(type) {
			case float64:
				numValue = v
				isNumber = true
			case int:
				numValue = float64(v)
				isNumber = true
			case string:
				if f, err := strconv.ParseFloat(v, 64); err == nil {
					numValue = f
					isNumber = true
				}
			}

			if isNumber {
				if numValue < ranges.min || numValue > ranges.max {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Resource '%s' property '%s' value %v is outside valid range [%v, %v]", resName, propName, numValue, ranges.min, ranges.max),
						Path:    []string{"Resources", resName, "Properties", propName},
					})
				}
			}
		}
	}

	return matches
}

func (r *W3034) checkParameterRange(paramName, propName string, minRange, maxRange float64, resName, resType string, tmpl *template.Template, matches *[]rules.Match) {
	param, exists := tmpl.Parameters[paramName]
	if !exists {
		return
	}

	// Check if parameter constraints allow values outside the property range
	var paramMin, paramMax float64
	hasParamMin, hasParamMax := false, false

	if param.MinValue != nil {
		paramMin = *param.MinValue
		hasParamMin = true
	}

	if param.MaxValue != nil {
		paramMax = *param.MaxValue
		hasParamMax = true
	}

	// Check if parameter range allows values outside property range
	if hasParamMin && paramMin < minRange {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("Parameter '%s' used for '%s.%s' allows minimum value %v which is below property minimum %v", paramName, resType, propName, paramMin, minRange),
			Path:    []string{"Parameters", paramName},
		})
	}

	if hasParamMax && paramMax > maxRange {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("Parameter '%s' used for '%s.%s' allows maximum value %v which exceeds property maximum %v", paramName, resType, propName, paramMax, maxRange),
			Path:    []string{"Parameters", paramName},
		})
	}

	// Check default value
	if param.Default != nil {
		var defaultVal float64
		switch v := param.Default.(type) {
		case float64:
			defaultVal = v
		case int:
			defaultVal = float64(v)
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				defaultVal = f
			} else {
				return
			}
		default:
			return
		}

		if defaultVal < minRange || defaultVal > maxRange {
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf("Parameter '%s' default value %v is outside valid range [%v, %v] for '%s.%s'", paramName, defaultVal, minRange, maxRange, resType, propName),
				Path:    []string{"Parameters", paramName, "Default"},
			})
		}
	}
}
