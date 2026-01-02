// Package schema provides CloudFormation resource schema validation.
//
// This package wraps cloudformation-schema-go/spec to provide convenient access
// to the CloudFormation Resource Specification for validating resource types,
// properties, and attributes.
//
// # Loading the Schema
//
// The schema is lazily loaded and cached:
//
//	spec, err := schema.Load()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Force a refresh:
//
//	spec, err := schema.LoadWithOptions(&schema.Options{Force: true})
//
// # Validating Resource Types
//
// Check if a resource type exists:
//
//	exists, _ := schema.HasResourceType("AWS::S3::Bucket")
//
// Get required properties:
//
//	required, _ := schema.GetRequiredProperties("AWS::Lambda::Function")
//	// returns: []string{"Code", "Role"}
//
// # Validating Properties
//
// Check if a property exists:
//
//	exists, _ := schema.HasProperty("AWS::S3::Bucket", "BucketName")
//
// Get property definition:
//
//	prop, _ := schema.GetProperty("AWS::Lambda::Function", "Runtime")
//	if prop != nil {
//	    fmt.Println(prop.PrimitiveType) // "String"
//	}
//
// # Validating Attributes (for GetAtt)
//
// Check if an attribute is valid for GetAtt:
//
//	exists, _ := schema.HasAttribute("AWS::S3::Bucket", "Arn")
//	exists, _ := schema.HasAttribute("AWS::S3::Bucket", "InvalidAttr") // false
package schema
