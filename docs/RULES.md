# Rules Reference

cfn-lint-go implements rules from [aws-cloudformation/cfn-lint](https://github.com/aws-cloudformation/cfn-lint).

## Current Status

**v0.4.0**: 42 rules implemented (Phase 1 + Phase 2: Reference Validation complete).

## Rule Categories

| Prefix | Category | Implemented | Planned |
|--------|----------|-------------|---------|
| E0xxx | Template Errors | 3 | ~6 |
| E1xxx | Functions | 11 | ~30 |
| E2xxx | Parameters | 4 | ~7 |
| E3xxx | Resources | 9 | ~40+ |
| E4xxx | Metadata | 1 | ~2 |
| E6xxx | Outputs | 5 | ~9 |
| E7xxx | Mappings | 2 | ~3 |
| E8xxx | Conditions | 7 | ~7 |
| W* | Warnings | 0 | ~40 |
| I* | Informational | 0 | ~20 |

## Implemented Rules

### E0xxx - Template Errors

| Rule | Description | Status |
|------|-------------|--------|
| E0000 | Template parse error | ✅ Implemented |
| E0001 | Template transformation error | ✅ Implemented |
| E0002 | Rule processing error | ✅ Implemented |

### E1xxx - Functions

| Rule | Description | Status |
|------|-------------|--------|
| E1001 | Ref to undefined resource or parameter | ✅ Implemented |
| E1002 | Template size limit exceeded | ✅ Implemented |
| E1005 | Transform configuration error | ✅ Implemented |
| E1010 | GetAtt to undefined resource | ✅ Implemented |
| E1011 | FindInMap references undefined mapping | ✅ Implemented |
| E1019 | Sub function validation | ✅ Implemented |
| E1020 | Ref value must be a string | ✅ Implemented |
| E1028 | Fn::If structure error | ✅ Implemented |
| E1040 | GetAtt format error | ✅ Implemented |
| E1041 | Ref format error | ✅ Implemented |
| E1050 | Dynamic reference syntax error | ✅ Implemented |

### E2xxx - Parameters

| Rule | Description | Status |
|------|-------------|--------|
| E2001 | Parameter configuration error | ✅ Implemented |
| E2002 | Invalid parameter type | ✅ Implemented |
| E2010 | Parameter limit exceeded (200) | ✅ Implemented |
| E2015 | Default value within constraints | ✅ Implemented |

### E3xxx - Resources

| Rule | Description | Status |
|------|-------------|--------|
| E3001 | Resource configuration error | ✅ Implemented |
| E3002 | Resource Properties structure error | ✅ Implemented |
| E3003 | Required properties present | 🚧 Partial (common resources) |
| E3004 | Circular resource dependency detected | ✅ Implemented |
| E3005 | DependsOn references undefined resource | ✅ Implemented |
| E3006 | Invalid resource type format | ✅ Implemented |
| E3007 | Duplicate resource logical ID | ✅ Implemented |
| E3010 | Resource limit exceeded (500) | ✅ Implemented |
| E3015 | Resource condition references undefined condition | ✅ Implemented |

### E4xxx - Metadata

| Rule | Description | Status |
|------|-------------|--------|
| E4002 | Metadata section structure | ✅ Implemented |

### E6xxx - Outputs

| Rule | Description | Status |
|------|-------------|--------|
| E6001 | Output property structure error | ✅ Implemented |
| E6002 | Output has required Value property | ✅ Implemented |
| E6003 | Output property type error | ✅ Implemented |
| E6005 | Output condition references undefined condition | ✅ Implemented |
| E6010 | Output limit exceeded (200) | ✅ Implemented |

### E7xxx - Mappings

| Rule | Description | Status |
|------|-------------|--------|
| E7001 | Mapping configuration valid | ✅ Implemented |
| E7010 | Mapping limit exceeded (200) | ✅ Implemented |

### E8xxx - Conditions

| Rule | Description | Status |
|------|-------------|--------|
| E8001 | Condition configuration error | ✅ Implemented |
| E8002 | Referenced conditions are defined | ✅ Implemented |
| E8003 | Fn::Equals structure error | ✅ Implemented |
| E8004 | Fn::And structure error | ✅ Implemented |
| E8005 | Fn::Not structure error | ✅ Implemented |
| E8006 | Fn::Or structure error | ✅ Implemented |
| E8007 | Condition intrinsic function error | ✅ Implemented |

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
