# Rules Reference

cfn-lint-go implements rules from [aws-cloudformation/cfn-lint](https://github.com/aws-cloudformation/cfn-lint).

## Rule Categories

| Prefix | Category | Count | Status |
|--------|----------|-------|--------|
| E0xxx | Template Errors | ~30 | 🟡 In Progress |
| E1xxx | Functions | ~25 | 🟡 In Progress |
| E2xxx | Parameters | ~20 | ⚪ Planned |
| E3xxx | Resources | ~100 | ⚪ Planned |
| E4xxx | Metadata | ~10 | ⚪ Planned |
| E6xxx | Outputs | ~15 | ⚪ Planned |
| E7xxx | Mappings | ~10 | ⚪ Planned |
| E8xxx | Conditions | ~15 | ⚪ Planned |
| W* | Warnings | ~40 | ⚪ Planned |
| I* | Informational | ~20 | ⚪ Planned |

## Implemented Rules

### E0xxx - Template Errors

| Rule | Description | Status |
|------|-------------|--------|
| E0000 | Template parse error | ✅ |

### E1xxx - Functions

| Rule | Description | Status |
|------|-------------|--------|
| E1001 | Ref to undefined resource or parameter | ✅ |

## Rule Severity Levels

- **E (Error)**: Must be fixed for valid CloudFormation
- **W (Warning)**: Best practice violations
- **I (Informational)**: Suggestions and tips

## Ignoring Rules

### CLI

```bash
cfn-lint template.yaml --ignore-rules E1001,W3002
```

### Template Metadata

```yaml
Metadata:
  cfn-lint:
    config:
      ignore_checks:
        - E1001
        - W3002
```

### Inline Comments

```yaml
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Metadata:
      cfn-lint:
        config:
          ignore_checks:
            - E3002
```

## Adding Custom Rules

See [CONTRIBUTING.md](CONTRIBUTING.md) for how to add new rules.

## Rule Parity with Python cfn-lint

This project aims for full parity with Python cfn-lint rules, excluding:

- SAM transform rules (requires aws-sam-translator)
- Dynamic rule loading (`--append-rules`)

See [RESEARCH.md](RESEARCH.md) for the porting strategy.
