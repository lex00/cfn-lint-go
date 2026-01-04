# Custom Rules

This guide explains how to create and use custom rules with cfn-lint-go.

## Overview

cfn-lint-go allows you to create custom rules that are specific to your organization's needs. Custom rules can enforce:

- Naming conventions
- Required tags
- Allowed resource types
- Security policies
- Cost optimization rules
- And more

## Creating Custom Rules as a Library

The most flexible approach is to create custom rules as a Go library.

### Step 1: Create Your Rule Package

```go
// mycompany/cfnrules/naming.go
package cfnrules

import (
    "fmt"
    "strings"

    "github.com/lex00/cfn-lint-go/pkg/rules"
    "github.com/lex00/cfn-lint-go/pkg/template"
)

// C0001 enforces company naming conventions.
type C0001 struct{}

func (r *C0001) ID() string          { return "C0001" }
func (r *C0001) ShortDesc() string   { return "Resource naming convention" }
func (r *C0001) Description() string { return "Resources must follow company naming convention" }
func (r *C0001) Source() string      { return "https://wiki.mycompany.com/cfn-naming" }
func (r *C0001) Tags() []string      { return []string{"custom", "naming"} }

func (r *C0001) Match(tmpl *template.Template) []rules.Match {
    var matches []rules.Match

    for name := range tmpl.Resources {
        // Example: require resources to start with company prefix
        if !strings.HasPrefix(name, "Acme") {
            matches = append(matches, rules.Match{
                Message: fmt.Sprintf("Resource '%s' must start with 'Acme' prefix", name),
                Path:    []string{"Resources", name},
            })
        }
    }

    return matches
}
```

### Step 2: Create a Custom CLI

```go
// cmd/acme-cfn-lint/main.go
package main

import (
    "fmt"
    "os"

    "github.com/lex00/cfn-lint-go/pkg/lint"

    // Import standard rules
    _ "github.com/lex00/cfn-lint-go/internal/rules/conditions"
    _ "github.com/lex00/cfn-lint-go/internal/rules/errors"
    _ "github.com/lex00/cfn-lint-go/internal/rules/functions"
    // ... other standard rules

    // Import custom rules
    _ "mycompany/cfnrules"
)

func main() {
    linter := lint.New(lint.Options{})

    matches, err := linter.LintFile(os.Args[1])
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    for _, m := range matches {
        fmt.Printf("%s: %s [%s]\n", m.Location.Filename, m.Message, m.Rule.ID)
    }

    if len(matches) > 0 {
        os.Exit(2)
    }
}
```

### Step 3: Register Custom Rules

In your custom rules package, register rules in `init()`:

```go
// mycompany/cfnrules/init.go
package cfnrules

import "github.com/lex00/cfn-lint-go/pkg/rules"

func init() {
    rules.Register(&C0001{})
    rules.Register(&C0002{})
    // ... more custom rules
}
```

## Using cfn-lint-go as a Library

You can also use cfn-lint-go programmatically without creating a custom CLI:

```go
package main

import (
    "fmt"

    "github.com/lex00/cfn-lint-go/pkg/lint"
    "github.com/lex00/cfn-lint-go/pkg/rules"
    "github.com/lex00/cfn-lint-go/pkg/template"

    // Import standard rules
    _ "github.com/lex00/cfn-lint-go/internal/rules/errors"
    _ "github.com/lex00/cfn-lint-go/internal/rules/functions"
)

// Custom rule implemented inline
type RequiredTagsRule struct{}

func (r *RequiredTagsRule) ID() string          { return "C0002" }
func (r *RequiredTagsRule) ShortDesc() string   { return "Required tags" }
func (r *RequiredTagsRule) Description() string { return "Resources must have required tags" }
func (r *RequiredTagsRule) Source() string      { return "" }
func (r *RequiredTagsRule) Tags() []string      { return []string{"custom", "tags"} }

func (r *RequiredTagsRule) Match(tmpl *template.Template) []rules.Match {
    var matches []rules.Match
    requiredTags := []string{"Environment", "Owner", "CostCenter"}

    for name, res := range tmpl.Resources {
        props, ok := res.Properties.(map[string]interface{})
        if !ok {
            continue
        }

        tags, ok := props["Tags"].([]interface{})
        if !ok {
            matches = append(matches, rules.Match{
                Message: fmt.Sprintf("Resource '%s' is missing required tags", name),
                Path:    []string{"Resources", name},
            })
            continue
        }

        tagKeys := make(map[string]bool)
        for _, tag := range tags {
            if tagMap, ok := tag.(map[string]interface{}); ok {
                if key, ok := tagMap["Key"].(string); ok {
                    tagKeys[key] = true
                }
            }
        }

        for _, required := range requiredTags {
            if !tagKeys[required] {
                matches = append(matches, rules.Match{
                    Message: fmt.Sprintf("Resource '%s' is missing required tag '%s'", name, required),
                    Path:    []string{"Resources", name, "Properties", "Tags"},
                })
            }
        }
    }

    return matches
}

func main() {
    // Register custom rule
    rules.Register(&RequiredTagsRule{})

    // Create linter
    linter := lint.New(lint.Options{})

    // Lint template
    matches, err := linter.LintFile("template.yaml")
    if err != nil {
        panic(err)
    }

    for _, m := range matches {
        fmt.Printf("%s [%s]\n", m.Message, m.Rule.ID)
    }
}
```

## Common Custom Rule Examples

### Enforce Resource Limits

```go
type MaxResourcesRule struct {
    MaxResources int
}

func (r *MaxResourcesRule) Match(tmpl *template.Template) []rules.Match {
    if len(tmpl.Resources) > r.MaxResources {
        return []rules.Match{{
            Message: fmt.Sprintf("Template has %d resources (max: %d)",
                len(tmpl.Resources), r.MaxResources),
        }}
    }
    return nil
}
```

### Deny Specific Resource Types

```go
type DenyResourceTypesRule struct {
    DeniedTypes []string
}

func (r *DenyResourceTypesRule) Match(tmpl *template.Template) []rules.Match {
    var matches []rules.Match
    denied := make(map[string]bool)
    for _, t := range r.DeniedTypes {
        denied[t] = true
    }

    for name, res := range tmpl.Resources {
        if denied[res.Type] {
            matches = append(matches, rules.Match{
                Message: fmt.Sprintf("Resource type '%s' is not allowed", res.Type),
                Path:    []string{"Resources", name, "Type"},
            })
        }
    }
    return matches
}
```

### Require Encryption

```go
type RequireEncryptionRule struct{}

func (r *RequireEncryptionRule) Match(tmpl *template.Template) []rules.Match {
    var matches []rules.Match

    for name, res := range tmpl.Resources {
        props, ok := res.Properties.(map[string]interface{})
        if !ok {
            continue
        }

        switch res.Type {
        case "AWS::S3::Bucket":
            if _, ok := props["BucketEncryption"]; !ok {
                matches = append(matches, rules.Match{
                    Message: fmt.Sprintf("S3 bucket '%s' must have encryption enabled", name),
                    Path:    []string{"Resources", name},
                })
            }
        case "AWS::RDS::DBInstance":
            if encrypted, ok := props["StorageEncrypted"].(bool); !ok || !encrypted {
                matches = append(matches, rules.Match{
                    Message: fmt.Sprintf("RDS instance '%s' must have storage encryption", name),
                    Path:    []string{"Resources", name},
                })
            }
        }
    }

    return matches
}
```

## Filtering Rules

Use the linter options to control which rules run:

```go
linter := lint.New(lint.Options{
    // Ignore specific rules
    IgnoreRules: []string{"E1001", "W3002"},

    // Only run specific rules (not yet implemented)
    // OnlyRules: []string{"C0001", "C0002"},
})
```

## Testing Custom Rules

```go
func TestRequiredTagsRule(t *testing.T) {
    tmpl, _ := template.ParseString(`
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      Tags:
        - Key: Environment
          Value: prod
`)

    rule := &RequiredTagsRule{}
    matches := rule.Match(tmpl)

    // Should fail - missing Owner and CostCenter tags
    if len(matches) != 2 {
        t.Errorf("expected 2 matches, got %d", len(matches))
    }
}
```

## Distribution

### As a Go Module

Publish your custom rules as a Go module:

```bash
# In your rules repository
go mod init github.com/mycompany/cfn-lint-rules
git tag v0.1.0
git push origin v0.1.0
```

Users can then import:

```go
import _ "github.com/mycompany/cfn-lint-rules"
```

### As a Custom Binary

Build and distribute a custom binary with your rules baked in:

```bash
go build -o acme-cfn-lint ./cmd/acme-cfn-lint
```

## Limitations

Currently, cfn-lint-go does not support:

- Dynamic rule loading at runtime (`--append-rules`)
- Plugin architecture for rules
- Configuration-based rule creation

These features may be added in future versions. For now, custom rules must be compiled into the binary.

## See Also

- [Rule Creation Guide](getting_started/rules.md) - Technical details on implementing rules
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contributing rules to the main project
- [Rules Reference](RULES.md) - All built-in rules
