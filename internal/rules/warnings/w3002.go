// Package warnings contains warning-level rules (Wxxx).
package warnings

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W3002{})
}

// W3002 warns about properties that require packaging (aws cloudformation package).
type W3002 struct{}

func (r *W3002) ID() string { return "W3002" }

func (r *W3002) ShortDesc() string {
	return "Package-required property with local path"
}

func (r *W3002) Description() string {
	return "Warns when a property that requires packaging (Code, Content, etc.) contains a local file path instead of an S3 URI."
}

func (r *W3002) Source() string {
	return "https://docs.aws.amazon.com/cli/latest/reference/cloudformation/package.html"
}

func (r *W3002) Tags() []string {
	return []string{"warnings", "resources", "packaging"}
}

// packageRequiredProps maps resource types to properties that need packaging
var packageRequiredProps = map[string][]string{
	"AWS::Lambda::Function": {"Code"},
	"AWS::Lambda::LayerVersion": {"Content"},
	"AWS::Serverless::Function": {"CodeUri"},
	"AWS::Serverless::LayerVersion": {"ContentUri"},
	"AWS::ApiGateway::RestApi": {"BodyS3Location"},
	"AWS::AppSync::GraphQLSchema": {"DefinitionS3Location"},
	"AWS::AppSync::Resolver": {"RequestMappingTemplateS3Location", "ResponseMappingTemplateS3Location"},
	"AWS::CloudFormation::Stack": {"TemplateURL"},
	"AWS::Glue::Job": {"Command"},
	"AWS::StepFunctions::StateMachine": {"DefinitionS3Location"},
}

func (r *W3002) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		props, ok := packageRequiredProps[res.Type]
		if !ok {
			continue
		}

		for _, propName := range props {
			if propVal, exists := res.Properties[propName]; exists {
				if isLocalPath(propVal) {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Property '%s' in resource '%s' appears to be a local path. Use 'aws cloudformation package' or provide an S3 URI.", propName, resName),
						Path:    []string{"Resources", resName, "Properties", propName},
					})
				}
			}
		}
	}

	return matches
}

func isLocalPath(v any) bool {
	switch val := v.(type) {
	case string:
		// Check for local path indicators
		if strings.HasPrefix(val, "./") ||
			strings.HasPrefix(val, "../") ||
			strings.HasPrefix(val, "/") ||
			(len(val) > 0 && !strings.HasPrefix(val, "s3://") && !strings.HasPrefix(val, "https://") && !strings.Contains(val, "{{")) {
			// Exclude S3 URIs and HTTPS URLs and dynamic refs
			// Check if it looks like a file path
			return strings.Contains(val, "/") || strings.HasSuffix(val, ".zip") || strings.HasSuffix(val, ".jar") || strings.HasSuffix(val, ".yaml") || strings.HasSuffix(val, ".json")
		}
	case map[string]any:
		// Check nested properties like Code: { S3Bucket: ..., S3Key: ... } vs Code: { ZipFile: ... }
		// If it has ZipFile, it's inline code (valid)
		// If it has local path strings, warn
		if _, hasZipFile := val["ZipFile"]; hasZipFile {
			return false
		}
		if _, hasS3Bucket := val["S3Bucket"]; hasS3Bucket {
			return false
		}
		// Check for string values that look like paths
		for _, child := range val {
			if isLocalPath(child) {
				return true
			}
		}
	}
	return false
}
