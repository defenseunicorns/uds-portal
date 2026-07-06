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

### Package CR annotations

Per-endpoint annotations on `spec.network.expose[]` entries control tile appearance and visibility. `portal.uds.dev/` takes precedence over `uds.dev/` — use the portal namespace to override upstream package defaults at the bundle level.

| | `uds.dev/` | `portal.uds.dev/` |
|---|---|---|
| Title | `uds.dev/title` | `portal.uds.dev/title` |
| Icon | `uds.dev/icon` | `portal.uds.dev/icon` |
| Visibility | `uds.dev/visible` | `portal.uds.dev/visible` |

```yaml
spec:
  network:
    expose:
      - host: my-app
        annotations:
          uds.dev/title: "My App"
          uds.dev/icon: "data:image/svg+xml;base64,<base64-encoded-svg>"
      - host: my-app-admin
        annotations:
          uds.dev/visible: "false"  # hidden; use portal.uds.dev/visible to override at bundle level
```

- **Title** falls back to `dev.uds.title` in Zarf package metadata, then the formatted package name (e.g. `uds-registry` → `UDS Registry`).
- **Icon** falls back to `dev.uds.icon` in Zarf package metadata, then a default logo.
- **Visibility** — missing or any value other than `"false"` (case-insensitive) defaults to visible. Wildcard hosts (e.g. `*.pages`) are always excluded. To hide endpoints by host name across a package, set `uds.dev/portal-hide-apps` on the Package CR metadata to a comma-separated list of host values.

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
