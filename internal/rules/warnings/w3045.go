package warnings

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W3045{})
}

// W3045 warns about S3 bucket access control configurations.
type W3045 struct{}

func (r *W3045) ID() string { return "W3045" }

func (r *W3045) ShortDesc() string {
	return "S3 bucket policies for access control"
}

func (r *W3045) Description() string {
	return "Warns about S3 bucket access control issues, such as public access or missing encryption requirements."
}

func (r *W3045) Source() string {
	return "https://docs.aws.amazon.com/AmazonS3/latest/userguide/access-control-best-practices.html"
}

func (r *W3045) Tags() []string {
	return []string{"warnings", "s3", "security", "access-control"}
}

func (r *W3045) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::S3::Bucket" {
			continue
		}

		r.checkBucket(resName, res, &matches)
	}

	// Also check S3 bucket policies
	for resName, res := range tmpl.Resources {
		if res.Type == "AWS::S3::BucketPolicy" {
			r.checkBucketPolicy(resName, res, &matches)
		}
	}

	return matches
}

func (r *W3045) checkBucket(resName string, res *template.Resource, matches *[]rules.Match) {
	// Check for public access settings
	if publicAccessConfig, ok := res.Properties["PublicAccessBlockConfiguration"].(map[string]any); ok {
		// Check if any public access is allowed
		blockPublicAcls, _ := publicAccessConfig["BlockPublicAcls"].(bool)
		blockPublicPolicy, _ := publicAccessConfig["BlockPublicPolicy"].(bool)
		ignorePublicAcls, _ := publicAccessConfig["IgnorePublicAcls"].(bool)
		restrictPublicBuckets, _ := publicAccessConfig["RestrictPublicBuckets"].(bool)

		if !blockPublicAcls || !blockPublicPolicy || !ignorePublicAcls || !restrictPublicBuckets {
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf("S3 bucket '%s' does not have all public access block settings enabled; consider blocking all public access", resName),
				Path:    []string{"Resources", resName, "Properties", "PublicAccessBlockConfiguration"},
			})
		}
	} else {
		// No PublicAccessBlockConfiguration at all
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("S3 bucket '%s' does not have PublicAccessBlockConfiguration; consider adding it to prevent public access", resName),
			Path:    []string{"Resources", resName, "Properties"},
		})
	}

	// Check for bucket encryption
	if _, hasEncryption := res.Properties["BucketEncryption"]; !hasEncryption {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("S3 bucket '%s' does not have BucketEncryption configured; consider enabling default encryption", resName),
			Path:    []string{"Resources", resName, "Properties"},
		})
	}

	// Check for versioning
	if versioningConfig, ok := res.Properties["VersioningConfiguration"].(map[string]any); ok {
		status, _ := versioningConfig["Status"].(string)
		if status != "Enabled" {
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf("S3 bucket '%s' versioning is not enabled; consider enabling versioning for data protection", resName),
				Path:    []string{"Resources", resName, "Properties", "VersioningConfiguration"},
			})
		}
	} else {
		*matches = append(*matches, rules.Match{
			Message: fmt.Sprintf("S3 bucket '%s' does not have VersioningConfiguration; consider enabling versioning", resName),
			Path:    []string{"Resources", resName, "Properties"},
		})
	}

	// Check for deprecated ACL usage
	if acl, ok := res.Properties["AccessControl"].(string); ok {
		if strings.Contains(strings.ToLower(acl), "public") {
			*matches = append(*matches, rules.Match{
				Message: fmt.Sprintf("S3 bucket '%s' uses public ACL '%s'; consider using bucket policies instead", resName, acl),
				Path:    []string{"Resources", resName, "Properties", "AccessControl"},
			})
		}
	}
}

func (r *W3045) checkBucketPolicy(resName string, res *template.Resource, matches *[]rules.Match) {
	policyDoc, ok := res.Properties["PolicyDocument"].(map[string]any)
	if !ok {
		return
	}

	statements, ok := policyDoc["Statement"].([]any)
	if !ok {
		return
	}

	for i, stmt := range statements {
		stmtMap, ok := stmt.(map[string]any)
		if !ok {
			continue
		}

		effect, _ := stmtMap["Effect"].(string)
		if !strings.EqualFold(effect, "Allow") {
			continue
		}

		// Check for public principal
		if principal, ok := stmtMap["Principal"]; ok {
			if principalStr, ok := principal.(string); ok && principalStr == "*" {
				// Check if there are conditions that limit access
				_, hasCondition := stmtMap["Condition"]
				if !hasCondition {
					*matches = append(*matches, rules.Match{
						Message: fmt.Sprintf("S3 bucket policy '%s' statement %d allows public access without conditions", resName, i),
						Path:    []string{"Resources", resName, "Properties", "PolicyDocument", "Statement", fmt.Sprintf("[%d]", i)},
					})
				}
			}
		}

		// Check for HTTP access (should require HTTPS)
		if actions, ok := stmtMap["Action"]; ok {
			hasGetObject := false
			switch a := actions.(type) {
			case string:
				hasGetObject = strings.Contains(a, "GetObject") || a == "s3:*" || a == "*"
			case []any:
				for _, act := range a {
					if actStr, ok := act.(string); ok {
						if strings.Contains(actStr, "GetObject") || actStr == "s3:*" || actStr == "*" {
							hasGetObject = true
							break
						}
					}
				}
			}

			if hasGetObject {
				// Check for secure transport condition
				hasSecureTransport := false
				if conditions, ok := stmtMap["Condition"].(map[string]any); ok {
					if boolCond, ok := conditions["Bool"].(map[string]any); ok {
						if _, ok := boolCond["aws:SecureTransport"]; ok {
							hasSecureTransport = true
						}
					}
				}

				if !hasSecureTransport {
					*matches = append(*matches, rules.Match{
						Message: fmt.Sprintf("S3 bucket policy '%s' allows GetObject but does not require secure transport (HTTPS)", resName),
						Path:    []string{"Resources", resName, "Properties", "PolicyDocument", "Statement", fmt.Sprintf("[%d]", i)},
					})
				}
			}
		}
	}
}
