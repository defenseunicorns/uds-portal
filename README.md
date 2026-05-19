# UDS Portal

<img align="right"  alt="Unicorn Delivery Service logo" src="ui/static/uds.svg"  height="256" />

[![UDS Documentation](https://img.shields.io/badge/docs-docs.defenseunicorns.com-775ba1)](https://docs.defenseunicorns.com/)

UDS Portal is the landing page for UDS users — a single point of discovery for every application deployed in a UDS environment. It ships as a UDS Core layer that depends on `base` and `identity-authorization`.

Apps are derived from `UDS Package` Custom Resources in the cluster. An app is shown if it has an `sso` section; if a `groups` section is also present, only members of those groups can see it.

<br><br>

## Prerequisites

Install [mise](https://mise.jdx.dev/getting-started.html) to manage development dependencies at the versions pinned in `mise.toml`, then install the pre-commit hook:

```sh
mise install
hk install
```

[uds-cli](https://github.com/defenseunicorns/uds-cli#install) is required to run tasks.

## Development

Spin up a local UDS Core k3d cluster with seed apps:

```sh
uds run dev-setup
```

Deploy portal changes to the running cluster:

```sh
uds run dev-deploy
```

Run the backend and frontend dev servers separately (useful for rapid UI iteration):

```sh
uds run dev-server   # Go API with hot-reload on :8080
uds run dev-ui       # Svelte dev server
```

## Building

```sh
uds run build:ui           # build Svelte frontend
uds run build:api          # build Go API for local platform
uds run build:container    # build container image (local arch)
uds run build:zarf-package # build Zarf package
```

## Testing

```sh
uds run test:unit       # all unit tests (backend + frontend)
uds run test:e2e-setup  # deploy test cluster
uds run test:e2e        # run e2e tests against running cluster
```

## Linting

```sh
uds run lint
```

## Authentication Model

UDS Portal runs behind AuthService (from UDS Core), the trust boundary for identity.

- AuthService mints and validates JWTs; the portal does not perform independent JWT signature verification.
- Group claims from the validated Keycloak JWT determine which apps each user can see.

## Deploying

For a test deployment:

```sh
uds run test:e2e-setup
```

For local dev clusters, UDS Portal is accessible at <https://portal.uds.dev>.

For production, UDS Portal is distributed as a UDS Core layer via OCI-hosted Zarf packages at `oci://ghcr.io/defenseunicorns/packages/uds`.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).
