# cfn-lint-go

CloudFormation Linter for Go - a native Go port of [aws-cloudformation/cfn-lint](https://github.com/aws-cloudformation/cfn-lint).

## Status

**v0.5.0 - Phase 4: Best Practice Rules**

This is a Go port of the Python cfn-lint tool. Implements core framework with 64 rules covering template structure, intrinsic functions, and best practices. See [docs/RESEARCH.md](docs/RESEARCH.md) for the full porting strategy.

### What's Implemented

- YAML/JSON template parsing with line number tracking
- CloudFormation intrinsic function support (!Ref, !GetAtt, !Sub, etc.)
- Rule interface and registry system
- DOT graph generation for resource dependencies
- CLI with text and JSON output formats
- CLI `graph` command for dependency visualization
- CLI `list-rules` command
- `--ignore-rules` flag
- 64 rules covering foundation, structure, intrinsics, and best practices:
  - **E0xxx**: E0000-E0003 (parse, transform, processing, config)
  - **E1xxx**: 20 rules for intrinsic functions (Ref, GetAtt, Sub, Join, Select, Split, Base64, Cidr, GetAZs, ImportValue, dynamic refs)
  - **E2xxx**: E2001-E2015 (param config, type, naming, length, limits, defaults)
  - **E3xxx**: E3001-E3036 (resource config, properties, dependencies, policies)
  - **E4xxx**: E4001-E4002 (interface metadata, structure)
  - **E6xxx**: E6001-E6102 (output structure, types, naming, exports)
  - **E7xxx**: E7001-E7010 (mapping config, naming, limits)
  - **E8xxx**: E8001-E8007 (condition functions)

### What's Planned

- 50+ additional rules (see [docs/RULES.md](docs/RULES.md))
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

    gen := &graph.Generator{IncludeParameters: true}
    gen.Generate(tmpl, os.Stdout)
}
```

## Rule Categories

| Range | Category | Status |
|-------|----------|--------|
| E0xxx | Template errors | 4 rules |
| E1xxx | Functions (Ref, GetAtt, Sub, etc.) | 20 rules |
| E2xxx | Parameters | 6 rules |
| E3xxx | Resources | 12 rules |
| E4xxx | Metadata | 2 rules |
| E6xxx | Outputs | 9 rules |
| E7xxx | Mappings | 3 rules |
| E8xxx | Conditions | 7 rules |

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
│   └── rules/          # Rule interface and registry
├── internal/           # Private implementation
│   └── rules/          # Rule implementations
│       ├── errors/     # E0xxx
│       ├── functions/  # E1xxx
│       ├── parameters/ # E2xxx
│       ├── resources/  # E3xxx
│       ├── metadata/   # E4xxx
│       ├── outputs/    # E6xxx
│       ├── mappings/   # E7xxx
│       └── conditions/ # E8xxx
├── testdata/           # Test fixtures
└── docs/               # Documentation
```

## Contributing

See [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) for how to contribute.

## License

MIT License - see [LICENSE](LICENSE)

## Related Projects

- [aws-cloudformation/cfn-lint](https://github.com/aws-cloudformation/cfn-lint) - Original Python implementation
