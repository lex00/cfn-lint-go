// Package rules defines the rule interface and registry.
package rules

import (
	"sync"

	"github.com/lex00/cfn-lint-go/pkg/template"
)

// Rule is the interface that all linting rules must implement.
type Rule interface {
	// ID returns the rule identifier (e.g., "E1001", "W3002").
	ID() string

	// ShortDesc returns a brief description of what the rule checks.
	ShortDesc() string

	// Description returns a detailed description of the rule.
	Description() string

	// Source returns a URL to documentation about this rule.
	Source() string

	// Tags returns searchable tags for the rule.
	Tags() []string

	// Match checks the template and returns any matches.
	Match(tmpl *template.Template) []Match
}

// Match represents a single rule violation.
type Match struct {
	Message string
	Line    int
	Column  int
	Path    []string // JSON path to the problematic element
}

// registry holds all registered rules.
var (
	registry     []Rule
	registryLock sync.RWMutex
)

// Register adds a rule to the global registry.
// This is typically called from init() functions.
func Register(r Rule) {
	registryLock.Lock()
	defer registryLock.Unlock()
	registry = append(registry, r)
}

// All returns all registered rules.
func All() []Rule {
	registryLock.RLock()
	defer registryLock.RUnlock()
	result := make([]Rule, len(registry))
	copy(result, registry)
	return result
}

// Get returns a rule by ID, or nil if not found.
func Get(id string) Rule {
	registryLock.RLock()
	defer registryLock.RUnlock()
	for _, r := range registry {
		if r.ID() == id {
			return r
		}
	}
	return nil
}

// Count returns the number of registered rules.
func Count() int {
	registryLock.RLock()
	defer registryLock.RUnlock()
	return len(registry)
}
