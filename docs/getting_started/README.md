# Getting Started

Welcome to cfn-lint-go! This guide will help you get started with validating your CloudFormation templates.

## Quick Links

- [Installation](#installation)
- [Basic Usage](#basic-usage)
- [SAM Templates](#sam-templates)
- [Configuration](#configuration)
- [Integration Guide](integration.md) - CI/CD, IDE, and pre-commit integration
- [Rule Creation Guide](rules.md) - How to create new rules
- [Custom Rules](../custom_rules.md) - Loading custom rules
- [SAM Template Support](../SAM.md) - Full SAM documentation

## Installation

### Go Install

```bash
go install github.com/lex00/cfn-lint-go/cmd/cfn-lint@latest
```

### From Source

```bash
git clone https://github.com/lex00/cfn-lint-go.git
cd cfn-lint-go
go build -o cfn-lint ./cmd/cfn-lint
```

### Verify Installation

```bash
cfn-lint --version
```

## Basic Usage

### Lint a Template

```bash
cfn-lint template.yaml
```

### Multiple Templates

```bash
cfn-lint template1.yaml template2.yaml
cfn-lint templates/*.yaml
```

### Output Formats

cfn-lint-go supports multiple output formats:

```bash
# Default text format
cfn-lint template.yaml

# JSON format
cfn-lint template.yaml --format json

# SARIF format (for GitHub Code Scanning)
cfn-lint template.yaml --format sarif --output results.sarif

# JUnit XML format (for CI/CD)
cfn-lint template.yaml --format junit --output results.xml

# Pretty format (colorized with context)
cfn-lint template.yaml --format pretty
```

### Ignore Rules

```bash
# Ignore specific rules
cfn-lint template.yaml --ignore-rules E1001,W3002

# Include rules even if ignored elsewhere
cfn-lint template.yaml --ignore-rules E1001 --include-checks E1001
```

## SAM Templates

cfn-lint-go natively supports AWS SAM templates through [aws-sam-translator-go](https://github.com/lex00/aws-sam-translator-go).

### Automatic Detection

SAM templates are automatically detected and transformed to CloudFormation before linting:

```bash
# Just lint like any template - SAM is auto-detected
cfn-lint sam-template.yaml
```

### SAM-Specific Options

```bash
# Skip SAM transformation (lint as-is)
cfn-lint sam-template.yaml --no-sam-transform

# Output the transformed CloudFormation (for debugging)
cfn-lint sam-template.yaml --show-transformed
```

### SAM Configuration

Configure SAM behavior in `.cfnlintrc.yaml`:

```yaml
sam:
  auto_transform: true
  transform_options:
    region: us-east-1
    account_id: "123456789012"
    stack_name: my-app
```

For complete SAM documentation, see [SAM Template Support](../SAM.md).

## Configuration

### Configuration File

Create a `.cfnlintrc.yaml` file in your project root:

```yaml
# Templates to lint (supports globs)
templates:
  - templates/**/*.yaml
  - infrastructure/*.yml

# Templates to ignore
ignore_templates:
  - test/**
  - .aws-sam/**

# AWS regions to validate against
regions:
  - us-east-1
  - us-west-2

# Rules to ignore
ignore_checks:
  - E1001
  - W3002

# Rules to include (even if ignored)
include_checks:
  - I1001

# Include experimental rules
include_experimental: false

# Output format
format: pretty
```

### Configuration File Discovery

cfn-lint-go automatically searches for configuration files:

1. `.cfnlintrc.yaml` in current directory
2. `.cfnlintrc.yml` in current directory
3. `.cfnlintrc.json` in current directory
4. `.cfnlintrc` in current directory
5. Same search up to git root

Override with `--config`:

```bash
cfn-lint template.yaml --config custom-config.yaml
```

## Commands

### lint (default)

Validate CloudFormation templates:

```bash
cfn-lint template.yaml
```

### list-rules

List all available rules:

```bash
cfn-lint list-rules
cfn-lint list-rules --format json
```

### graph

Generate dependency graph:

```bash
cfn-lint graph template.yaml > deps.dot
dot -Tpng deps.dot -o deps.png
```

### update-documentation

Update RULES.md from registered rules:

```bash
cfn-lint update-documentation
cfn-lint update-documentation --output docs/RULES.md
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success, no issues found |
| 1 | Error running the tool |
| 2 | Issues found in templates |

## Next Steps

- Read the [Integration Guide](integration.md) to set up CI/CD
- Check out [Rule Creation](rules.md) to contribute new rules
- See the full [Rules Reference](../RULES.md) for all available rules
