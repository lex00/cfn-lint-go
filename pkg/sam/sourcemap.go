package sam

import "fmt"

// SourceMap tracks the relationship between transformed CloudFormation
// resources and their original SAM template locations.
type SourceMap struct {
	// ResourceMapping maps CloudFormation resource names to SAM source locations.
	ResourceMapping map[string]SourceLocation
}

// SourceLocation represents a position in the original SAM template.
type SourceLocation struct {
	// Line is the 1-indexed line number in the original SAM template.
	Line int

	// Column is the 1-indexed column number in the original SAM template.
	Column int

	// OriginalResource is the name of the SAM resource this was derived from.
	OriginalResource string
}

// NewSourceMap creates a new empty SourceMap.
func NewSourceMap() *SourceMap {
	return &SourceMap{
		ResourceMapping: make(map[string]SourceLocation),
	}
}

// AddResourceMapping adds a mapping from a CloudFormation resource to its
// original SAM resource location.
func (sm *SourceMap) AddResourceMapping(cfnResource, samResource string, line, column int) {
	sm.ResourceMapping[cfnResource] = SourceLocation{
		Line:             line,
		Column:           column,
		OriginalResource: samResource,
	}
}

// GetResourceLocation returns the source location for a CloudFormation resource.
func (sm *SourceMap) GetResourceLocation(cfnResource string) (SourceLocation, bool) {
	loc, ok := sm.ResourceMapping[cfnResource]
	return loc, ok
}

// MapError translates an error location from the transformed template
// back to the original SAM template.
func (sm *SourceMap) MapError(cfnResource, propertyPath string, cfnLine int) SourceLocation {
	// Try resource-level mapping
	if loc, ok := sm.ResourceMapping[cfnResource]; ok {
		return loc
	}

	// No mapping found, return original line with empty original resource
	return SourceLocation{Line: cfnLine}
}

// String returns a string representation of the source location.
func (loc SourceLocation) String() string {
	return fmt.Sprintf("%s:%d:%d", loc.OriginalResource, loc.Line, loc.Column)
}
