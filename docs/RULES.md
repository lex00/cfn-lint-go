# Rules Reference

cfn-lint-go implements rules from [aws-cloudformation/cfn-lint](https://github.com/aws-cloudformation/cfn-lint).

## Current Status

**v0.1.0**: 2 rules implemented out of 265 planned.

## Rule Categories

| Prefix | Category | Implemented | Planned |
|--------|----------|-------------|---------|
| E0xxx | Template Errors | 1 | ~30 |
| E1xxx | Functions | 1 | ~25 |
| E2xxx | Parameters | 0 | ~20 |
| E3xxx | Resources | 0 | ~100 |
| E4xxx | Metadata | 0 | ~10 |
| E6xxx | Outputs | 0 | ~15 |
| E7xxx | Mappings | 0 | ~10 |
| E8xxx | Conditions | 0 | ~15 |
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
