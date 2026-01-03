# Integrations

This document describes how to integrate cfn-lint-go into your development workflow.

## Table of Contents

- [GitHub Actions](#github-actions)
- [Pre-commit Hooks](#pre-commit-hooks)
- [CI/CD Integration](#cicd-integration)
- [IDE Integration](#ide-integration)

## GitHub Actions

Use the cfn-lint-go GitHub Action to automatically validate CloudFormation templates in your pull requests.

### Basic Usage

```yaml
name: Validate CloudFormation Templates

on:
  pull_request:
    paths:
      - '**.yaml'
      - '**.yml'
      - '**.json'

jobs:
  cfn-lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - uses: lex00/cfn-lint-go@main
        with:
          templates: 'templates/*.yaml'
```

### SARIF Output with Code Scanning

Use SARIF format to integrate with GitHub Code Scanning:

```yaml
name: CloudFormation Lint

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

permissions:
  security-events: write
  contents: read

jobs:
  cfn-lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - uses: lex00/cfn-lint-go@main
        with:
          templates: '**/*.yaml'
          format: sarif
          ignore-rules: E1001,W3002
```

### Configuration File

Use a configuration file for consistent settings:

```yaml
- uses: lex00/cfn-lint-go@main
  with:
    templates: 'templates/*.yaml'
    config: .cfnlintrc.yaml
```

## Pre-commit Hooks

Integrate cfn-lint-go into your pre-commit workflow to catch issues before committing.

### Installation

1. Install [pre-commit](https://pre-commit.com/):

```bash
pip install pre-commit
```

2. Add to `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: https://github.com/lex00/cfn-lint-go
    rev: v0.15.0
    hooks:
      - id: cfn-lint-go
```

3. Install the hook:

```bash
pre-commit install
```

### Configuration

Customize the hook behavior:

```yaml
repos:
  - repo: https://github.com/lex00/cfn-lint-go
    rev: v0.15.0
    hooks:
      - id: cfn-lint-go
        args:
          - --format=pretty
          - --ignore-rules=E1001,W3002
        files: ^templates/.*\.(yaml|yml)$
```

## CI/CD Integration

### JUnit XML for CI Systems

Many CI systems (Jenkins, GitLab CI, CircleCI) support JUnit XML format for test reporting:

```bash
cfn-lint template.yaml --format junit --output results.xml
```

#### Jenkins

```groovy
pipeline {
    agent any
    stages {
        stage('Validate CloudFormation') {
            steps {
                sh 'cfn-lint templates/*.yaml --format junit --output cfn-lint-results.xml'
            }
            post {
                always {
                    junit 'cfn-lint-results.xml'
                }
            }
        }
    }
}
```

#### GitLab CI

```yaml
cfn-lint:
  stage: test
  script:
    - go install github.com/lex00/cfn-lint-go/cmd/cfn-lint@latest
    - cfn-lint templates/*.yaml --format junit --output cfn-lint-results.xml
  artifacts:
    reports:
      junit: cfn-lint-results.xml
```

#### CircleCI

```yaml
version: 2.1
jobs:
  cfn-lint:
    docker:
      - image: cimg/go:1.21
    steps:
      - checkout
      - run:
          name: Install cfn-lint-go
          command: go install github.com/lex00/cfn-lint-go/cmd/cfn-lint@latest
      - run:
          name: Run cfn-lint
          command: cfn-lint templates/*.yaml --format junit --output cfn-lint-results.xml
      - store_test_results:
          path: cfn-lint-results.xml
```

### AWS CodeBuild

```yaml
version: 0.2

phases:
  install:
    commands:
      - go install github.com/lex00/cfn-lint-go/cmd/cfn-lint@latest
  build:
    commands:
      - cfn-lint templates/*.yaml --format pretty

reports:
  cfn-lint:
    files:
      - cfn-lint-results.xml
    file-format: JunitXml
```

## IDE Integration

### Visual Studio Code

While a dedicated VS Code extension is planned, you can use the Tasks feature to run cfn-lint-go:

1. Create `.vscode/tasks.json`:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Validate CloudFormation",
      "type": "shell",
      "command": "cfn-lint",
      "args": ["${file}", "--format", "pretty"],
      "problemMatcher": [],
      "presentation": {
        "reveal": "always",
        "panel": "new"
      }
    }
  ]
}
```

2. Run with `Ctrl+Shift+P` > "Tasks: Run Task" > "Validate CloudFormation"

### Vim/Neovim

Use ALE (Asynchronous Lint Engine):

```vim
let g:ale_linters = {
\   'yaml': ['cfn-lint-go'],
\}

let g:ale_yaml_cfn_lint_go_executable = 'cfn-lint'
let g:ale_yaml_cfn_lint_go_options = '--format json'
```

## Configuration File

Create a `.cfnlintrc.yaml` file in your project root for consistent linting settings:

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
  - E1001  # Free-form properties are allowed
  - W3002  # Hardcoded resources

# Rules to include (even if ignored elsewhere)
include_checks:
  - I1001

# Include experimental rules
include_experimental: false

# Output format
format: pretty

# Output file (optional)
# output_file: lint-results.json
```

The config file is automatically discovered by searching from the current directory up to the git root.

Override with `--config`:

```bash
cfn-lint template.yaml --config custom-config.yaml
```

## Output Formats

cfn-lint-go supports multiple output formats for different use cases:

### Text (Default)

Human-readable output:

```bash
cfn-lint template.yaml
```

### JSON

Machine-readable output:

```bash
cfn-lint template.yaml --format json
```

### SARIF

For GitHub Code Scanning and other SARIF-compatible tools:

```bash
cfn-lint template.yaml --format sarif --output results.sarif
```

### JUnit XML

For CI/CD test reporting:

```bash
cfn-lint template.yaml --format junit --output results.xml
```

### Pretty

Colorized output with code context:

```bash
cfn-lint template.yaml --format pretty
```

Disable colors:

```bash
cfn-lint template.yaml --format pretty --no-color
```

## Best Practices

1. **Use Configuration Files**: Store project-specific settings in `.cfnlintrc.yaml` for consistency across team.

2. **Integrate Early**: Run linting in pre-commit hooks to catch issues before they reach CI.

3. **CI/CD Validation**: Always validate templates in your CI/CD pipeline before deployment.

4. **SARIF for PRs**: Use SARIF format in GitHub Actions to annotate PRs with inline comments.

5. **Fail on Warnings**: In production pipelines, consider failing on warnings:
   ```yaml
   - run: cfn-lint template.yaml
     # Exit code 2 means issues found
   ```

6. **Template Metadata**: Use template metadata to override rules per-template:
   ```yaml
   Metadata:
     cfn-lint:
       config:
         ignore_checks:
           - E3012
   ```

## Troubleshooting

### Pre-commit Hook Not Running

Ensure Go is in your PATH:

```bash
which cfn-lint
```

If not found, install cfn-lint-go:

```bash
go install github.com/lex00/cfn-lint-go/cmd/cfn-lint@latest
```

### GitHub Action Fails

Check that:
1. Go is set up with `actions/setup-go@v5`
2. Templates path is correct
3. Required permissions are set for SARIF upload

### Config File Not Found

cfn-lint-go searches from current directory up to git root. Verify:

```bash
cfn-lint template.yaml --config .cfnlintrc.yaml
```

## Future Integrations

Planned integrations:

- **VS Code Extension**: Real-time linting in the editor
- **IntelliJ IDEA Plugin**: Support for JetBrains IDEs
- **Terraform CloudFormation Provider**: Lint templates in Terraform workflows
