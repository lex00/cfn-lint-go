# cfn-lint-go

CloudFormation Linter for Go - a native Go port of [aws-cloudformation/cfn-lint](https://github.com/aws-cloudformation/cfn-lint).

## Status

**v0.13.0 - 262 rules implemented**

This is a Go port of the Python cfn-lint tool. Implements core framework with 262 rules covering template structure, intrinsic functions, schema-based validation, best practices, warnings, and informational rules. Uses `cloudformation-schema-go` for CloudFormation resource specification and enum validation. See [docs/RESEARCH.md](docs/RESEARCH.md) for the full porting strategy.

### What's Implemented

- YAML/JSON template parsing with line number tracking
- CloudFormation intrinsic function support (!Ref, !GetAtt, !Sub, etc.)
- Rule interface and registry system
- DOT graph generation for resource dependencies
- CLI with text and JSON output formats
- CLI `graph` command for dependency visualization
- CLI `list-rules` command
- `--ignore-rules` flag
- 262 rules across all categories:
  - **E0xxx**: 6 rules (parse, transform, processing, config, deployment/parameter files)
  - **E1xxx**: 27 rules for intrinsic functions and schema validation
  - **E2xxx**: 6 rules (param config, type, naming, length, limits, defaults)
  - **E3xxx**: 124 rules (resource config, properties, type validation, enum validation, dependencies, policies, constraints)
  - **E4xxx**: 2 rules (interface metadata, structure)
  - **E5xxx**: 1 rule (CloudFormation Modules validation)
  - **E6xxx**: 9 rules (output structure, types, naming, exports)
  - **E7xxx**: 3 rules (mapping config, naming, limits)
  - **E8xxx**: 7 rules (condition functions)
  - **Wxxx**: 46 warning rules (security, best practices, deprecations)
  - **Ixxx**: 19 informational rules

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

| Range | Category | Count |
|-------|----------|-------|
| E0xxx | Template errors | 6 |
| E1xxx | Functions & schema validation | 27 |
| E2xxx | Parameters | 6 |
| E3xxx | Resources & properties | 124 |
| E4xxx | Metadata | 2 |
| E5xxx | Modules | 1 |
| E6xxx | Outputs | 9 |
| E7xxx | Mappings | 3 |
| E8xxx | Conditions | 7 |
| W1xxx | Template warnings | 15 |
| W2xxx | Parameter warnings | 10 |
| W3xxx | Resource warnings | 15 |
| W4xxx | Metadata warnings | 2 |
| W6xxx | Output warnings | 1 |
| W7xxx | Mapping warnings | 1 |
| W8xxx | Condition warnings | 2 |
| Ixxx | Informational | 19 |
| **Total** | | **262** |

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
