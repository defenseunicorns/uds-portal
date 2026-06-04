# CLAUDE.md

## Sub-guides

Read the relevant sub-guide before modifying code of that type:

- [.ai/GO_CODE_STANDARDS.md](.ai/GO_CODE_STANDARDS.md) — any `.go` file
- [.ai/SVELTE_STANDARDS.md](.ai/SVELTE_STANDARDS.md) — any `.svelte` or `ui/src/**/*.ts` file
- [.ai/FRONTEND_TESTING.md](.ai/FRONTEND_TESTING.md) — any `*.test.ts` or `*.spec.ts` file

## Copyright

When modifying any file that has a copyright header, ensure the year includes the current year. If the header says `Copyright 2024` and the current year is 2026, update it to `Copyright 2024-2026`. If it already includes the current year, leave it as-is.

## Behavior

- Smallest change that solves the problem. Fix root causes, no temporary workarounds.
- Verify with tests, logs, or runtime output before declaring done.
- Ask before destructive or shared-state actions (cluster ops, force pushes, deletions).

## Tooling

- Lint: `uds run lint`
- Unit tests: `uds run test:unit`
- E2E tests: `uds run test:e2e`
