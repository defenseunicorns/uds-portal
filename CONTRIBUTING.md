# Contributing to UDS Portal

Welcome :unicorn: to UDS Portal!

To report a bug, request a feature, or ask a question, review [open issues](/issues) or open a [new issue](/issues/new/choose).

## Submissions

Track your work in Linear. Reference the Linear issue (e.g. `CORE-123`) in your PR description.

Recommended workflow:

1. Clone the repo
2. Create a feature branch from `main`
3. Set up your environment (see [Prerequisites](#prerequisites))
4. Make your changes (add tests and docs as appropriate)
5. Open a PR against `main`

## Prerequisites

Install [mise](https://mise.jdx.dev/getting-started.html) to manage development dependencies. Recommended: enable [shell activation](https://mise.jdx.dev/getting-started.html#activate-mise).

## Setup

```sh
# install pinned tools (mise.toml)
mise install

# install pre-commit hook (hk.pkl)
hk install
```

## Tasks

```sh
# run all linters (matches CI)
uds run lint:check

# auto-fix where possible
uds run lint:fix
```

## Technical Standards

- **Testing**: All features and bug fixes require automated tests unless blocked by a technical issue.
- **Readability**: Meaningful names, reasonable function sizes, simple solutions.
- **Design Documentation**: Use [Architectural Decision Records](https://adr.github.io/) for significant design decisions.

### Pre-Commit Hooks and Linting

This project uses [hk](https://hk.jdx.dev/) for pre-commit hooks to automate validation and linting. Follow [hk's Getting Started guide](https://hk.jdx.dev/getting_started.html) to install the git hook.

CI runs the same checks via `uds run lint:check`, so passing hk locally means passing CI lint.

### Continuous Delivery

Continuous Delivery is core to our development philosophy. See [https://continuousdelivery.com/principles/](https://continuousdelivery.com/principles/).

Specifically:

- Trunk-based development on `main` with short-lived feature branches that merge and delete after merge
- Do not merge code into `main` that isn't releasable
- Automated testing on all changes before merge
- Immutable release artifacts
