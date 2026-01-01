# Contributing to cfn-lint-go

Thank you for your interest in contributing to cfn-lint-go!

## Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/cfn-lint-go.git
   cd cfn-lint-go
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Run tests:
   ```bash
   go test ./...
   ```

## Development Workflow

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./pkg/template/...
```

### Linting

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

### Building

```bash
# Build CLI
go build -o cfn-lint ./cmd/cfn-lint

# Run CLI
./cfn-lint --help
```

## Adding a New Rule

1. **Find the Python implementation** at:
   ```
   https://github.com/aws-cloudformation/cfn-lint/tree/main/src/cfnlint/rules
   ```

2. **Create the Go file** in the appropriate category:
   ```
   internal/rules/{category}/{rule_id}.go
   ```

3. **Implement the Rule interface**:
   ```go
   package functions

   import (
       "github.com/lex00/cfn-lint-go/pkg/rules"
       "github.com/lex00/cfn-lint-go/pkg/template"
   )

   func init() {
       rules.Register(&E1001{})
   }

   type E1001 struct{}

   func (r *E1001) ID() string          { return "E1001" }
   func (r *E1001) ShortDesc() string   { return "Brief description" }
   func (r *E1001) Description() string { return "Detailed description" }
   func (r *E1001) Source() string      { return "https://..." }
   func (r *E1001) Tags() []string      { return []string{"functions", "ref"} }

   func (r *E1001) Match(tmpl *template.Template) []rules.Match {
       // Implementation
       return nil
   }
   ```

4. **Add tests** in `{rule_id}_test.go`

5. **Add test fixtures** if needed in `testdata/`

## Rule Categories

| Directory | Rule Range | Description |
|-----------|------------|-------------|
| `errors/` | E0xxx | Template parse and structure errors |
| `functions/` | E1xxx | Intrinsic function validation |
| `parameters/` | E2xxx | Parameter validation |
| `resources/` | E3xxx | Resource property validation |
| `metadata/` | E4xxx | Metadata validation |
| `outputs/` | E6xxx | Output validation |
| `mappings/` | E7xxx | Mapping validation |
| `conditions/` | E8xxx | Condition validation |

## Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Match error messages to Python cfn-lint when possible
- Add doc comments to exported functions and types

## Pull Request Guidelines

1. Create a feature branch from `main`
2. Add tests for new functionality
3. Ensure all tests pass: `go test ./...`
4. Ensure linter passes: `golangci-lint run`
5. Update documentation if needed
6. Submit PR with clear description

## Questions?

Open an issue for any questions or concerns.
