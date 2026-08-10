# Frontend Testing

Follow the [Vitest docs](https://vitest.dev) and [Playwright docs](https://playwright.dev). This file captures project-specific caveats and non-obvious behavior.

## Vitest (Unit/Component)

Run: `uds run test:unit` or `pnpm -C ui run test:unit`

**SvelteKit runtime is absent in Vitest.** `$app/*` modules don't exist — mock them. This project uses the legacy `$app/stores` API (not `$app/state`):

```ts
import { readable } from 'svelte/store'

vi.mock('$app/stores', () => ({
  page: readable({ url: new URL('http://localhost/') }),
}))

vi.mock('$app/paths', () => ({
  resolve: (p: string) => p,
}))
```

**Tailwind does not compile in Vitest.** `element.toBeVisible()` misreads elements hidden via `hidden`, `opacity-0`, or `invisible` — assert class presence instead:

```ts
function isHidden(el: HTMLElement) {
  return el.classList.contains('hidden') || !!el.closest('.hidden')
}
expect(isHidden(screen.getByTestId('my-element'))).toBe(true)
```

## Playwright (E2E)

Run: `uds run test:e2e`

- Tests run in parallel (`fullyParallel: true`) — every test must be race-condition-safe.
- Navigate explicitly at the start of every test (`test.beforeEach`). Never assume navigation state from a prior test.
- Use `test.describe.configure({ mode: 'serial' })` only for intra-file sequential dependencies; prefer isolation fixes.
