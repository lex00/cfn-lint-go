# cfn-lint-go

CloudFormation Linter for Go - a native Go port of [aws-cloudformation/cfn-lint](https://github.com/aws-cloudformation/cfn-lint).

## Status

**Implementation: In Progress**

This is a Go port of the Python cfn-lint tool. See [docs/RESEARCH.md](docs/RESEARCH.md) for the porting strategy and progress.

## Features

- Validate CloudFormation templates (YAML/JSON)
- 265 linting rules (porting in progress)
- DOT graph generation for resource dependencies
- Multiple output formats: text, JSON, SARIF, JUnit
- Embeddable as a Go library

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

# Generate dependency graph
cfn-lint graph template.yaml > deps.dot
dot -Tpng deps.dot -o deps.png

# List available rules
cfn-lint list-rules
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

### Graph Generation

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
| E0xxx | Template errors | 🟡 In Progress |
| E1xxx | Functions (Ref, GetAtt) | 🟡 In Progress |
| E2xxx | Parameters | ⚪ Planned |
| E3xxx | Resources | ⚪ Planned |
| E4xxx | Metadata | ⚪ Planned |
| E6xxx | Outputs | ⚪ Planned |
| E7xxx | Mappings | ⚪ Planned |
| E8xxx | Conditions | ⚪ Planned |

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
│   ├── rules/          # Rule implementations
│   │   ├── errors/     # E0xxx
│   │   ├── functions/  # E1xxx
│   │   └── ...
│   ├── decode/         # YAML/JSON parsing
│   ├── schema/         # CF resource schemas
│   └── formatters/     # Output formatters
├── schemas/            # Embedded CF schemas
└── docs/               # Documentation
```

## Contributing

See [docs/RESEARCH.md](docs/RESEARCH.md) for the porting strategy and how to help.

## License

MIT License - see [LICENSE](LICENSE)

## Related Projects

- [aws-cloudformation/cfn-lint](https://github.com/aws-cloudformation/cfn-lint) - Original Python implementation
- [wetwire](https://github.com/lex00/wetwire) - Infrastructure as code framework
