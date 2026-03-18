# J016 - Team notifications via ntfy

Implements the mob-consensus portion of cswg coordination repo TODO 020.

## Summary

Add ntfy push notifications to mob-consensus so collaborators get async
push notifications when workflow events happen (merges, joins, auto-commits)
without requiring Discord, SMS, email, or any accounts/API keys.

## Architecture

- **Notifier interface** with a `Notify(ctx, Event)` method
- **ntfyNotifier** concrete implementation: reads `git config --local`, sends
  async HTTP POST via goroutine, fire-and-forget (never blocks or errors)
- **No-op notifier** returned when no topic is configured
- **Event struct** holds event type, user, branch, twig, human-readable message
- **Config via `git config --local`** (never committed):
  - `mob-consensus.ntfyTopic` — required, no default
  - `mob-consensus.ntfyServer` — optional, defaults to `https://ntfy.sh`
  - `mob-consensus.ntfyEnabled` — optional, disable without removing topic
- **No new Go dependencies** — uses `net/http` from stdlib
- **No new CLI flags for MVP** — config only

## Notification insertion points

| Call site | Function | Event |
|-----------|----------|-------|
| After merge + push | `runMerge()` main.go:1251 | "alice merged bob/feature-x and pushed" |
| After start completes | `runStart()` main.go:905 | "alice started twig feature-x" |
| After join completes | `runJoin()` main.go:1007 | "bob joined twig feature-x" |
| After auto-commit + push | caller of `ensureClean()` | "alice committed dirty changes and pushed" |

Note: notifications go in the *callers* of `ensureClean()`, not inside it,
so each caller has proper context for the message.

## Subtasks

- [ ] J016.1 — Notifier interface + ntfyNotifier implementation (cswg 020.2)
  - [ ] J016.1.1 — Define `Notifier` interface and `Event` type
  - [ ] J016.1.2 — Implement `ntfyNotifier`: read git config, async POST with 5s timeout
  - [ ] J016.1.3 — Constructor returns ntfyNotifier or no-op based on config
- [ ] J016.2 — Add notification call sites (cswg 020.3)
  - [ ] J016.2.1 — After merge + push in `runMerge`
  - [ ] J016.2.2 — After `runStart` completes
  - [ ] J016.2.3 — After `runJoin` completes
  - [ ] J016.2.4 — After auto-commit + push (callers of `ensureClean`)
- [ ] J016.3 — Unit tests (cswg 020.4)
  - [ ] J016.3.1 — No-op when topic is not configured
  - [ ] J016.3.2 — No-op when `ntfyEnabled` is `false`
  - [ ] J016.3.3 — Sends correct POST when configured (use `httptest.NewServer`)
- [ ] J016.4 — Document configuration in `usage.tmpl` (cswg 020.6)
- [ ] J016.5 — Git branch overlap detection in discovery (cswg 020.12) — future, depends on overlap detection design
- [ ] J016.6 — Repo-tracked config (cswg 020.9) — future, blocked on mob-consensus TODO 008
