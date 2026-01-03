# cfn-lint-go

CloudFormation Linter for Go - a native Go port of [aws-cloudformation/cfn-lint](https://github.com/aws-cloudformation/cfn-lint).

## Status

**v0.16.0 - Phase 9: Parameter Extensions**

This is a Go port of the Python cfn-lint tool. Implements core framework with 131 rules covering template structure, intrinsic functions, schema-based validation, best practices, and warnings. Uses `cloudformation-schema-go` for CloudFormation resource specification and enum validation. See [docs/RESEARCH.md](docs/RESEARCH.md) for the full porting strategy.

### What's Implemented

- YAML/JSON template parsing with line number tracking
- CloudFormation intrinsic function support (!Ref, !GetAtt, !Sub, etc.)
- Rule interface and registry system
- DOT graph generation for resource dependencies
- CLI with text and JSON output formats
- CLI `graph` command for dependency visualization
- CLI `list-rules` command
- `--ignore-rules` flag
- 131 rules covering foundation, structure, intrinsics, schema validation, best practices, and warnings:
  - **E0xxx**: E0000-E0003, E0100, E0200 (parse, transform, processing, config, deployment/parameter files)
  - **E1xxx**: 34 rules for intrinsic functions and schema validation
  - **E2xxx**: E2001-E2015, E2529-E2533, E2900 (param config, type, naming, length, limits, defaults, Lambda runtime, SubscriptionFilters, deployment files)
  - **E3xxx**: E3001-E3040 (resource config, properties, type validation, enum validation, dependencies, policies, constraints)
  - **E4xxx**: E4001-E4002 (interface metadata, structure)
  - **E5xxx**: E5001 (CloudFormation Modules validation)
  - **E6xxx**: E6001-E6102 (output structure, types, naming, exports)
  - **E7xxx**: E7001-E7010 (mapping config, naming, limits)
  - **E8xxx**: E8001-E8007 (condition functions)
  - **W1xxx-W8xxx**: 12 warning rules (unused resources, security, best practices)
  - **I1xxx-I7xxx**: 20 informational rules

### What's Planned

- SARIF, JUnit output formats
- Rule ignoring via template metadata

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

# Lint with JSON output
cfn-lint template.yaml --format json

# Ignore specific rules
cfn-lint template.yaml --ignore-rules E1001,W3002

# Generate dependency graph
cfn-lint graph template.yaml > deps.dot
dot -Tpng deps.dot -o deps.png

# Include parameters in graph
cfn-lint graph template.yaml --include-parameters

# List available rules
cfn-lint list-rules

# List rules as JSON
cfn-lint list-rules --format json

# Show help
cfn-lint --help
```

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

| Range | Category | Status |
|-------|----------|--------|
| E0xxx | Template errors | 4 rules |
| E1xxx | Functions & schema validation | 34 rules |
| E2xxx | Parameters | 11 rules |
| E3xxx | Resources & properties | 25 rules |
| E4xxx | Metadata | 2 rules |
| E6xxx | Outputs | 9 rules |
| E7xxx | Mappings | 3 rules |
| E8xxx | Conditions | 7 rules |
| W1xxx | Template warnings | 2 rules |
| W2xxx | Parameter warnings | 2 rules |
| W3xxx | Resource warnings | 3 rules |
| W4xxx | Metadata warnings | 1 rule |
| W6xxx | Output warnings | 1 rule |
| W7xxx | Mapping warnings | 1 rule |
| W8xxx | Condition warnings | 2 rules |

## NOT in Scope

- SAM transform support (use `sam build` first)
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
│   ├── rules/          # Rule interface and registry
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
└── docs/               # Documentation
```

## Contributing

See [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) for how to contribute.

## License

MIT License - see [LICENSE](LICENSE)

## Related Projects

- [aws-cloudformation/cfn-lint](https://github.com/aws-cloudformation/cfn-lint) - Original Python implementation
