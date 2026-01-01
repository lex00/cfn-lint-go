# Rules Reference

cfn-lint-go implements rules from [aws-cloudformation/cfn-lint](https://github.com/aws-cloudformation/cfn-lint).

## Current Status

**v0.2.0**: 8 rules implemented (one critical rule per category prefix).

## Rule Categories

| Prefix | Category | Implemented | Planned |
|--------|----------|-------------|---------|
| E0xxx | Template Errors | 1 | ~6 |
| E1xxx | Functions | 1 | ~30 |
| E2xxx | Parameters | 1 | ~7 |
| E3xxx | Resources | 1 | ~40+ |
| E4xxx | Metadata | 1 | ~2 |
| E6xxx | Outputs | 1 | ~9 |
| E7xxx | Mappings | 1 | ~3 |
| E8xxx | Conditions | 1 | ~7 |
| W* | Warnings | 0 | ~40 |
| I* | Informational | 0 | ~20 |

## Implemented Rules

### E0xxx - Template Errors

| Rule | Description | Status |
|------|-------------|--------|
| E0000 | Template parse error | ✅ Implemented |

### E1xxx - Functions

| Rule | Description | Status |
|------|-------------|--------|
| E1001 | Ref to undefined resource or parameter | ✅ Implemented |

### E2xxx - Parameters

| Rule | Description | Status |
|------|-------------|--------|
| E2015 | Default value within constraints | ✅ Implemented |

### E3xxx - Resources

| Rule | Description | Status |
|------|-------------|--------|
| E3003 | Required properties present | 🚧 Partial (common resources) |

### E4xxx - Metadata

| Rule | Description | Status |
|------|-------------|--------|
| E4002 | Metadata section structure | ✅ Implemented |

### E6xxx - Outputs

| Rule | Description | Status |
|------|-------------|--------|
| E6002 | Output has required Value property | ✅ Implemented |

### E7xxx - Mappings

| Rule | Description | Status |
|------|-------------|--------|
| E7001 | Mapping configuration valid | ✅ Implemented |

### E8xxx - Conditions

| Rule | Description | Status |
|------|-------------|--------|
| E8002 | Referenced conditions are defined | ✅ Implemented |

## Rule Severity Levels

- **E (Error)**: Must be fixed for valid CloudFormation
- **W (Warning)**: Best practice violations
- **I (Informational)**: Suggestions and tips

## Ignoring Rules

### CLI

```bash
cfn-lint template.yaml --ignore-rules E1001,W3002
```

### Library API

```go
linter := lint.New(lint.Options{
    IgnoreRules: []string{"E1001", "W3002"},
})
```

### Template Metadata (Planned)

```yaml
Metadata:
  cfn-lint:
    config:
      ignore_checks:
        - E1001
        - W3002
```

## Adding Custom Rules

See [CONTRIBUTING.md](CONTRIBUTING.md) for how to add new rules.

## Rule Parity with Python cfn-lint

This project aims for full parity with Python cfn-lint rules, excluding:

- SAM transform rules (requires aws-sam-translator)
- Dynamic rule loading (`--append-rules`)

See [RESEARCH.md](RESEARCH.md) for the porting strategy.
