# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

[0.1.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.1.0
