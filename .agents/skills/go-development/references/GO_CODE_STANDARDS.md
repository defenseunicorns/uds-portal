# Go Code Standards

Follow [Effective Go](https://go.dev/doc/effective_go) and the [Google Go Style Guide](https://google.github.io/styleguide/go/). This file captures only project-specific rules and non-obvious conventions.

## Project Rules

- Lint: `uds run lint`. Suppress narrowly with `//nolint:linter // reason` — never disable globally.
- Run `go mod tidy` after dependency changes.
- `log/slog` for structured logs. Pass key/value pairs, not formatted strings. Never log secrets.
- Guard clauses over nested happy paths — handle the error/empty case first.
- Group related `const (...)` and `var (...)` declarations.
- Return errors; do not panic except at `main` init or truly unrecoverable invariants.
- Wrap errors with context: `fmt.Errorf("thing %s: %w", name, err)`. Never log and return the same error.
- `context.Context` is always the first parameter, named `ctx`. Never store in a struct.
- `t.Parallel()` for independent tests. `t.Helper()` on test helpers. Table-driven with subtests.
