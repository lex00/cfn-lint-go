package sam

import (
	"encoding/json"
	"fmt"

	"github.com/lex00/aws-sam-translator-go/pkg/translator"
	"github.com/lex00/cfn-lint-go/pkg/template"
	"gopkg.in/yaml.v3"
)

// TransformOptions configures the SAM transformation.
type TransformOptions struct {
	// Region is the AWS region for transformation context.
	Region string

	// AccountID is the AWS account ID for transformation context.
	AccountID string

	// StackName is the CloudFormation stack name for transformation context.
	StackName string

	// Partition is the AWS partition (aws, aws-cn, aws-us-gov).
	Partition string
}

// TransformResult contains the transformed template and metadata.
type TransformResult struct {
	// Template is the transformed CloudFormation template.
	Template *template.Template

	// SourceMap maps transformed resources back to original SAM template locations.
	SourceMap *SourceMap

	// Warnings contains any warnings generated during transformation.
	Warnings []string
}

// DefaultTransformOptions returns sensible default options for transformation.
func DefaultTransformOptions() *TransformOptions {
	return &TransformOptions{
		Region:    "us-east-1",
		AccountID: "123456789012",
		StackName: "sam-app",
		Partition: "aws",
	}
}

// Transform converts a SAM template to CloudFormation.
// If the template is not a SAM template, it returns it unchanged.
func Transform(tmpl *template.Template, opts *TransformOptions) (*TransformResult, error) {
	if tmpl == nil {
		return nil, fmt.Errorf("template is nil")
	}

	// If not a SAM template, return it unchanged with empty source map
	if !IsSAMTemplate(tmpl) {
		return &TransformResult{
			Template:  tmpl,
			SourceMap: NewSourceMap(),
			Warnings:  nil,
		}, nil
	}

	// Use default options if none provided
	if opts == nil {
		opts = DefaultTransformOptions()
	}

	// Build source map from original SAM resources before transformation
	sourceMap := buildSourceMap(tmpl)

	// Convert template to YAML bytes for the translator
	yamlBytes, err := templateToYAML(tmpl)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize template: %w", err)
	}

	// Create translator with options
	tr := translator.NewWithOptions(translator.Options{
		Region:    opts.Region,
		AccountID: opts.AccountID,
		StackName: opts.StackName,
		Partition: opts.Partition,
	})

	// Transform SAM to CloudFormation
	cfnBytes, err := tr.TransformBytes(yamlBytes)
	if err != nil {
		return nil, fmt.Errorf("SAM transform failed: %w", err)
	}

	// Parse the transformed CloudFormation template
	cfnTmpl, err := template.Parse(cfnBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse transformed template: %w", err)
	}

	// Update source map with new resource mappings
	updateSourceMapWithTransformedResources(sourceMap, tmpl, cfnTmpl)

	return &TransformResult{
		Template:  cfnTmpl,
		SourceMap: sourceMap,
		Warnings:  nil,
	}, nil
}

// buildSourceMap creates a source map from the original SAM template.
func buildSourceMap(tmpl *template.Template) *SourceMap {
	sm := NewSourceMap()

	for name, res := range tmpl.Resources {
		if res.Node != nil {
			sm.AddResourceMapping(name, name, res.Node.Line, res.Node.Column)
		}
	}

	return sm
}

// updateSourceMapWithTransformedResources adds mappings for generated CFN resources.
func updateSourceMapWithTransformedResources(sm *SourceMap, samTmpl, cfnTmpl *template.Template) {
	// For each CFN resource that wasn't in the original SAM template,
	// try to find the SAM resource it was generated from
	for cfnName := range cfnTmpl.Resources {
		// Skip if we already have a mapping (it's an original resource)
		if _, exists := sm.ResourceMapping[cfnName]; exists {
			continue
		}

		// Try to find the source SAM resource
		// SAM typically generates resources with names like:
		// - MyFunctionRole (from MyFunction)
		// - MyApiDeployment (from MyApi)
		// - MyTableScalingTarget (from MyTable)
		for samName, samRes := range samTmpl.Resources {
			if IsSAMResourceType(samRes.Type) && isGeneratedFrom(cfnName, samName) {
				if samRes.Node != nil {
					sm.AddResourceMapping(cfnName, samName, samRes.Node.Line, samRes.Node.Column)
				}
				break
			}
		}
	}
}

// isGeneratedFrom checks if a CFN resource name was likely generated from a SAM resource.
func isGeneratedFrom(cfnName, samName string) bool {
	// Common patterns:
	// - Exact match (when SAM resource is passed through)
	// - CFN name starts with SAM name (e.g., MyFunctionRole from MyFunction)
	if cfnName == samName {
		return true
	}
	if len(cfnName) > len(samName) && cfnName[:len(samName)] == samName {
		return true
	}
	return false
}

// templateToYAML converts a template back to YAML bytes.
func templateToYAML(tmpl *template.Template) ([]byte, error) {
	// Build a map representation of the template
	data := make(map[string]any)

	if tmpl.AWSTemplateFormatVersion != "" {
		data["AWSTemplateFormatVersion"] = tmpl.AWSTemplateFormatVersion
	}
	if tmpl.Description != "" {
		data["Description"] = tmpl.Description
	}
	if tmpl.Transform != nil {
		data["Transform"] = tmpl.Transform
	}

	if len(tmpl.Parameters) > 0 {
		params := make(map[string]any)
		for name, p := range tmpl.Parameters {
			param := make(map[string]any)
			if p.Type != "" {
				param["Type"] = p.Type
			}
			if p.Description != "" {
				param["Description"] = p.Description
			}
			if p.Default != nil {
				param["Default"] = p.Default
			}
			if len(p.AllowedValues) > 0 {
				param["AllowedValues"] = p.AllowedValues
			}
			if p.AllowedPattern != "" {
				param["AllowedPattern"] = p.AllowedPattern
			}
			if p.MinValue != nil {
				param["MinValue"] = *p.MinValue
			}
			if p.MaxValue != nil {
				param["MaxValue"] = *p.MaxValue
			}
			if p.MinLength != nil {
				param["MinLength"] = *p.MinLength
			}
			if p.MaxLength != nil {
				param["MaxLength"] = *p.MaxLength
			}
			if p.NoEcho {
				param["NoEcho"] = true
			}
			if p.ConstraintDescription != "" {
				param["ConstraintDescription"] = p.ConstraintDescription
			}
			params[name] = param
		}
		data["Parameters"] = params
	}

	if len(tmpl.Mappings) > 0 {
		mappings := make(map[string]any)
		for name, m := range tmpl.Mappings {
			mappings[name] = m.Values
		}
		data["Mappings"] = mappings
	}

	if len(tmpl.Conditions) > 0 {
		conditions := make(map[string]any)
		for name, c := range tmpl.Conditions {
			conditions[name] = c.Expression
		}
		data["Conditions"] = conditions
	}

	if len(tmpl.Resources) > 0 {
		resources := make(map[string]any)
		for name, r := range tmpl.Resources {
			res := make(map[string]any)
			res["Type"] = r.Type
			if len(r.Properties) > 0 {
				res["Properties"] = r.Properties
			}
			if len(r.DependsOn) > 0 {
				if len(r.DependsOn) == 1 {
					res["DependsOn"] = r.DependsOn[0]
				} else {
					res["DependsOn"] = r.DependsOn
				}
			}
			if r.Condition != "" {
				res["Condition"] = r.Condition
			}
			if len(r.Metadata) > 0 {
				res["Metadata"] = r.Metadata
			}
			resources[name] = res
		}
		data["Resources"] = resources
	}

	if len(tmpl.Outputs) > 0 {
		outputs := make(map[string]any)
		for name, o := range tmpl.Outputs {
			out := make(map[string]any)
			out["Value"] = o.Value
			if o.Description != "" {
				out["Description"] = o.Description
			}
			if len(o.Export) > 0 {
				out["Export"] = o.Export
			}
			if o.Condition != "" {
				out["Condition"] = o.Condition
			}
			outputs[name] = out
		}
		data["Outputs"] = outputs
	}

	if len(tmpl.Metadata) > 0 {
		data["Metadata"] = tmpl.Metadata
	}

	// Use JSON for the translator (it handles both YAML and JSON)
	return json.Marshal(data)
}

// TransformBytes transforms SAM template bytes to CloudFormation template bytes.
// This is a convenience function for direct byte-to-byte transformation.
func TransformBytes(input []byte, opts *TransformOptions) ([]byte, *SourceMap, error) {
	tmpl, err := template.Parse(input)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse template: %w", err)
	}

	result, err := Transform(tmpl, opts)
	if err != nil {
		return nil, nil, err
	}

	// Convert back to YAML
	output, err := yaml.Marshal(templateToMap(result.Template))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to serialize transformed template: %w", err)
	}

	return output, result.SourceMap, nil
}

// templateToMap converts a template to a map for serialization.
func templateToMap(tmpl *template.Template) map[string]any {
	data := make(map[string]any)

	if tmpl.AWSTemplateFormatVersion != "" {
		data["AWSTemplateFormatVersion"] = tmpl.AWSTemplateFormatVersion
	}
	if tmpl.Description != "" {
		data["Description"] = tmpl.Description
	}

	if len(tmpl.Resources) > 0 {
		resources := make(map[string]any)
		for name, r := range tmpl.Resources {
			res := make(map[string]any)
			res["Type"] = r.Type
			if len(r.Properties) > 0 {
				res["Properties"] = r.Properties
			}
			resources[name] = res
		}
		data["Resources"] = resources
	}

	return data
}
