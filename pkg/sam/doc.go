// Package sam provides SAM (Serverless Application Model) template detection
// and transformation capabilities for cfn-lint-go.
//
// # SAM Detection
//
// SAM templates are identified by the presence of:
//   - Transform: AWS::Serverless-2016-10-31
//   - Any AWS::Serverless::* resource type
//
// Supported SAM resource types:
//   - AWS::Serverless::Function
//   - AWS::Serverless::Api
//   - AWS::Serverless::HttpApi
//   - AWS::Serverless::SimpleTable
//   - AWS::Serverless::LayerVersion
//   - AWS::Serverless::Application
//   - AWS::Serverless::StateMachine
//   - AWS::Serverless::Connector
//   - AWS::Serverless::GraphQLApi
//
// Example usage:
//
//	tmpl, err := template.ParseFile("template.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	if sam.IsSAMTemplate(tmpl) {
//	    fmt.Println("This is a SAM template")
//	}
package sam
