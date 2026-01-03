package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W2533{})
}

// W2533 warns about Lambda .zip deployment configuration issues.
type W2533 struct{}

func (r *W2533) ID() string { return "W2533" }

func (r *W2533) ShortDesc() string {
	return "Lambda .zip deployment properties"
}

func (r *W2533) Description() string {
	return "Warns about Lambda .zip deployment configuration issues, such as missing or conflicting properties."
}

func (r *W2533) Source() string {
	return "https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-package.html"
}

func (r *W2533) Tags() []string {
	return []string{"warnings", "lambda", "deployment", "zip"}
}

func (r *W2533) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::Lambda::Function" {
			continue
		}

		// Check PackageType - default is Zip
		packageType, _ := res.Properties["PackageType"].(string)
		if packageType == "Image" {
			continue // Not a .zip deployment
		}

		code, hasCode := res.Properties["Code"].(map[string]any)
		if !hasCode {
			continue
		}

		// Check for S3 deployment without S3ObjectVersion
		s3Bucket, hasS3Bucket := code["S3Bucket"]
		_, hasS3Key := code["S3Key"]
		_, hasS3ObjectVersion := code["S3ObjectVersion"]

		if hasS3Bucket && hasS3Key && !hasS3ObjectVersion {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Lambda function '%s' uses S3 deployment without S3ObjectVersion; consider using versioned deployments for rollback capability", resName),
				Path:    []string{"Resources", resName, "Properties", "Code"},
			})
		}

		// Check for ZipFile with large inline code
		if zipFile, ok := code["ZipFile"].(string); ok {
			if len(zipFile) > 4096 {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Lambda function '%s' has large inline ZipFile (%d chars); consider using S3 deployment for better manageability", resName, len(zipFile)),
					Path:    []string{"Resources", resName, "Properties", "Code", "ZipFile"},
				})
			}
		}

		// Check for Handler when using S3 deployment
		_, hasHandler := res.Properties["Handler"]
		if (hasS3Bucket || hasS3Key) && !hasHandler {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Lambda function '%s' uses S3 deployment but Handler is not specified", resName),
				Path:    []string{"Resources", resName, "Properties"},
			})
		}

		// Check for Runtime when using .zip
		_, hasRuntime := res.Properties["Runtime"]
		if !hasRuntime && packageType != "Image" {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Lambda function '%s' is a .zip deployment but Runtime is not specified", resName),
				Path:    []string{"Resources", resName, "Properties"},
			})
		}

		// Check for both S3 and ZipFile (conflicting)
		_, hasZipFile := code["ZipFile"]
		if hasZipFile && (hasS3Bucket || hasS3Key) {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Lambda function '%s' has both ZipFile and S3 properties in Code; use only one", resName),
				Path:    []string{"Resources", resName, "Properties", "Code"},
			})
		}

		// Check for ImageUri in .zip deployment
		if _, hasImageUri := code["ImageUri"]; hasImageUri && packageType != "Image" {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Lambda function '%s' has ImageUri but PackageType is not 'Image'", resName),
				Path:    []string{"Resources", resName, "Properties", "Code", "ImageUri"},
			})
		}

		// Check S3Bucket without S3Key
		if hasS3Bucket && !hasS3Key {
			// This could be valid if using intrinsic functions
			if _, isBucketString := s3Bucket.(string); isBucketString {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Lambda function '%s' has S3Bucket but no S3Key specified", resName),
					Path:    []string{"Resources", resName, "Properties", "Code"},
				})
			}
		}
	}

	return matches
}
