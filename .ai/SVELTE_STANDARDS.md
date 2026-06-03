# Svelte/SvelteKit Standards

Follow the [Svelte docs](https://svelte.dev/docs/svelte) and [SvelteKit docs](https://svelte.dev/docs/kit). This file captures project-specific rules and non-obvious pitfalls.

## Reactivity

This project uses Svelte 4 syntax on Svelte 5 (backward-compatible mode) with the legacy `svelte/store` API.

- Props declared with `export let`. Type with TypeScript where useful.
- Reactive statements use `$:` — use for derived values and side effects tied to reactive state.
- Store subscriptions use the `$store` auto-subscription shorthand in `.svelte` files.
- Use `writable` for mutable state.

## SvelteKit

- This project uses `ssr = false` (browser-only). All load functions are universal (`+layout.ts`), no server load files.
- Don't `await parent()` before fetching independent data — causes unnecessary waterfalls.
- Don't use `redirect()` inside `try` blocks — it throws and will be caught.

## TypeScript

- Type component props explicitly. Prefer interfaces over inline types for reusable shapes.
- Prefer explicit types at component boundaries; let inference handle internals.

## Tailwind

- Utility classes only — no custom CSS unless a utility genuinely can't express it.
- Responsive variants (`sm:`, `md:`) and state variants (`hover:`, `focus:`) over JS-driven style logic.
- Dynamic runtime values (e.g. colors from a map, calculated heights) may use inline `style=` — acceptable when Tailwind arbitrary values can't represent the value.
