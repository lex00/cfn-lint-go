// Package informational contains informational-level rules (Ixxx).
package informational

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&I1002{})
}

// I1002 checks if the template is approaching the size limit.
type I1002 struct{}

func (r *I1002) ID() string { return "I1002" }

func (r *I1002) ShortDesc() string {
	return "Template size approaching limit"
}

func (r *I1002) Description() string {
	return "Checks if the template size is approaching the CloudFormation limit of 51,200 bytes. This rule provides an informational message when the template exceeds 80% of the limit (40,960 bytes)."
}

func (r *I1002) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cloudformation-limits.html"
}

func (r *I1002) Tags() []string {
	return []string{"template", "limits", "size"}
}

// MaxTemplateSize is the CloudFormation limit for template size in bytes (51,200 bytes).
const MaxTemplateSize = 51200

// TemplateSizeWarningThreshold is 80% of the maximum template size.
const TemplateSizeWarningThreshold = int(float64(MaxTemplateSize) * 0.8)

func (r *I1002) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	// Get the template size from the filename if available
	// For now, we'll use the serialized YAML size as a proxy
	if tmpl.Root == nil {
		return matches
	}

	// Estimate size by marshaling the root node
	// This is an approximation of the actual template size
	templateSize := estimateTemplateSize(tmpl)

	if templateSize > TemplateSizeWarningThreshold {
		percentage := int(float64(templateSize) / float64(MaxTemplateSize) * 100)
		matches = append(matches, rules.Match{
			Message: fmt.Sprintf("Template size is %d bytes (%d%% of %d byte limit). Consider splitting into nested stacks or reducing content.", templateSize, percentage, MaxTemplateSize),
			Line:    1,
			Column:  1,
			Path:    []string{},
		})
	}

	return matches
}

// estimateTemplateSize estimates the size of the template by counting YAML nodes.
// This is a rough approximation; actual file size may vary.
func estimateTemplateSize(tmpl *template.Template) int {
	// Simple heuristic: estimate based on number of resources, parameters, outputs, etc.
	// Each resource/parameter/output typically adds ~100-500 bytes
	size := 0

	// Base template overhead
	size += 100

	// Parameters
	size += len(tmpl.Parameters) * 150

	// Resources (typically larger)
	size += len(tmpl.Resources) * 300

	// Outputs
	size += len(tmpl.Outputs) * 150

	// Mappings
	size += len(tmpl.Mappings) * 200

	// Conditions
	for name, cond := range tmpl.Conditions {
		size += len(name) + 100
		if cond.Node != nil {
			size += 50
		}
	}

	// Description
	if tmpl.Description != "" {
		size += len(tmpl.Description)
	}

	return size
}
