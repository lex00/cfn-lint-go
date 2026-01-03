package resources

import (
	"testing"
)

func TestE3025_Metadata(t *testing.T) {
	rule := &E3025{}

	if rule.ID() != "E3025" {
		t.Errorf("Expected ID E3025, got %s", rule.ID())
	}
	if rule.ShortDesc() == "" {
		t.Error("ShortDesc should not be empty")
	}
	if len(rule.Tags()) == 0 {
		t.Error("Tags should not be empty")
	}
}
