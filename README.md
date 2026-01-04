# cfn-lint-go

CloudFormation Linter for Go - a native Go port of [aws-cloudformation/cfn-lint](https://github.com/aws-cloudformation/cfn-lint).

## Status

**v1.0.0 - 270 rules implemented (100% feature parity with Python cfn-lint)**

This is a Go port of the Python cfn-lint tool. Implements complete feature parity with 270 rules covering template structure, intrinsic functions, schema-based validation, best practices, warnings, and informational rules. Uses `cloudformation-schema-go` for CloudFormation resource specification and enum validation.

### What's Implemented

- YAML/JSON template parsing with line number tracking
- CloudFormation intrinsic function support (!Ref, !GetAtt, !Sub, etc.)
- Rule interface and registry system
- DOT graph generation for resource dependencies
- **Multiple output formats**: text, JSON, SARIF, JUnit XML, pretty (colorized with context)
- **Configuration file support**: `.cfnlintrc`, `.cfnlintrc.yaml`, `.cfnlintrc.json`
- **SAM template support**: automatic transformation to CloudFormation with source mapping
- **GitHub Action** for CI/CD integration with SARIF support
- **Pre-commit hooks** for local validation
- CLI `graph` command for dependency visualization
- CLI `list-rules` command
- Complete CLI options matching Python cfn-lint
- 270 rules across all categories:
  - **E0xxx**: 7 rules (parse, transform, processing, config, SAM, deployment/parameter files)
  - **E1xxx**: 39 rules (intrinsic functions, schema validation, format validation)
  - **E2xxx**: 14 rules (param config, type, naming, length, limits, defaults, NoEcho, SSM types, constraints)
  - **E3xxx**: 119 rules (resource config, properties, type validation, enum validation, dependencies, policies, constraints)
  - **E4xxx**: 2 rules (interface metadata, structure)
  - **E5xxx**: 1 rule (CloudFormation Modules validation)
  - **E6xxx**: 9 rules (output structure, types, naming, exports)
  - **E7xxx**: 3 rules (mapping config, naming, limits)
  - **E8xxx**: 7 rules (condition functions)
  - **Wxxx**: 49 warning rules (security, best practices, deprecations)
  - **Ixxx**: 20 informational rules

## Installation

```bash
go install github.com/lex00/cfn-lint-go/cmd/cfn-lint@latest
```

Or add as a library:

```bash
go get github.com/lex00/cfn-lint-go
```

## Usage

### CLI

```bash
# Lint a template
cfn-lint template.yaml

# Lint with different output formats
cfn-lint template.yaml --format json
cfn-lint template.yaml --format sarif --output results.sarif
cfn-lint template.yaml --format junit --output results.xml
cfn-lint template.yaml --format pretty  # Colorized output with code context

# Use configuration file
cfn-lint template.yaml --config .cfnlintrc.yaml

# Ignore specific rules
cfn-lint template.yaml --ignore-rules E1001,W3002

# Include specific rules (even if ignored elsewhere)
cfn-lint template.yaml --ignore-rules E1001 --include-checks E1001

# Specify AWS regions
cfn-lint template.yaml --regions us-east-1,us-west-2

# Write output to file
cfn-lint template.yaml --output results.txt

# Generate dependency graph
cfn-lint graph template.yaml > deps.dot
dot -Tpng deps.dot -o deps.png

# Include parameters in graph
cfn-lint graph template.yaml --include-parameters

# List available rules
cfn-lint list-rules

# List rules as JSON
cfn-lint list-rules --format json

# Update RULES.md documentation
cfn-lint update-documentation

# SAM template support (auto-transforms to CloudFormation)
cfn-lint sam-template.yaml

# Skip SAM transformation (lint as-is)
cfn-lint sam-template.yaml --no-sam-transform

# Show transformed CloudFormation (for debugging)
cfn-lint sam-template.yaml --show-transformed

# Show help
cfn-lint --help
```

### Configuration File

Create a `.cfnlintrc.yaml` file in your project root:

```yaml
# Templates to lint (supports globs)
templates:
  - templates/**/*.yaml
  - infrastructure/*.yml

# Templates to ignore
ignore_templates:
  - test/**

# AWS regions to validate against
regions:
  - us-east-1
  - us-west-2

# Rules to ignore
ignore_checks:
  - E1001
  - W3002

# Rules to include (even if ignored)
include_checks:
  - I1001

# Output format
format: pretty

# Output file (optional)
# output_file: lint-results.json

# SAM template configuration
sam:
  auto_transform: true  # Set to false to skip SAM transformation
  transform_options:
    region: us-east-1
    account_id: "123456789012"
    stack_name: my-sam-app
```

### GitHub Actions

```yaml
name: Validate CloudFormation

on: [pull_request]

permissions:
  security-events: write
  contents: read

jobs:
  cfn-lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - uses: lex00/cfn-lint-go@main
        with:
          templates: 'templates/*.yaml'
          format: sarif
```

### Pre-commit Hooks

Add to `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: https://github.com/lex00/cfn-lint-go
    rev: v1.0.0
    hooks:
      - id: cfn-lint-go
```

See [docs/getting_started/integration.md](docs/getting_started/integration.md) for more integration examples.

### Library

```go
package main

import (
    "fmt"
    "log"

    "github.com/lex00/cfn-lint-go/pkg/lint"
)

func main() {
    linter := lint.New(lint.Options{})

    matches, err := linter.LintFile("template.yaml")
    if err != nil {
        log.Fatal(err)
    }

    for _, m := range matches {
        fmt.Printf("%s:%d: %s [%s]\n", m.Filename, m.Line, m.Message, m.Rule)
    }
}
```

### Graph Generation (Library)

```go
package main

import (
    "os"

    "github.com/lex00/cfn-lint-go/pkg/graph"
    "github.com/lex00/cfn-lint-go/pkg/template"
)

func main() {
    tmpl, _ := template.ParseFile("template.yaml")

    gen := &graph.Generator{
        IncludeParameters: true,
        Format:            graph.FormatDOT,      // or graph.FormatMermaid
        ClusterByType:     true,                 // group by AWS service
    }
    gen.Generate(tmpl, os.Stdout)
}
```

#### Graph Output Formats

- **DOT** (default): Graphviz format, render with `dot -Tpng`
- **Mermaid**: Renders natively in GitHub markdown

#### Edge Colors (DOT format)
- Black: Ref references
- Blue: GetAtt references
- Gray dashed: DependsOn dependencies

## Rule Categories

| Range | Category | Count |
|-------|----------|-------|
| E0xxx | Template errors | 7 |
| E1xxx | Functions & schema validation | 39 |
| E2xxx | Parameters | 14 |
| E3xxx | Resources & properties | 119 |
| E4xxx | Metadata | 2 |
| E5xxx | Modules | 1 |
| E6xxx | Outputs | 9 |
| E7xxx | Mappings | 3 |
| E8xxx | Conditions | 7 |
| W1xxx | Template warnings | 15 |
| W2xxx | Parameter warnings | 10 |
| W3xxx | Resource warnings | 18 |
| W4xxx | Metadata warnings | 2 |
| W6xxx | Output warnings | 1 |
| W7xxx | Mapping warnings | 1 |
| W8xxx | Condition warnings | 2 |
| Ixxx | Informational | 20 |
| **Total** | | **270** |

## NOT in Scope

- Dynamic rule loading (`--append-rules`)

## Development

```bash
# Run tests
go test -v ./...

# Run linter
golangci-lint run

# Build CLI
go build -o cfn-lint ./cmd/cfn-lint

# Run CLI
./cfn-lint --help
```

## Project Structure

```
cfn-lint-go/
├── cmd/cfn-lint/       # CLI application
├── pkg/                # Public API (importable)
│   ├── lint/           # Main linting interface
│   ├── template/       # Template parsing
│   ├── graph/          # DOT graph generation
│   ├── output/         # Output formatters (SARIF, JUnit, pretty)
│   ├── config/         # Configuration file support
│   ├── docgen/         # Documentation generator
│   ├── rules/          # Rule interface and registry
│   ├── sam/            # SAM template detection and transformation
│   └── schema/         # CloudFormation spec access (via cloudformation-schema-go)
├── internal/           # Private implementation
│   └── rules/          # Rule implementations
│       ├── errors/     # E0xxx
│       ├── functions/  # E1xxx
│       ├── parameters/ # E2xxx
│       ├── resources/  # E3xxx
│       ├── metadata/   # E4xxx
│       ├── outputs/    # E6xxx
│       ├── mappings/   # E7xxx
│       ├── conditions/ # E8xxx
│       └── warnings/   # Wxxx
├── testdata/           # Test fixtures
├── docs/               # Documentation
├── action.yml          # GitHub Action definition
└── .pre-commit-hooks.yaml  # Pre-commit hook configuration
```

## Documentation

- [Getting Started](docs/getting_started/README.md) - Installation and basic usage
- [Integration Guide](docs/getting_started/integration.md) - CI/CD, IDE, and pre-commit integration
- [Rules Reference](docs/RULES.md) - All available rules
- [Rule Creation Guide](docs/getting_started/rules.md) - How to create new rules
- [Custom Rules](docs/custom_rules.md) - Creating organization-specific rules
- [API Reference](docs/API.md) - Library API documentation

## Contributing

See [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) for how to contribute.

## License

Apache License 2.0 - see [LICENSE](LICENSE)

## Related Projects

- [aws-cloudformation/cfn-lint](https://github.com/aws-cloudformation/cfn-lint) - Original Python implementation
