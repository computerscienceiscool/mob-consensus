# Repository Guidelines

## Project Structure & Module Organization
- The Go CLI lives at the module root (`main.go`) (preferred).
- `x/` holds experimental prototypes and legacy helpers; `x/mob-consensus` is the original Bash implementation kept for reference.
- Keep packages at the module root or under `x/` (avoid `internal/` and `pkg/`).
- Do not commit local state (e.g., `.grok/`, `.grok`) or generated binaries; keep these ignored locally.

## Build, Test, and Development Commands
- `go test ./...`: run unit tests.
- `go run . -h`: run locally (requires Go 1.24.0+).
- `go install github.com/stevegt/mob-consensus@latest`: install/upgrade the CLI.
- `mob-consensus -h`: show help and flags.

Notes: the tool runs `git fetch`, uses `git mergetool`/`git difftool` (defaulting to `vimdiff`), and pushes after commits unless `-n` is set.

## Coding Style & Naming Conventions
- Bash: keep changes small and readable; validate with `bash -n x/mob-consensus` (and `shellcheck x/mob-consensus` if available).
- Go: run `gofmt`; keep package names short and lower-case. Minimum supported Go is 1.24.0 (see `go.mod`).
- Always add detailed comments to code.

## Testing Guidelines
- Prefer deterministic tests using Go’s standard `testing` package when adding Go code.
- When tests interact with Git workflows, keep them as realistic as possible by using the same commands shown in the `mob-consensus` help (`usage.tmpl`) where practical (e.g., `git switch -c`, `git fetch`, `git push -u`). If a test must deviate (compatibility, determinism, or focus), explain why in the test code comments.
- For shell changes, run `bash -n x/mob-consensus` and include a short manual repro in the PR description.

## TODO Tracking

Track work in `TODO/` with an index at `TODO/TODO.md`; number using letter-prefixed IDs; don't renumber; sort by priority.

### Determining YOUR prefix

Run `git config user.name` to see your git username. Your TODO prefix is the **FIRST LETTER** (uppercase) of that name.

Examples:
- `git config user.name` → "Alice" → use prefix **A** (create A001, A002, etc.)
- `git config user.name` → "Bob" → use prefix **B** (create B001, B002, etc.)
- `git config user.name` → "Sam" → use prefix **S** (create S001, S002, etc.)

### TODO ID format

- Use `LNNN` format: letter prefix + 3 digits (e.g., A015, B016, S023)
- The prefix is based on who **created** that TODO entry
- Keep integer parts globally unique during transition

### Transition from old format

- When bulk-renaming existing TODO files to add prefixes, use `git mv` in one commit without mixing other work

### Usage

- Mark completion with checkboxes: `- [ ] A005 - ...` → `- [x] A005 - ...`
- Legacy: root `TODO.md` exists for historical reference; update `TODO/TODO.md` going forward

## Commit & Pull Request Guidelines
- Keep commit messages short and imperative; existing history often uses a `mob-consensus:` prefix for script changes.
- PRs should include: a concise summary, test commands run (e.g., `bash -n x/mob-consensus`), and before/after notes for behavior or output changes.
- When staging, list files explicitly (avoid `git add .` / `git add -A`).

## Agent-Specific Notes
- Check `~/.codex/AGENTS.md` for updated local workflows and keep `~/.codex/meta-context.md` current.
- Treat a line containing only `commit` as “add and commit all changes with an AGENTS.md-compliant message”.
