// Package rules defines the rule interface and registry for cfn-lint-go.
//
// All linting rules implement the Rule interface and register themselves
// with the global registry at init time.
//
// # Rule Interface
//
// Rules must implement:
//
//	type Rule interface {
//	    ID() string                              // e.g., "E1001"
//	    ShortDesc() string                       // Brief description
//	    Description() string                     // Detailed description
//	    Source() string                          // Documentation URL
//	    Tags() []string                          // Searchable tags
//	    Match(*template.Template) []Match        // Check template
//	}
//
// # Creating a Custom Rule
//
//	type NoHardcodedBuckets struct{}
//
//	func init() {
//	    rules.Register(&NoHardcodedBuckets{})
//	}
//
//	func (r *NoHardcodedBuckets) ID() string { return "C0001" }
//	func (r *NoHardcodedBuckets) ShortDesc() string { return "No hardcoded bucket names" }
//	// ... implement other methods
//
//	func (r *NoHardcodedBuckets) Match(tmpl *template.Template) []rules.Match {
//	    var matches []rules.Match
//	    for name, res := range tmpl.Resources {
//	        if res.Type == "AWS::S3::Bucket" {
//	            if bn, ok := res.Properties["BucketName"].(string); ok {
//	                matches = append(matches, rules.Match{
//	                    Message: "Hardcoded bucket name: " + bn,
//	                    Path:    []string{"Resources", name, "Properties", "BucketName"},
//	                })
//	            }
//	        }
//	    }
//	    return matches
//	}
//
// # Registry Functions
//
// Query registered rules:
//
//	allRules := rules.All()           // Get all rules
//	rule := rules.Get("E1001")        // Get by ID
//	count := rules.Count()            // Count registered rules
//
// # Rule ID Convention
//
// Rule IDs follow the Python cfn-lint convention:
//   - E0xxx: Parse/transform errors
//   - E1xxx: Intrinsic function errors
//   - E2xxx: Parameter errors
//   - E3xxx: Resource/property errors
//   - E4xxx: Metadata errors
//   - E6xxx: Output errors
//   - E7xxx: Mapping errors
//   - E8xxx: Condition errors
//   - W####: Warnings
//   - I####: Informational
package rules
