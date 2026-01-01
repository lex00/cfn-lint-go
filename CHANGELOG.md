# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

[0.2.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.2.0
[0.1.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.1.0
