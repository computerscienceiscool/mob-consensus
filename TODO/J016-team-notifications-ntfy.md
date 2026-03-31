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

## Integration plan (pending architect review)

To wire in notifications, three changes are needed:

### 1. Add notifier to options struct (main.go)

Add a `notifier Notifier` field to the `options` struct.

### 2. Create notifier in CLI handlers (cli.go)

In each Cobra `RunE` handler that calls runMerge/runStart/runJoin, set:

```go
opts.notifier = newNotifier(cmd.Context())
```

### 3. Add Notify calls at each event site (main.go)

**a) runMerge — after successful smartPush (line ~1251):**

```go
opts.notifier.Notify(ctx, Event{
    Type:    EventMerged,
    User:    user,
    Twig:    twigFromBranch(currentBranch),
    Branch:  currentBranch,
    Message: fmt.Sprintf("%s merged %s and pushed", user, mergeTarget),
})
```

**b) runStart — after runGitPlan returns successfully (line ~905):**

```go
opts.notifier.Notify(ctx, Event{
    Type:    EventStarted,
    User:    user,
    Twig:    twig,
    Branch:  userBranch,
    Message: fmt.Sprintf("%s started twig %s", user, twig),
})
```

**c) runJoin — after runGitPlan returns successfully (line ~1007):**

```go
opts.notifier.Notify(ctx, Event{
    Type:    EventJoined,
    User:    user,
    Twig:    twig,
    Branch:  userBranch,
    Message: fmt.Sprintf("%s joined twig %s", user, twig),
})
```

**d) Auto-commit callers — after ensureClean + smartPush succeeds:**

The caller (not ensureClean itself) notifies with `EventAutoCommitted`,
since each caller knows its own context for the message.

### Note

The notifier is a no-op when `mob-consensus.ntfyTopic` is not set in
git config, so wiring the calls in has zero effect until configured.

## Subtasks

- [x] J016.1 — Notifier interface + ntfyNotifier implementation (cswg 020.2) — done in notify.go (b6228af)
  - [x] J016.1.1 — Define `Notifier` interface and `Event` type
  - [x] J016.1.2 — Implement `ntfyNotifier`: read git config, async POST with 5s timeout
  - [x] J016.1.3 — Constructor returns ntfyNotifier or no-op based on config
- [x] J016.2 — Add notification call sites (cswg 020.3) — done in cli.go + main.go (5a9338f)
  - [x] J016.2.1 — After merge + push in `runMerge`
  - [x] J016.2.2 — After `runStart` completes
  - [x] J016.2.3 — After `runJoin` completes
  - [x] J016.2.4 — After auto-commit + push in `ensureClean`
- [x] J016.3 — Unit tests (cswg 020.4) — done in notify_test.go (a3fd264)
  - [x] J016.3.1 — No-op when topic is not configured (TestNewNotifier/NoTopic)
  - [x] J016.3.2 — No-op when `ntfyEnabled` is `false` (TestNewNotifier/Disabled)
  - [x] J016.3.3 — Sends correct POST when configured (TestNtfyNotifierSendsRequest, uses httptest)
  - [x] J016.3.4 — Empty message falls back to event type (TestNtfyNotifierEmptyMessageFallback)
  - [x] J016.3.5 — Default server is ntfy.sh (TestNewNotifier/DefaultServer)
  - [x] J016.3.6 — Custom server from git config (TestNewNotifier/WithTopicAndServer)
- [x] J016.4 — Document configuration in `usage.tmpl` (cswg 020.6) — done (12b360b), pending Steve's review
- [ ] J016.5 — Git branch overlap detection in discovery (cswg 020.12) — future, depends on overlap detection design
- [ ] J016.6 — Repo-tracked config (cswg 020.9) — future, blocked on mob-consensus TODO 008
