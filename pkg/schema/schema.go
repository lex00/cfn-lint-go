// Package schema provides CloudFormation resource schema validation.
// It wraps the cloudformation-schema-go/spec package for cfn-lint-go.
package schema

import (
	"sync"

	"github.com/lex00/cloudformation-schema-go/spec"
)

var (
	// globalSpec is the lazily loaded CloudFormation spec.
	globalSpec     *spec.Spec
	globalSpecErr  error
	globalSpecOnce sync.Once
)

// Options configures how the schema is loaded.
type Options struct {
	// Force re-download even if cached.
	Force bool
	// Quiet suppresses progress output.
	Quiet bool
}

// Load fetches and returns the CloudFormation resource specification.
// The spec is cached after the first call. Use LoadWithOptions to force refresh.
func Load() (*spec.Spec, error) {
	return LoadWithOptions(nil)
}

// LoadWithOptions fetches the CloudFormation spec with custom options.
// If opts is nil, default options are used.
func LoadWithOptions(opts *Options) (*spec.Spec, error) {
	if opts != nil && opts.Force {
		// Force refresh - reset the once guard
		globalSpecOnce = sync.Once{}
	}

	globalSpecOnce.Do(func() {
		fetchOpts := &spec.FetchOptions{
			Quiet: true, // Default to quiet for linting
		}
		if opts != nil {
			fetchOpts.Force = opts.Force
			fetchOpts.Quiet = opts.Quiet
		}
		globalSpec, globalSpecErr = spec.FetchSpec(fetchOpts)
	})

	return globalSpec, globalSpecErr
}

// GetRequiredProperties returns the required property names for a resource type.
// Returns nil if the resource type is not found in the spec.
func GetRequiredProperties(resourceType string) ([]string, error) {
	s, err := Load()
	if err != nil {
		return nil, err
	}

	rt := s.GetResourceType(resourceType)
	if rt == nil {
		return nil, nil // Unknown resource type
	}

	return rt.GetRequiredProperties(), nil
}

// HasResourceType returns true if the resource type exists in the spec.
func HasResourceType(resourceType string) (bool, error) {
	s, err := Load()
	if err != nil {
		return false, err
	}
	return s.HasResourceType(resourceType), nil
}

// GetResourceType returns the resource type definition.
// Returns nil if not found.
func GetResourceType(resourceType string) (*spec.ResourceType, error) {
	s, err := Load()
	if err != nil {
		return nil, err
	}
	return s.GetResourceType(resourceType), nil
}

// GetProperty returns the property definition for a resource type.
// Returns nil if the resource or property is not found.
func GetProperty(resourceType, propertyName string) (*spec.Property, error) {
	s, err := Load()
	if err != nil {
		return nil, err
	}

	rt := s.GetResourceType(resourceType)
	if rt == nil {
		return nil, nil
	}

	return rt.GetProperty(propertyName), nil
}

// HasProperty returns true if the resource type has the given property.
func HasProperty(resourceType, propertyName string) (bool, error) {
	s, err := Load()
	if err != nil {
		return false, err
	}

	rt := s.GetResourceType(resourceType)
	if rt == nil {
		return false, nil
	}

	return rt.HasProperty(propertyName), nil
}

// GetAttribute returns the attribute definition for GetAtt validation.
// Returns nil if not found.
func GetAttribute(resourceType, attributeName string) (*spec.Attribute, error) {
	s, err := Load()
	if err != nil {
		return nil, err
	}

	rt := s.GetResourceType(resourceType)
	if rt == nil {
		return nil, nil
	}

	return rt.GetAttribute(attributeName), nil
}

// HasAttribute returns true if the resource type has the given attribute.
func HasAttribute(resourceType, attributeName string) (bool, error) {
	s, err := Load()
	if err != nil {
		return false, err
	}

	rt := s.GetResourceType(resourceType)
	if rt == nil {
		return false, nil
	}

	return rt.HasAttribute(attributeName), nil
}
