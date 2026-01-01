# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Phase 2: Reference Validation rules (12 new rules):
  - E1010: GetAtt to undefined resource
  - E1011: FindInMap references undefined mapping
  - E1019: Sub function validation (undefined variables)
  - E1020: Ref value must be a string
  - E1028: Fn::If structure error (undefined condition)
  - E1040: GetAtt format error
  - E1041: Ref format error
  - E1050: Dynamic reference syntax error
  - E3004: Circular resource dependency detected
  - E3005: DependsOn references undefined resource
  - E3015: Resource condition references undefined condition
  - E6005: Output condition references undefined condition
- Comprehensive tests for all Phase 2 rules

### Changed

- Total rule count: 30 -> 42

## [0.3.0] - 2026-01-01

### Added

- Phase 1: Foundation & Structure rules (22 new rules):
  - E0001: Template transformation error
  - E0002: Rule processing error
  - E1002: Template size limit exceeded (51KB direct, 460KB S3)
  - E1005: Transform configuration error
  - E2001: Parameter configuration error (valid properties, required Type)
  - E2002: Invalid parameter type
  - E2010: Parameter limit exceeded (200 max)
  - E3001: Resource configuration error (valid properties, required Type)
  - E3002: Resource Properties structure error (must be object)
  - E3006: Invalid resource type format (AWS::*, Custom::*, Alexa::*)
  - E3007: Duplicate resource logical ID
  - E3010: Resource limit exceeded (500 max)
  - E6001: Output property structure error (valid properties only)
  - E6003: Output property type error (Export requires Name)
  - E6010: Output limit exceeded (200 max)
  - E7010: Mapping limit exceeded (200 max)
  - E8001: Condition configuration error (valid condition functions)
  - E8003: Fn::Equals structure error (must have exactly 2 elements)
  - E8004: Fn::And structure error (2-10 conditions)
  - E8005: Fn::Not structure error (exactly 1 condition)
  - E8006: Fn::Or structure error (2-10 conditions)
  - E8007: Condition intrinsic function error
- Comprehensive tests for all Phase 1 rules

### Changed

- Total rule count: 8 -> 30

## [0.2.0] - 2026-01-01

### Added

- One critical rule per category prefix (6 new rules):
  - E2015: Parameter default value must satisfy constraints (AllowedValues, AllowedPattern, MinValue/MaxValue, MinLength/MaxLength)
  - E3003: Required resource properties are present (covers common AWS resources)
  - E4002: Metadata section has valid structure (no null values)
  - E6002: Output has required Value property
  - E7001: Mapping configuration is valid (proper structure and keys)
  - E8002: Referenced conditions are defined
- Enhanced template parser:
  - Mappings section parsing with nested key structure
  - Conditions section parsing with expression support
  - Metadata section parsing
  - Parameter constraint fields (AllowedValues, AllowedPattern, MinValue, MaxValue, MinLength, MaxLength)
  - Output Value and Export field parsing
- Complete rules matrix in docs/RESEARCH.md with criticality levels and ordering dependencies
- Comprehensive tests for all new rules

### Changed

- Total rule count: 2 -> 8
- All category prefixes now have at least one critical rule implemented

## [0.1.0] - 2026-01-01

### Added

- Initial release of cfn-lint-go
- Core framework for CloudFormation template linting
- YAML/JSON template parsing with line number tracking
- CloudFormation intrinsic function tag support (!Ref, !GetAtt, !Sub, etc.)
- DOT graph generation for resource dependencies
- Rule interface and registry system
- Two initial rules:
  - E0000: Template parse error detection
  - E1001: Ref to undefined resource or parameter
- CLI with multiple output formats (text, JSON, SARIF, JUnit)
- Comprehensive test suite
- Documentation (API.md, RULES.md, CONTRIBUTING.md)
- GitHub Actions CI workflow

### Notes

- This is a Go port of [aws-cloudformation/cfn-lint](https://github.com/aws-cloudformation/cfn-lint)
- SAM transform support is not included (users should run `sam build` first)
- See [docs/RESEARCH.md](docs/RESEARCH.md) for the full porting strategy

[Unreleased]: https://github.com/lex00/cfn-lint-go/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.3.0
[0.2.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.2.0
[0.1.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.1.0
