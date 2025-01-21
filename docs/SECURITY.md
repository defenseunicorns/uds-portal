# Security

This document describes the security posture and related security automation for UDS Runtime.

## Static Code Analysis

This project uses a variety of static code analysis tools to ensure the quality and security of the codebase, the specific tools are detailed below.

### Golang

`golangci-lint` is used to perform static code analysis on the Go backend code. This tool is configured to use a variety of linters including `gosec`, `govet` and many others to ensure the security of the codebase. The configuration for this tool can be found in [.golangci.yml](../.golangci.yml).

### Typescript

`prettier`, `eslint` and `svelte-check` are used to perform static code analysis on the frontend Typescript code. The configuration for these tools can be found in [.eslintrc.js](../ui/eslint.config.js), [.prettierrc.js](../ui/.prettierrc.cjs) and [tsconfig.json](../ui/tsconfig.json).

### YAML

`yamllint` is used to perform static code analysis on the YAML files in the project. The configuration for this tool can be found in [.yamllint](../.yamllint).


### Pre-Commit and CI

All of the above tools are run as part of the pre-commit hooks and CI pipelines to ensure issues are detected early and only code that passes the analyzers can be merged into `main`.

## Dependency Scanning

This project uses `grype` to scan both the Go and Node.js dependencies for vulnerabilities in CI. The configuration for this tool can be found in [.grype.yaml](../.grype.yaml) and [dependency-check.yaml](../.github/workflows/dependency-check.yaml).

## Supply Chain Security

SBOMS for this project are generated using Zarf; they can be viewed using the Zarf SBOM viewer.
