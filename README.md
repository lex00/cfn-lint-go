# cfn-lint-go

CloudFormation Linter for Go - a native Go port of [aws-cloudformation/cfn-lint](https://github.com/aws-cloudformation/cfn-lint).

## Status

**v0.2.0 - Critical Rules**

This is a Go port of the Python cfn-lint tool. Implements core framework with one critical rule per category. See [docs/RESEARCH.md](docs/RESEARCH.md) for the full porting strategy.

### What's Implemented

- YAML/JSON template parsing with line number tracking
- CloudFormation intrinsic function support (!Ref, !GetAtt, !Sub, etc.)
- Rule interface and registry system
- DOT graph generation for resource dependencies
- CLI with text and JSON output formats
- CLI `graph` command for dependency visualization
- CLI `list-rules` command
- `--ignore-rules` flag
- 8 rules covering all category prefixes:
  - E0000: Parse errors
  - E1001: Undefined Ref
  - E2015: Parameter default constraints
  - E3003: Required resource properties
  - E4002: Metadata validation
  - E6002: Output Value required
  - E7001: Mapping configuration
  - E8002: Undefined conditions

### What's Planned

- 100+ additional rules (see [docs/RULES.md](docs/RULES.md))
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
| E0xxx | Template errors | 1 rule (E0000) |
| E1xxx | Functions (Ref, GetAtt) | 1 rule (E1001) |
| E2xxx | Parameters | 1 rule (E2015) |
| E3xxx | Resources | 1 rule (E3003) |
| E4xxx | Metadata | 1 rule (E4002) |
| E6xxx | Outputs | 1 rule (E6002) |
| E7xxx | Mappings | 1 rule (E7001) |
| E8xxx | Conditions | 1 rule (E8002) |

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
│       └── functions/  # E1xxx
├── testdata/           # Test fixtures
└── docs/               # Documentation
```

## Contributing

See [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) for how to contribute.

## License

MIT License - see [LICENSE](LICENSE)

## Related Projects

- [aws-cloudformation/cfn-lint](https://github.com/aws-cloudformation/cfn-lint) - Original Python implementation
