package rules

import (
	"testing"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

// mockRule is a simple rule for testing the registry.
type mockRule struct {
	id string
}

func (r *mockRule) ID() string                            { return r.id }
func (r *mockRule) ShortDesc() string                     { return "Mock rule" }
func (r *mockRule) Description() string                   { return "A mock rule for testing" }
func (r *mockRule) Source() string                        { return "" }
func (r *mockRule) Tags() []string                        { return []string{"test"} }
func (r *mockRule) Match(tmpl *template.Template) []Match { return nil }

func TestRegister(t *testing.T) {
	// Clear registry for test
	registryLock.Lock()
	registry = nil
	registryLock.Unlock()

	Register(&mockRule{id: "TEST001"})
	Register(&mockRule{id: "TEST002"})

	if Count() != 2 {
		t.Errorf("Expected 2 rules, got %d", Count())
	}
}

func TestGet(t *testing.T) {
	// Clear and add test rules
	registryLock.Lock()
	registry = nil
	registryLock.Unlock()

	Register(&mockRule{id: "TEST001"})
	Register(&mockRule{id: "TEST002"})

	rule := Get("TEST001")
	if rule == nil {
		t.Error("Expected to find TEST001")
	}
	if rule.ID() != "TEST001" {
		t.Errorf("Expected ID TEST001, got %s", rule.ID())
	}

	notFound := Get("NOTEXIST")
	if notFound != nil {
		t.Error("Expected nil for non-existent rule")
	}
}

func TestAll(t *testing.T) {
	// Clear and add test rules
	registryLock.Lock()
	registry = nil
	registryLock.Unlock()

	Register(&mockRule{id: "TEST001"})
	Register(&mockRule{id: "TEST002"})

	all := All()
	if len(all) != 2 {
		t.Errorf("Expected 2 rules, got %d", len(all))
	}

	// Verify it's a copy
	all[0] = nil
	if Get("TEST001") == nil {
		t.Error("All() should return a copy, not the original slice")
	}
}
