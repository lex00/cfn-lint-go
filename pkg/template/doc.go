// Package template provides CloudFormation template parsing with line number tracking.
//
// Templates are parsed into a structured representation that preserves YAML node
// information for accurate line number reporting in linting errors.
//
// # Parsing Templates
//
// Parse from a file:
//
//	tmpl, err := template.ParseFile("template.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Parse from bytes:
//
//	tmpl, err := template.Parse([]byte(yamlContent))
//
// # Accessing Template Sections
//
// The Template struct provides access to all CloudFormation sections:
//
//	// Resources
//	for name, res := range tmpl.Resources {
//	    fmt.Printf("%s: %s\n", name, res.Type)
//	}
//
//	// Parameters
//	for name, param := range tmpl.Parameters {
//	    fmt.Printf("%s: %s (default: %v)\n", name, param.Type, param.Default)
//	}
//
//	// Outputs, Mappings, Conditions are also available
//
// # Line Number Tracking
//
// Each parsed element includes a Node field with YAML position information:
//
//	res := tmpl.Resources["MyBucket"]
//	line := res.Node.Line
//	column := res.Node.Column
//
// # Intrinsic Functions
//
// CloudFormation intrinsic functions are parsed into their long-form map representation:
//
//	// YAML: BucketName: !Ref MyParam
//	// Parsed as: map[string]any{"Ref": "MyParam"}
//
// Supported intrinsic tags: !Ref, !GetAtt, !Sub, !Join, !Select, !If, !Condition,
// !GetAZs, !Base64, !Cidr, !FindInMap, !ImportValue, !Split, and more.
package template
