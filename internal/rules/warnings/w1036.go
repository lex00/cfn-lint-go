package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W1036{})
}

// W1036 warns about GetAZs function value issues.
type W1036 struct{}

func (r *W1036) ID() string { return "W1036" }

func (r *W1036) ShortDesc() string {
	return "GetAZs function value validation"
}

func (r *W1036) Description() string {
	return "Warns about potential issues with Fn::GetAZs values, such as hardcoded regions."
}

func (r *W1036) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getavailabilityzones.html"
}

func (r *W1036) Tags() []string {
	return []string{"warnings", "functions", "getazs", "availability-zones"}
}

func (r *W1036) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check resources
	for resName, res := range tmpl.Resources {
		r.checkValue(res.Properties, []string{"Resources", resName, "Properties"}, &matches)
	}

	// Check outputs
	for outName, out := range tmpl.Outputs {
		r.checkValue(out.Value, []string{"Outputs", outName, "Value"}, &matches)
	}

	return matches
}

func (r *W1036) checkValue(v any, path []string, matches *[]rules.Match) {
	switch val := v.(type) {
	case map[string]any:
		if getAZs, ok := val["Fn::GetAZs"]; ok {
			r.checkGetAZs(getAZs, path, matches)
		}
		for key, child := range val {
			r.checkValue(child, append(path, key), matches)
		}
	case []any:
		for i, child := range val {
			r.checkValue(child, append(path, fmt.Sprintf("[%d]", i)), matches)
		}
	}
}

func (r *W1036) checkGetAZs(getAZs any, path []string, matches *[]rules.Match) {
	// Fn::GetAZs takes a region string or "" for current region
	region, ok := getAZs.(string)
	if !ok {
		// Could be a Ref or other intrinsic, which is fine
		return
	}

	// Check if a specific region is hardcoded
	if region != "" {
		regions := map[string]bool{
			"us-east-1": true, "us-east-2": true, "us-west-1": true, "us-west-2": true,
			"eu-west-1": true, "eu-west-2": true, "eu-west-3": true, "eu-central-1": true,
			"ap-northeast-1": true, "ap-northeast-2": true, "ap-southeast-1": true, "ap-southeast-2": true,
			"ap-south-1": true, "sa-east-1": true, "ca-central-1": true,
			"af-south-1": true, "ap-east-1": true, "ap-northeast-3": true, "eu-north-1": true,
			"eu-south-1": true, "me-south-1": true,
		}
		if regions[region] {
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf("Fn::GetAZs uses hardcoded region '%s'; consider using '' or { Ref: AWS::Region } for portability", region),
				Path:    path,
			})
		}
	}
}
