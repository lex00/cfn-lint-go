// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3514{})
}

// E3514 validates IAM resource policy ARNs.
type E3514 struct{}

func (r *E3514) ID() string { return "E3514" }

func (r *E3514) ShortDesc() string {
	return "IAM resource policy ARNs"
}

func (r *E3514) Description() string {
	return "Validates that resource ARNs in IAM policies follow compliance standards and proper format."
}

func (r *E3514) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3514"
}

func (r *E3514) Tags() []string {
	return []string{"resources", "properties", "iam", "policy", "arn"}
}

// General AWS ARN pattern
var arnPattern = regexp.MustCompile(`^arn:(aws|aws-cn|aws-us-gov):([a-z0-9-]+):([a-z0-9-]*):(\d{12}|):(.+)$`)

func (r *E3514) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Check IAM policies
	for resName, res := range tmpl.Resources {
		switch res.Type {
		case "AWS::IAM::Policy", "AWS::IAM::Role", "AWS::IAM::User", "AWS::IAM::Group":
			r.checkIAMPolicies(res, resName, &matches)
		case "AWS::S3::BucketPolicy", "AWS::SQS::QueuePolicy", "AWS::SNS::TopicPolicy":
			r.checkResourcePolicy(res, resName, &matches)
		case "AWS::KMS::Key":
			r.checkKMSPolicy(res, resName, &matches)
		}
	}

	return matches
}

func (r *E3514) checkIAMPolicies(res *template.Resource, resName string, matches *[]rules.Match) {
	// Check PolicyDocument for AWS::IAM::Policy
	if res.Type == "AWS::IAM::Policy" {
		if policyDoc, hasDoc := res.Properties["PolicyDocument"]; hasDoc {
			r.validatePolicyARNs(policyDoc, resName, matches, res.Node.Line, res.Node.Column, []string{"Resources", resName, "Properties", "PolicyDocument"})
		}
	}

	// Check Policies array for other IAM resources
	if policies, hasPolicies := res.Properties["Policies"]; hasPolicies {
		if policiesList, ok := policies.([]interface{}); ok {
			for i, policy := range policiesList {
				if policyMap, ok := policy.(map[string]interface{}); ok {
					if policyDoc, hasDoc := policyMap["PolicyDocument"]; hasDoc {
						r.validatePolicyARNs(policyDoc, resName, matches, res.Node.Line, res.Node.Column, []string{"Resources", resName, "Properties", "Policies", fmt.Sprintf("[%d]", i), "PolicyDocument"})
					}
				}
			}
		}
	}
}

func (r *E3514) checkResourcePolicy(res *template.Resource, resName string, matches *[]rules.Match) {
	if policyDoc, hasDoc := res.Properties["PolicyDocument"]; hasDoc {
		r.validatePolicyARNs(policyDoc, resName, matches, res.Node.Line, res.Node.Column, []string{"Resources", resName, "Properties", "PolicyDocument"})
	}
}

func (r *E3514) checkKMSPolicy(res *template.Resource, resName string, matches *[]rules.Match) {
	if keyPolicy, hasPolicy := res.Properties["KeyPolicy"]; hasPolicy {
		r.validatePolicyARNs(keyPolicy, resName, matches, res.Node.Line, res.Node.Column, []string{"Resources", resName, "Properties", "KeyPolicy"})
	}
}

func (r *E3514) validatePolicyARNs(policyDoc interface{}, resName string, matches *[]rules.Match, line, column int, path []string) {
	policyMap, ok := policyDoc.(map[string]interface{})
	if !ok {
		// Try to parse as JSON string
		if policyStr, ok := policyDoc.(string); ok {
			var parsedPolicy map[string]interface{}
			if err := json.Unmarshal([]byte(policyStr), &parsedPolicy); err != nil {
				return
			}
			policyMap = parsedPolicy
		} else {
			return
		}
	}

	// Get statements
	statements, hasStatements := policyMap["Statement"]
	if !hasStatements {
		return
	}

	stmtList, ok := statements.([]interface{})
	if !ok {
		return
	}

	// Check each statement for Resource ARNs
	for i, stmt := range stmtList {
		stmtMap, ok := stmt.(map[string]interface{})
		if !ok {
			continue
		}

		// Check Resource field
		if resource, hasResource := stmtMap["Resource"]; hasResource {
			r.validateResourceARNs(resource, resName, matches, line, column, append(path, "Statement", fmt.Sprintf("[%d]", i), "Resource"))
		}

		// Check NotResource field
		if notResource, hasNotResource := stmtMap["NotResource"]; hasNotResource {
			r.validateResourceARNs(notResource, resName, matches, line, column, append(path, "Statement", fmt.Sprintf("[%d]", i), "NotResource"))
		}
	}
}

func (r *E3514) validateResourceARNs(resource interface{}, resName string, matches *[]rules.Match, line, column int, path []string) {
	switch v := resource.(type) {
	case string:
		r.validateSingleARN(v, resName, matches, line, column, path)
	case []interface{}:
		for _, arn := range v {
			if arnStr, ok := arn.(string); ok {
				r.validateSingleARN(arnStr, resName, matches, line, column, path)
			}
		}
	}
}

func (r *E3514) validateSingleARN(arn string, resName string, matches *[]rules.Match, line, column int, path []string) {
	// Skip wildcards and intrinsic functions
	if arn == "*" || len(arn) == 0 {
		return
	}

	// Only validate if it starts with "arn:"
	if len(arn) < 4 || arn[:4] != "arn:" {
		return
	}

	// Validate ARN format
	if !arnPattern.MatchString(arn) {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf(
				"Resource '%s': Invalid ARN format '%s' in policy Resource field",
				resName, arn,
			),
			Line:   line,
			Column: column,
			Path:   path,
		})
	}
}
