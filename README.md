# cfn-lint-go

CloudFormation Linter for Go - a native Go port of [aws-cloudformation/cfn-lint](https://github.com/aws-cloudformation/cfn-lint).

## Status

**v0.4.0 - Phase 2: Reference Validation**

This is a Go port of the Python cfn-lint tool. Implements core framework with 42 rules covering template and reference validation. See [docs/RESEARCH.md](docs/RESEARCH.md) for the full porting strategy.

### What's Implemented

- YAML/JSON template parsing with line number tracking
- CloudFormation intrinsic function support (!Ref, !GetAtt, !Sub, etc.)
- Rule interface and registry system
- DOT graph generation for resource dependencies
- CLI with text and JSON output formats
- CLI `graph` command for dependency visualization
- CLI `list-rules` command
- `--ignore-rules` flag
- 42 rules covering foundation, structure, and reference validation:
  - **E0xxx**: E0000 (parse), E0001 (transform), E0002 (rule processing)
  - **E1xxx**: E1001 (Ref), E1002 (size), E1005 (transform), E1010 (GetAtt), E1011 (FindInMap), E1019 (Sub), E1020 (Ref type), E1028 (Fn::If), E1040 (GetAtt format), E1041 (Ref format), E1050 (dynamic refs)
  - **E2xxx**: E2001 (param config), E2002 (param type), E2010 (param limit), E2015 (defaults)
  - **E3xxx**: E3001 (resource config), E3002 (Properties), E3003 (required props), E3004 (circular deps), E3005 (DependsOn), E3006 (type format), E3007 (unique IDs), E3010 (limit), E3015 (conditions)
  - **E4xxx**: E4002 (metadata)
  - **E6xxx**: E6001 (output structure), E6002 (Value required), E6003 (output types), E6005 (conditions), E6010 (limit)
  - **E7xxx**: E7001 (mapping config), E7010 (limit)
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
| E0xxx | Template errors | 3 rules |
| E1xxx | Functions (Ref, GetAtt) | 3 rules |
| E2xxx | Parameters | 4 rules |
| E3xxx | Resources | 6 rules |
| E4xxx | Metadata | 1 rule |
| E6xxx | Outputs | 4 rules |
| E7xxx | Mappings | 2 rules |
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
