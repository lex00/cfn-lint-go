# SAM Template Support

cfn-lint-go provides native support for AWS SAM (Serverless Application Model) templates through integration with [aws-sam-translator-go](https://github.com/lex00/aws-sam-translator-go).

## How It Works

When you lint a SAM template, cfn-lint-go automatically:

1. **Detects** the template is a SAM template (via `AWS::Serverless-2016-10-31` transform or `AWS::Serverless::*` resource types)
2. **Transforms** the SAM template to CloudFormation using aws-sam-translator-go
3. **Lints** the resulting CloudFormation template with all 270 rules
4. **Maps** errors back to the original SAM template line numbers

This means you get full CloudFormation validation on your SAM templates without needing to run `sam build` first.

## Usage

### Basic Usage

SAM templates are automatically detected and transformed:

```bash
cfn-lint sam-template.yaml
```

### Skip Transformation

To lint SAM templates as-is without transformation:

```bash
cfn-lint sam-template.yaml --no-sam-transform
```

### Debug Transformation

To see the transformed CloudFormation template:

```bash
cfn-lint sam-template.yaml --show-transformed
```

## Configuration

### CLI Flags

| Flag | Description |
|------|-------------|
| `--no-sam-transform` | Skip SAM to CloudFormation transformation |
| `--show-transformed` | Output transformed CloudFormation template |

### Config File

Configure SAM behavior in `.cfnlintrc.yaml`:

```yaml
sam:
  # Disable automatic SAM transformation
  auto_transform: true

  # Configure transformation context
  transform_options:
    region: us-east-1
    account_id: "123456789012"
    stack_name: my-sam-app
    partition: aws  # aws, aws-cn, or aws-us-gov
```

## Supported SAM Resource Types

cfn-lint-go recognizes all SAM resource types:

| Resource Type | Description |
|---------------|-------------|
| `AWS::Serverless::Function` | Lambda function with simplified syntax |
| `AWS::Serverless::Api` | API Gateway REST API |
| `AWS::Serverless::HttpApi` | API Gateway HTTP API |
| `AWS::Serverless::SimpleTable` | DynamoDB table |
| `AWS::Serverless::LayerVersion` | Lambda layer |
| `AWS::Serverless::Application` | Nested serverless application |
| `AWS::Serverless::StateMachine` | Step Functions state machine |
| `AWS::Serverless::Connector` | Resource connector |
| `AWS::Serverless::GraphQLApi` | AppSync GraphQL API |

## SAM-Specific Rules

cfn-lint-go includes rules specifically for SAM templates:

| Rule | Level | Description |
|------|-------|-------------|
| E0010 | Error | SAM transform failed |
| W3100 | Warning | SAM Function missing MemorySize (default 128MB may not be optimal) |
| W3101 | Warning | SAM Function missing Timeout (default 3s may be too short) |
| W3102 | Warning | SAM Api missing StageName |
| I3101 | Informational | SAM resource expansion info |

## About aws-sam-translator-go

The SAM transformation is powered by [aws-sam-translator-go](https://github.com/lex00/aws-sam-translator-go), a native Go port of AWS's official [SAM Translator](https://github.com/aws/serverless-application-model).

### Key Features

- **Pure Go implementation** - No Python or external dependencies
- **Full SAM spec support** - Handles all SAM resource types and properties
- **Source mapping** - Tracks original line numbers through transformation
- **Configurable context** - Set region, account ID, stack name for accurate ARN generation

### How Transformation Works

```
SAM Template                    aws-sam-translator-go              CloudFormation
┌─────────────────┐            ┌────────────────────┐            ┌─────────────────┐
│ AWS::Serverless │            │                    │            │ AWS::Lambda::   │
│ ::Function      │ ────────▶  │    Transform()     │ ────────▶  │ Function        │
│                 │            │                    │            │ + IAM Role      │
│ Events:         │            │  Expands SAM       │            │ + Permissions   │
│   Api: ...      │            │  resources into    │            │ + API Gateway   │
└─────────────────┘            │  CloudFormation    │            └─────────────────┘
                               └────────────────────┘
```

A single `AWS::Serverless::Function` with an API event can expand into:
- `AWS::Lambda::Function`
- `AWS::IAM::Role`
- `AWS::Lambda::Permission`
- `AWS::ApiGateway::RestApi`
- `AWS::ApiGateway::Deployment`
- `AWS::ApiGateway::Stage`
- And more...

### Source Mapping

cfn-lint-go maintains a source map during transformation so that errors in the generated CloudFormation are reported at the correct line in your original SAM template:

```
sam-template.yaml:15:5: E3003 Required property 'Handler' missing
                        ^^^^^^
                        Points to original SAM template location,
                        not the generated CloudFormation
```

## Library Usage

Use SAM transformation in your Go code:

```go
package main

import (
    "fmt"
    "log"

    "github.com/lex00/cfn-lint-go/pkg/sam"
    "github.com/lex00/cfn-lint-go/pkg/template"
)

func main() {
    // Parse SAM template
    tmpl, err := template.ParseFile("sam-template.yaml")
    if err != nil {
        log.Fatal(err)
    }

    // Check if it's a SAM template
    if sam.IsSAMTemplate(tmpl) {
        fmt.Println("Detected SAM template")

        // Transform to CloudFormation
        result, err := sam.Transform(tmpl, sam.DefaultTransformOptions())
        if err != nil {
            log.Fatal(err)
        }

        // Use transformed template
        fmt.Printf("Transformed template has %d resources\n",
            len(result.Template.Resources))

        // Check for transformation warnings
        for _, w := range result.Warnings {
            fmt.Printf("Warning: %s\n", w)
        }
    }
}
```

## Limitations

- **Intrinsic functions in Globals**: Some complex intrinsic function usage in Globals may not transform correctly
- **Custom resource providers**: Third-party SAM extensions are not supported
- **SAM CLI plugins**: Plugin-based transformations are not available

## See Also

- [aws-sam-translator-go](https://github.com/lex00/aws-sam-translator-go) - The SAM translator library
- [AWS SAM Specification](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-specification.html) - Official SAM documentation
- [Rules Reference](RULES.md) - All available linting rules
