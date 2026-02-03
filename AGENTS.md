# Repository Guidelines

## Project Structure & Module Organization
- `x/` holds experimental prototypes. The current tool lives at `x/mob-consensus`.
- If you add Go code, keep packages at the module root or under `x/` (avoid `internal/` and `pkg/`).
- Do not commit local state (e.g., `.grok/`, `.grok`) or generated binaries; keep these ignored locally.

## Build, Test, and Development Commands
- `./x/mob-consensus -h`: show help and flags.
- `./x/mob-consensus`: compare the current branch with related branches matching `*/<twig>` (where `<twig>` is the basename of the current branch).
- `./x/mob-consensus OTHER_BRANCH`: merge `OTHER_BRANCH` onto the current branch with a prepared commit message and co-author lines.
- `./x/mob-consensus -b <base>`: create `$USER/<twig>` from `<base>` and push upstream.
- `./x/mob-consensus -c`: commit existing uncommitted changes before proceeding.
- `./x/mob-consensus -F`: force run even if not on a `$USER/` branch.

Notes: the script runs `git fetch`, uses `git mergetool`/`git difftool` (defaulting to `vimdiff`), and pushes after commits unless `-n` is set.

## Coding Style & Naming Conventions
- Bash: keep changes small and readable; validate with `bash -n x/mob-consensus` (and `shellcheck x/mob-consensus` if available).
- Go (if added): run `gofmt`; keep package names short and lower-case.

## Testing Guidelines
- Prefer deterministic tests using Go’s standard `testing` package when adding Go code.
- For shell changes, run `bash -n x/mob-consensus` and include a short manual repro in the PR description.

## TODO Tracking
- Track work in `TODO/` and keep an index at `TODO/TODO.md`.
- Number TODOs with 3 digits (e.g., `005`), do not renumber, and sort the index by priority (not number).
- In each `TODO/*.md`, use numbered checkboxes like `- [ ] 005.1 describe subtask`.

## Commit & Pull Request Guidelines
- Keep commit messages short and imperative; existing history often uses a `mob-consensus:` prefix for script changes.
- PRs should include: a concise summary, test commands run (e.g., `bash -n x/mob-consensus`), and before/after notes for behavior or output changes.
- When staging, list files explicitly (avoid `git add .` / `git add -A`).

## Agent-Specific Notes
- Check `~/.codex/AGENTS.md` for updated local workflows and keep `~/.codex/meta-context.md` current.
- Treat a line containing only `commit` as “add and commit all changes with an AGENTS.md-compliant message”.
