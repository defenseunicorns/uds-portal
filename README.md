# UDS Portal

<img align="right"  alt="Unicorn Delivery Service logo" src="ui/static/uds.svg"  height="256" />

[![UDS Documentation](https://img.shields.io/badge/docs-docs.defenseunicorns.com-775ba1)](https://docs.defenseunicorns.com/)

UDS Portal is the landing page for UDS users — a single point of discovery for every application deployed in a UDS environment. It ships as a UDS Core layer that depends on `base` and `identity-authorization`.

Apps are derived from `UDS Package` Custom Resources in the cluster. One tile is created per `network.expose` entry; a package must also have an `sso` section or its tiles are not shown. Group-based access works per `sso` entry via `spec.sso[].groups.anyOf`:

- If `groups.anyOf` lists one or more groups, only members of those groups see the tiles.
- If the `sso` entry omits `groups` (or leaves `anyOf` empty), the tiles are visible to all authenticated users.

```yaml
spec:
  sso:
    # Visible only to members of the listed groups
    - name: Admin App
      clientId: uds-admin-app
      groups:
        anyOf:
          - /UDS Core/Admin
    # No groups block: visible to all authenticated users
    - name: Shared App
      clientId: uds-shared-app
```

App tiles display the package name as a title and a default logo unless overridden. To customize, set annotations in the Zarf package metadata:

```yaml
metadata:
  annotations:
    dev.uds.title: My App        # display name shown on the tile
    dev.uds.icon: 'data:image/svg+xml;base64,<base64-encoded-svg>'  # tile logo
```

If `dev.uds.title` is absent, the title falls back to a formatted version of the package name (e.g. `uds-registry` → `UDS Registry`). If `dev.uds.icon` is absent, a default logo is shown.

Packages that [expose multiple endpoints](https://docs.defenseunicorns.com/core/reference/operator--crds/packages-v1alpha1-cr/#network) (e.g. GitLab) can control which ones appear as tiles:

- Expose entries with a wildcard `host` (e.g. `*.pages`) are always excluded.
- To hide specific endpoints, set the `uds.dev/portal-hide-apps` annotation on the Package CR to a comma-separated list of `network.expose[].host` values:

```yaml
metadata:
  annotations:
    uds.dev/portal-hide-apps: "registry,pages"
```

<br><br>

## Part of UDS Core

UDS Portal is a [UDS Core](https://github.com/defenseunicorns/uds-core) functional layer (`core-portal`), not a standalone application. It depends on the `base` and `identity-authorization` layers.

UDS Core consumes this repository as a git-based Helm chart, referenced by tag from [`src/portal/common/zarf.yaml`](https://github.com/defenseunicorns/uds-core/blob/main/src/portal/common/zarf.yaml) and kept in sync with releases here by Renovate. The layer ships in the `upstream` and `unicorn` flavors only, not `registry1`.

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
