# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.7.0] - 2026-01-01

### Added

- `pkg/schema` package integrating `cloudformation-schema-go` for CloudFormation resource specification access
- Schema-based required property validation (E3003) using full CloudFormation resource specification

### Changed

- E3003 now uses live CloudFormation resource specification instead of hardcoded property map

### Dependencies

- Added `github.com/lex00/cloudformation-schema-go` v0.4.0

## [0.6.1] - 2026-01-01

### Changed

- JSON output format now matches Python cfn-lint for drop-in compatibility
  - Nested `Rule` object with `Id`, `Description`, `ShortDescription`, `Source`
  - Nested `Location` with `Start`/`End` positions (`LineNumber`, `ColumnNumber`)
  - `Path` array for JSON path to issue location
  - PascalCase level values (`"Error"`, `"Warning"`, `"Informational"`)

## [0.6.0] - 2026-01-01

### Added

- Phase 5: Warning rules (12 new rules):
  - W1001: Ref/GetAtt to conditional resource
  - W1011: Use dynamic references for secrets
  - W2001: Unused parameter detection
  - W2010: NoEcho parameter may be exposed in outputs
  - W3002: Package-required property with local path
  - W3005: Redundant DependsOn (implicit dependency via Ref/GetAtt)
  - W3010: Hardcoded availability zone
  - W4001: Interface references undefined parameter
  - W6001: ImportValue in Output (circular dependency risk)
  - W7001: Unused mapping detection
  - W8001: Unused condition detection
  - W8003: Fn::Equals with static result (always true/false)
- New `warnings` package for all warning rules
- Comprehensive tests for all Phase 5 rules

### Changed

- Total rule count: 64 -> 76

## [0.5.0] - 2026-01-01

### Added

- Phase 4: Important Best Practice rules (20 new rules):
  - E0003: Configuration error (AWSTemplateFormatVersion, Resources required)
  - E1004: Description must be a string (max 1024 chars)
  - E1015: Fn::GetAZs function validation
  - E1016: Fn::ImportValue function validation
  - E1017: Fn::Select function validation
  - E1018: Fn::Split function validation
  - E1021: Fn::Base64 function validation
  - E1022: Fn::Join function validation
  - E1024: Fn::Cidr function validation
  - E1027: Dynamic reference in invalid location
  - E1029: Sub required for variable substitution
  - E2003: Parameter naming convention (alphanumeric only)
  - E2011: Parameter name length (max 255)
  - E3011: Invalid property name validation
  - E3035: Invalid DeletionPolicy validation
  - E3036: Invalid UpdateReplacePolicy validation
  - E4001: AWS::CloudFormation::Interface metadata validation
  - E6004: Output naming convention (alphanumeric only)
  - E6011: Output name length (max 255)
  - E7002: Mapping name length (max 255)
- Comprehensive tests for all Phase 4 rules

### Changed

- Total rule count: 44 -> 64

## [0.4.1] - 2026-01-01

### Added

- Phase 3: Value Validation rules (2 new rules):
  - E6101: Output Value must be a string
  - E6102: Export Name must be a string
- Comprehensive tests for all Phase 3 rules

### Changed

- Total rule count: 42 -> 44

### Notes

- E2900 (Deployment file parameters) deferred - requires parameter file infrastructure

## [0.4.0] - 2026-01-01

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

[Unreleased]: https://github.com/lex00/cfn-lint-go/compare/v0.6.0...HEAD
[0.6.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.6.0
[0.5.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.5.0
[0.4.1]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.4.1
[0.4.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.4.0
[0.3.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.3.0
[0.2.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.2.0
[0.1.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.1.0
