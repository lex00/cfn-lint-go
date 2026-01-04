# Rule Creation Guide

This guide explains how to create new rules for cfn-lint-go.

## Rule Interface

All rules implement the `rules.Rule` interface defined in `pkg/rules/rule.go`:

```go
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
```

## Rule Naming Convention

Rules follow a naming convention from the Python cfn-lint project:

| Prefix | Category | Description |
|--------|----------|-------------|
| E0xxx | Template Errors | Parse and structure errors |
| E1xxx | Functions | Intrinsic function validation |
| E2xxx | Parameters | Parameter validation |
| E3xxx | Resources | Resource property validation |
| E4xxx | Metadata | Metadata validation |
| E5xxx | Modules | CloudFormation Modules |
| E6xxx | Outputs | Output validation |
| E7xxx | Mappings | Mapping validation |
| E8xxx | Conditions | Condition validation |
| W1xxx | Template Warnings | Template best practices |
| W2xxx | Parameter Warnings | Parameter best practices |
| W3xxx | Resource Warnings | Resource best practices |
| I1xxx | Template Info | Informational rules |
| I2xxx | Parameter Info | Parameter suggestions |
| I3xxx | Resource Info | Resource suggestions |

## Creating a New Rule

### 1. Find the Python Implementation

Check the Python cfn-lint implementation for reference:

```
https://github.com/aws-cloudformation/cfn-lint/tree/main/src/cfnlint/rules
```

### 2. Create the Go File

Create a new file in the appropriate category directory:

```
internal/rules/{category}/{rule_id}.go
```

### 3. Implement the Rule

Here's a complete example for rule E1001 (Ref to undefined resource):

```go
package functions

import (
    "fmt"

    "github.com/lex00/cfn-lint-go/pkg/rules"
    "github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
    rules.Register(&E1001{})
}

// E1001 validates that Ref functions reference defined resources or parameters.
type E1001 struct{}

func (r *E1001) ID() string { return "E1001" }

func (r *E1001) ShortDesc() string {
    return "Ref to undefined resource or parameter"
}

func (r *E1001) Description() string {
    return "Validates that the Ref intrinsic function references a defined resource, parameter, or pseudo-parameter"
}

func (r *E1001) Source() string {
    return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-ref.html"
}

func (r *E1001) Tags() []string {
    return []string{"functions", "ref"}
}

func (r *E1001) Match(tmpl *template.Template) []rules.Match {
    var matches []rules.Match

    // Get all valid reference targets
    validRefs := make(map[string]bool)

    // Add parameters
    for name := range tmpl.Parameters {
        validRefs[name] = true
    }

    // Add resources
    for name := range tmpl.Resources {
        validRefs[name] = true
    }

    // Add pseudo-parameters
    pseudoParams := []string{
        "AWS::AccountId", "AWS::NotificationARNs", "AWS::NoValue",
        "AWS::Partition", "AWS::Region", "AWS::StackId",
        "AWS::StackName", "AWS::URLSuffix",
    }
    for _, p := range pseudoParams {
        validRefs[p] = true
    }

    // Walk template looking for Ref functions
    template.Walk(tmpl, func(path []string, value interface{}) {
        m, ok := value.(map[string]interface{})
        if !ok {
            return
        }

        ref, ok := m["Ref"]
        if !ok {
            return
        }

        refStr, ok := ref.(string)
        if !ok {
            return
        }

        if !validRefs[refStr] {
            matches = append(matches, rules.Match{
                Message: fmt.Sprintf("Ref '%s' references undefined resource or parameter", refStr),
                Path:    path,
            })
        }
    })

    return matches
}
```

### 4. Add Tests

Create a test file `{rule_id}_test.go`:

```go
package functions

import (
    "testing"

    "github.com/lex00/cfn-lint-go/pkg/template"
)

func TestE1001(t *testing.T) {
    tests := []struct {
        name     string
        template string
        wantErr  bool
        errCount int
    }{
        {
            name: "valid ref to parameter",
            template: `
Parameters:
  MyParam:
    Type: String
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !Ref MyParam
`,
            wantErr:  false,
            errCount: 0,
        },
        {
            name: "invalid ref to undefined",
            template: `
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !Ref UndefinedParam
`,
            wantErr:  true,
            errCount: 1,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tmpl, err := template.ParseString(tt.template)
            if err != nil {
                t.Fatalf("failed to parse template: %v", err)
            }

            rule := &E1001{}
            matches := rule.Match(tmpl)

            if tt.wantErr && len(matches) == 0 {
                t.Error("expected matches but got none")
            }
            if !tt.wantErr && len(matches) > 0 {
                t.Errorf("expected no matches but got %d", len(matches))
            }
            if len(matches) != tt.errCount {
                t.Errorf("expected %d matches but got %d", tt.errCount, len(matches))
            }
        })
    }
}
```

### 5. Add Test Fixtures (Optional)

For complex rules, add YAML test fixtures in `testdata/`:

```
testdata/
├── e1001_valid.yaml
├── e1001_invalid.yaml
└── ...
```

### 6. Register the Rule

Rules register themselves via `init()`:

```go
func init() {
    rules.Register(&E1001{})
}
```

Ensure the package is imported in `cmd/cfn-lint/main.go`:

```go
import (
    _ "github.com/lex00/cfn-lint-go/internal/rules/functions"
)
```

### 7. Update Documentation

Run the documentation generator to update RULES.md:

```bash
cfn-lint update-documentation
```

## Working with Templates

### Template Structure

The `template.Template` struct provides access to all template sections:

```go
type Template struct {
    AWSTemplateFormatVersion string
    Description              string
    Metadata                 map[string]interface{}
    Parameters               map[string]Parameter
    Mappings                 map[string]interface{}
    Conditions               map[string]interface{}
    Transform                interface{}
    Resources                map[string]Resource
    Outputs                  map[string]Output
    Rules                    map[string]interface{}
}
```

### Walking the Template

Use `template.Walk` to traverse all values:

```go
template.Walk(tmpl, func(path []string, value interface{}) {
    // path: ["Resources", "MyBucket", "Properties", "BucketName"]
    // value: the actual value at that path
})
```

### Getting Line Numbers

Use `template.GetLineNumber` for error locations:

```go
line := template.GetLineNumber(tmpl, path)
matches = append(matches, rules.Match{
    Message: "error message",
    Line:    line,
    Path:    path,
})
```

## Best Practices

1. **Match Python behavior**: When possible, match the error messages and behavior of Python cfn-lint
2. **Test edge cases**: Include tests for empty values, nested structures, and intrinsic functions
3. **Use descriptive messages**: Error messages should help users understand and fix the issue
4. **Tag appropriately**: Tags help users filter rules by category
5. **Document sources**: Include links to AWS documentation when applicable
6. **Handle intrinsic functions**: Many values can be intrinsic functions - handle them appropriately

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/rules/functions/...

# Run tests with verbose output
go test -v ./internal/rules/functions/...

# Run a specific test
go test -v -run TestE1001 ./internal/rules/functions/...
```

## Linting

Run the linter before submitting:

```bash
golangci-lint run
```

## Next Steps

- See [CONTRIBUTING.md](../CONTRIBUTING.md) for the full contribution workflow
- Check the [Rules Reference](../RULES.md) for existing implementations
- Join discussions in GitHub Issues for rule ideas
