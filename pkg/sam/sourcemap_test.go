package sam

import (
	"testing"
)

func TestNewSourceMap(t *testing.T) {
	sm := NewSourceMap()
	if sm == nil {
		t.Fatal("Expected NewSourceMap to return non-nil SourceMap")
	}
	if sm.ResourceMapping == nil {
		t.Error("Expected ResourceMapping to be initialized")
	}
}

func TestSourceMap_AddResourceMapping(t *testing.T) {
	sm := NewSourceMap()
	sm.AddResourceMapping("MyFunctionRole", "MyFunction", 10, 5)

	loc, ok := sm.GetResourceLocation("MyFunctionRole")
	if !ok {
		t.Fatal("Expected to find resource mapping for MyFunctionRole")
	}
	if loc.OriginalResource != "MyFunction" {
		t.Errorf("Expected OriginalResource 'MyFunction', got %q", loc.OriginalResource)
	}
	if loc.Line != 10 {
		t.Errorf("Expected Line 10, got %d", loc.Line)
	}
	if loc.Column != 5 {
		t.Errorf("Expected Column 5, got %d", loc.Column)
	}
}

func TestSourceMap_MapError_ExactMatch(t *testing.T) {
	sm := NewSourceMap()
	sm.AddResourceMapping("MyFunctionRole", "MyFunction", 10, 5)

	mappedLoc := sm.MapError("MyFunctionRole", "", 100)

	if mappedLoc.Line != 10 {
		t.Errorf("Expected mapped line 10, got %d", mappedLoc.Line)
	}
	if mappedLoc.Column != 5 {
		t.Errorf("Expected mapped column 5, got %d", mappedLoc.Column)
	}
	if mappedLoc.OriginalResource != "MyFunction" {
		t.Errorf("Expected OriginalResource 'MyFunction', got %q", mappedLoc.OriginalResource)
	}
}

func TestSourceMap_MapError_NoMapping(t *testing.T) {
	sm := NewSourceMap()

	// When no mapping exists, should return the original CFN line
	mappedLoc := sm.MapError("UnknownResource", "", 50)

	if mappedLoc.Line != 50 {
		t.Errorf("Expected fallback to original line 50, got %d", mappedLoc.Line)
	}
}

func TestSourceMap_MapError_OriginalResourcePreserved(t *testing.T) {
	sm := NewSourceMap()
	// Map a generated Lambda function to its source SAM function
	sm.AddResourceMapping("MyFunctionLambda", "MyFunction", 15, 3)

	mappedLoc := sm.MapError("MyFunctionLambda", "Properties.Code", 200)

	if mappedLoc.OriginalResource != "MyFunction" {
		t.Errorf("Expected OriginalResource 'MyFunction', got %q", mappedLoc.OriginalResource)
	}
}

func TestSourceLocation_String(t *testing.T) {
	loc := SourceLocation{
		Line:             10,
		Column:           5,
		OriginalResource: "MyFunction",
	}

	str := loc.String()
	if str != "MyFunction:10:5" {
		t.Errorf("Expected 'MyFunction:10:5', got %q", str)
	}
}

func TestSourceMap_MultipleResources(t *testing.T) {
	sm := NewSourceMap()

	// SAM Function expands to multiple CFN resources
	sm.AddResourceMapping("MyFunctionRole", "MyFunction", 10, 3)
	sm.AddResourceMapping("MyFunctionLambda", "MyFunction", 10, 3)
	sm.AddResourceMapping("MyFunctionApiPermission", "MyFunction", 10, 3)

	// All should map back to the same SAM resource
	for _, cfnResource := range []string{"MyFunctionRole", "MyFunctionLambda", "MyFunctionApiPermission"} {
		loc := sm.MapError(cfnResource, "", 100)
		if loc.OriginalResource != "MyFunction" {
			t.Errorf("Expected %s to map to MyFunction, got %s", cfnResource, loc.OriginalResource)
		}
		if loc.Line != 10 {
			t.Errorf("Expected %s to map to line 10, got %d", cfnResource, loc.Line)
		}
	}
}
