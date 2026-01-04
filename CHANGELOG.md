# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **SAM Integration Phase 1**: Foundation and SAM detection
  - New `pkg/sam` package for SAM template handling
  - `IsSAMTemplate()` function to detect SAM templates by Transform or resource types
  - `IsSAMResourceType()` function to check for AWS::Serverless::* resource types
  - Support for all 9 SAM resource types: Function, Api, HttpApi, SimpleTable, LayerVersion, Application, StateMachine, Connector, GraphQLApi
  - Helper functions: `GetSAMResourceTypes()`, `HasSAMResources()`, `HasServerlessTransform()`, `GetSAMResources()`

- **SAM Integration Phase 2**: SAM Transformation Pipeline
  - `Transform()` function to convert SAM templates to CloudFormation using aws-sam-translator-go
  - `TransformBytes()` convenience function for byte-to-byte transformation
  - `SourceMap` type for mapping transformed CFN resources back to original SAM template locations
  - `TransformOptions` for configuring region, account ID, stack name, and partition
  - Automatic SAM detection and transformation in the linter
  - Error line number mapping back to original SAM template
  - New `DisableSAMTransform` option to skip automatic SAM transformation

### Changed

- Updated `cloudformation-schema-go` from v0.6.0 to v1.0.0
- Added `aws-sam-translator-go` v1.1.0 dependency
- Linter now automatically transforms SAM templates before linting

## [0.15.0] - 2026-01-03

### Added

- **Phase 16**: Multiple output formats
  - SARIF 2.1.0 format for GitHub Code Scanning integration
  - JUnit XML format for CI/CD test reporting
  - Pretty format with colorized output and code context
  - `--output` flag to write results to file
  - `--no-color` flag to disable colors in pretty format

- **Phase 17**: Configuration file support
  - Support for `.cfnlintrc`, `.cfnlintrc.yaml`, `.cfnlintrc.yml`, `.cfnlintrc.json`
  - Auto-discovery from current directory up to git root
  - `--config` flag to specify explicit config file
  - Config options: templates, ignore_templates, regions, ignore_checks, include_checks, configure_rules, format, output_file
  - CLI flags override config file settings

- **Phase 18**: Complete CLI options for Python cfn-lint parity
  - `--include-checks` flag to include specific rules even if ignored
  - `--include-experimental` flag for experimental rules
  - `--regions` flag for AWS region validation
  - Templates can now be specified from config file

- **Phase 19**: Missing parameter rules
  - E2004: Parameter NoEcho configuration for sensitive parameters
  - E2012: Parameter Type validation with SSM parameter types
  - E2014: Parameter ConstraintDescription usage validation
  - Added ConstraintDescription field to Parameter struct

- **Phase 20**: Ecosystem integrations
  - GitHub Action (`action.yml`) with SARIF upload support
  - Pre-commit hook configuration (`.pre-commit-hooks.yaml`)
  - Comprehensive integration documentation (`docs/INTEGRATIONS.md`)

### Changed

- CLI now supports templates from config file, making templates argument optional
- Updated README with all new features and integration examples
- Updated rule count: 265 rules (was 262)
- E2xxx category now has 9 rules (was 6)

### Fixed

- Template parameter parsing now includes ConstraintDescription field

## [0.14.0] - 2026-01-03

### Added

- Tests for all 34 Phase 14 warning rules:
  - W1xxx: 13 tests (W1019-W1100)
  - W2xxx: 8 tests (W2030-W2533)
  - W3xxx: 12 tests (W3011-W3693)
  - W4xxx: 1 test (W4005)

### Changed

- README.md: Updated to reflect v0.13.0 with 262 rules
- docs/RULES.md: Updated category table with accurate counts, added Phase 14 warning rules
- CHANGELOG.md: Fixed v0.13.0 entry and version links

## [0.13.0] - 2026-01-03

### Added

- Phase 14: Warning Rules Extensions (34 new rules):
  - W1019: Unused Sub parameters
  - W1020: Sub not needed without variables
  - W1028: Fn::If unreachable path
  - W1030-W1036: Function value validations
  - W1040: ToJsonString function value validation
  - W1051: Secrets Manager ARN in dynamic ref
  - W1100: YAML merge usage
  - W2030: Parameter valid value check
  - W2031: Parameter AllowedPattern check
  - W2501: Password properties configuration
  - W2506: ImageId parameter type
  - W2511: IAM policy syntax
  - W2530: SnapStart configuration
  - W2531: Lambda EOL runtime warning
  - W2533: Lambda .zip deployment properties
  - W3011: UpdateReplacePolicy/DeletionPolicy both set
  - W3034: Parameter value range check
  - W3037: IAM permission configuration
  - W3045: S3 bucket policies for access control
  - W3660: Multiple resources modifying RestApi
  - W3663: SourceAccount required
  - W3687: Ports not for certain protocols
  - W3688-W3693: DB-related warnings
  - W4005: cfn-lint configuration in Metadata
- Phase 15: Deployment Files & Modules (3 new rules):
  - E0100: Deployment file syntax validation (placeholder)
  - E0200: Parameter file syntax validation (placeholder)
  - E5001: CloudFormation Modules resource validation

### Notes

- 262 rules now implemented (99% of Python cfn-lint's 265 rules)
- E0100 and E0200 are stubs pending file parsing infrastructure
- Warning rules category expanded from 12 to 46 rules

## [0.12.0] - 2026-01-02

### Changed

- Switched graph package from hand-rolled DOT generation to [emicklei/dot](https://pkg.go.dev/github.com/emicklei/dot) library
- Added Mermaid output format via `graph.FormatMermaid` - renders natively in GitHub markdown
- Added `ClusterByType` option to group resources by AWS service
- Color-coded edges by dependency type (blue=GetAtt, gray dashed=DependsOn, black=Ref)
- Uses `NodeInitializer` and `EdgeInitializer` for consistent styling

### Dependencies

- Added `github.com/emicklei/dot` v1.10.0

## [0.11.0] - 2026-01-02

### Added

- Added `doc.go` files to all packages for pkg.go.dev documentation:
  - `pkg/lint` - Public linting API documentation
  - `pkg/template` - Template parsing documentation
  - `pkg/graph` - DOT graph generation documentation
  - `pkg/schema` - Schema validation documentation
  - `pkg/rules` - Rule interface and registry documentation

## [0.10.0] - 2026-01-02

### Added

- Phase 6: Schema-based validation complete (8 new rules):
  - E1101: Schema validation - unknown properties
  - E3014: Mutually exclusive properties
  - E3017: Required anyOf properties
  - E3018: Required oneOf properties
  - E3020: Dependent property exclusions
  - E3021: Dependent property requirements
  - E3037: Unique list items
  - E3040: Read-only properties
- Extended `pkg/schema/constraints.go` with:
  - `ResourceConstraints` struct for resource-level validation
  - MutuallyExclusive, DependentRequired, DependentExcluded constraint maps
  - OneOf/AnyOf property sets
  - ReadOnlyProperties and UniqueItems tracking

### Changed

- Total rule count: 82 -> 90

## [0.7.2] - 2026-01-02

### Dependencies

- Upgrade cloudformation-schema-go to v0.6.0

## [0.7.1] - 2026-01-02

### Dependencies

- Update cloudformation-schema-go to v0.5.0

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

[Unreleased]: https://github.com/lex00/cfn-lint-go/compare/v0.15.0...HEAD
[0.15.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.15.0
[0.14.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.14.0
[0.13.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.13.0
[0.12.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.12.0
[0.11.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.11.0
[0.10.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.10.0
[0.7.2]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.7.2
[0.7.1]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.7.1
[0.7.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.7.0
[0.6.1]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.6.1
[0.6.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.6.0
[0.5.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.5.0
[0.4.1]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.4.1
[0.4.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.4.0
[0.3.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.3.0
[0.2.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.2.0
[0.1.0]: https://github.com/lex00/cfn-lint-go/releases/tag/v0.1.0
