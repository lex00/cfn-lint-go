package schema

import (
	"testing"
)

func TestGetRequiredProperties(t *testing.T) {
	tests := []struct {
		resourceType string
		wantRequired []string // subset of expected required properties
		wantNil      bool
	}{
		{
			resourceType: "AWS::Lambda::Function",
			wantRequired: []string{"Code", "Role"},
		},
		{
			resourceType: "AWS::IAM::Role",
			wantRequired: []string{"AssumeRolePolicyDocument"},
		},
		{
			resourceType: "AWS::S3::Bucket",
			wantRequired: nil, // S3 bucket has no required properties
		},
		{
			resourceType: "AWS::NonExistent::Resource",
			wantNil:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.resourceType, func(t *testing.T) {
			required, err := GetRequiredProperties(tt.resourceType)
			if err != nil {
				t.Fatalf("GetRequiredProperties() error = %v", err)
			}

			if tt.wantNil {
				if required != nil {
					t.Errorf("GetRequiredProperties() = %v, want nil", required)
				}
				return
			}

			// Check that expected required properties are present
			requiredSet := make(map[string]bool)
			for _, r := range required {
				requiredSet[r] = true
			}

			for _, want := range tt.wantRequired {
				if !requiredSet[want] {
					t.Errorf("GetRequiredProperties() missing expected property %q, got %v", want, required)
				}
			}
		})
	}
}

func TestHasResourceType(t *testing.T) {
	tests := []struct {
		resourceType string
		want         bool
	}{
		{"AWS::S3::Bucket", true},
		{"AWS::Lambda::Function", true},
		{"AWS::EC2::Instance", true},
		{"AWS::NonExistent::Resource", false},
		{"NotAResource", false},
	}

	for _, tt := range tests {
		t.Run(tt.resourceType, func(t *testing.T) {
			got, err := HasResourceType(tt.resourceType)
			if err != nil {
				t.Fatalf("HasResourceType() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("HasResourceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasProperty(t *testing.T) {
	tests := []struct {
		resourceType string
		propertyName string
		want         bool
	}{
		{"AWS::S3::Bucket", "BucketName", true},
		{"AWS::S3::Bucket", "NonExistentProperty", false},
		{"AWS::Lambda::Function", "Runtime", true},
		{"AWS::Lambda::Function", "Handler", true},
		{"AWS::NonExistent::Resource", "Property", false},
	}

	for _, tt := range tests {
		t.Run(tt.resourceType+"/"+tt.propertyName, func(t *testing.T) {
			got, err := HasProperty(tt.resourceType, tt.propertyName)
			if err != nil {
				t.Fatalf("HasProperty() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("HasProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasAttribute(t *testing.T) {
	tests := []struct {
		resourceType  string
		attributeName string
		want          bool
	}{
		{"AWS::S3::Bucket", "Arn", true},
		{"AWS::S3::Bucket", "DomainName", true},
		{"AWS::S3::Bucket", "NonExistentAttr", false},
		{"AWS::Lambda::Function", "Arn", true},
		{"AWS::NonExistent::Resource", "Arn", false},
	}

	for _, tt := range tests {
		t.Run(tt.resourceType+"/"+tt.attributeName, func(t *testing.T) {
			got, err := HasAttribute(tt.resourceType, tt.attributeName)
			if err != nil {
				t.Fatalf("HasAttribute() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("HasAttribute() = %v, want %v", got, tt.want)
			}
		})
	}
}
