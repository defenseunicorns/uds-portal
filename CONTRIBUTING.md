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

### Local Development Cluster

Use `dev-setup` when you want a full local UDS Core cluster with a variety of real UDS apps for manual portal development:

```sh
uds run dev-setup
```

This task:

- Deploys the full UDS Core k3d dev bundle via `setup:k3d-full-cluster`.
- Creates the standard local Keycloak test user via `setup:create-doug-user`.
- Creates and deploys seed bundles from `seed/bundles/`.
- Deploys shared dependencies first, then app bundles in parallel.

GitLab is not deployed by default. Enable it only when you need a complex representative app with multiple exposed entries to work with; otherwise skip it because it takes a long time to deploy and uses a lot of local resources:

```sh
uds run dev-setup --with gitlab=true
```

Use `deploy-seed-bundles` when the cluster already exists and you only want to rebuild and redeploy the seed bundles:

```sh
uds run deploy-seed-bundles
```

To include GitLab in a seed-only redeploy:

```sh
uds run deploy-seed-bundles --with gitlab=true
```

Use `dev-deploy` for quick portal iteration after the dev cluster is up.

```sh
uds run dev-deploy
```

This builds the portal container image for your local architecture, creates a local Zarf package with the same image tag, and deploys that package to the current cluster.

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

### Release Process

Releases are automated via [release-please](https://github.com/googleapis/release-please). Merge conventional commits to `main`; release-please opens a release PR with an updated `CHANGELOG.md` and bumped versions. Merging that PR creates the tag and triggers the publish workflow.

**What gets published on release:**

| Artifact | Destination |
|---|---|
| Docker image (linux/amd64 + linux/arm64) | `ghcr.io/defenseunicorns/uds-portal:<version>` |
| Zarf package amd64 | `oci://ghcr.io/defenseunicorns/packages/uds` |
| Zarf package arm64 | `oci://ghcr.io/defenseunicorns/packages/uds` |

**Version files:** release-please bumps versions in lockstep across `chart/Chart.yaml`, `tasks/build.yaml`, and `tasks/release.yaml`.

The container image tag and Zarf package version are driven by the `VERSION` task variable, matching the UDS Registry pattern. Local `dev-deploy` defaults `VERSION` to `0.0.0-dev` so it does not collide with a published release image. It builds the local container image, creates a local Zarf package with that image tag, then deploys the package tarball with `--set VERSION=${VERSION}`. Release publishes pass the release-please version as `VERSION`. The chart image tag uses `###ZARF_CONST_VERSION###`, which is backed by the package version constant in `zarf.yaml`.

**Override version for a release:**

Add `Release-As: x.y.z` to any commit body landing on `main`:

```sh
git commit --allow-empty -m "chore: release 0.2.0" -m "Release-As: 0.2.0"
```
