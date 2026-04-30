# AGENTS.md

## Purpose

UDS Portal is a lightweight end-user portal for UDS Core clusters. It discovers installed apps from Kubernetes `UDS Package` custom resources and only presents apps/endpoints the current user is allowed to access.

Core product behavior:
- Backend reads `uds.dev/v1alpha1` `packages` resources from the cluster.
- Backend filters apps by endpoint presence, auth context, and SSO group membership.
- Frontend shows a searchable app grid and links users to the permitted app endpoints.

## Architecture Snapshot

- Entrypoint: `main.go`
  - Embeds built UI assets from `ui/build/*`.
  - Starts API/UI server via `src/pkg/api`.
- API/router: `src/pkg/api/start.go`
  - Routes:
    - `GET /healthz`
    - `GET /api/v1/auth`
    - `GET /api/v1/apps`
  - Uses auth middleware (`src/pkg/api/middleware/auth.go`).
- App discovery/filter logic: `src/pkg/api/apps/apps.go`
  - Queries Kubernetes dynamic client for UDS Package CRs.
  - Filters to packages with endpoints and excludes `uds-portal`.
  - Applies user-group filtering using auth context (`incluster.GroupKey`).
  - Removes `.admin.` endpoints for non-admin users.
- Frontend: `ui/` (SvelteKit + Vite + Tailwind)
  - Main page fetches `GET /api/v1/apps` and renders app tiles.

## Important Directories

- `src/pkg/api/` - backend handlers, routing, middleware, auth, app filtering.
- `src/pkg/config/` - env-based runtime config.
- `ui/` - Svelte app, unit tests (Vitest), e2e tests (Playwright).
- `tasks/` + `tasks.yaml` - canonical build/test workflows (used by CI).
- `chart/` - Helm chart for deployment.
- `ui/tests/packages/` - integration/e2e deployment assets.
- `adr/` - design decisions and engineering rationale.

## Local Development Workflow

Use `uds run` tasks as the source of truth where possible.

Common commands:
- `uds run --list-all` - list all tasks.
- `uds run dev-server` - run backend locally.
- `uds run dev-ui` - run frontend locally.
- `uds run build:ui` - build frontend assets.
- `uds run build:api` - build backend binary.
- `uds run test:unit` - run backend + frontend unit tests.
- `uds run test:e2e` - run end-to-end tests.

Notes:
- Backend binary embeds `ui/build`; many Go test/build flows require running `build:ui` first.
- Local auth and in-cluster auth are controlled by env vars in `src/pkg/config/config.go`.

## Testing Expectations

When changing behavior, add or update tests close to the changed code.

- Backend unit tests: `go test ./src/pkg/...` (via `uds run test:go-unit`).
- Frontend unit tests: `pnpm run test:unit` in `ui/` (via `uds run test:ui-unit`).
- E2E tests: Playwright suites in `ui/tests/` (via `uds run test:e2e*`).

Preferred validation order for iterative work:
1. Smallest relevant unit tests.
2. Broader package/unit suite.
3. E2E tests for user-visible flows.

## Coding Guidance For Agents

- Keep changes focused and minimal; preserve existing patterns unless intentionally refactoring.
- Fix root causes, not symptoms.
- Do not change public API behavior without updating tests and docs.
- Do not commit generated artifacts unless the repo convention requires it.
- Avoid introducing new dependencies unless necessary and justified.
- Preserve AGPL/commercial dual-license headers in source files where already present.

For app visibility/auth changes specifically:
- Treat `src/pkg/api/apps/apps.go` as the source of truth for package filtering.
- Be explicit about admin group bypass behavior.
- Ensure endpoint filtering and group filtering semantics remain test-covered.

## Definition of Done For Changes

A change is considered done when:
- Relevant unit/integration tests pass.
- Behavior is validated in local dev flow.
- Docs are updated when behavior or workflow changes.
- No unrelated files are modified.
