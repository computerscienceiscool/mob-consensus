# 005-JJ: Branch creation not idempotent

**Source:** consensus-bugs.md #3
**Severity:** High
**Category:** Flow Logic

## Problem

`runCreateBranch()` at `main.go:235` uses `git checkout -b` which fails if the
branch already exists. The flow is not idempotent — running `mob-consensus -b main`
twice errors out on the second run.

## Expected Behavior

If the user's branch already exists, switch to it (or offer to reset it) instead
of failing with a fatal error.

## Current Code

```go
// main.go:235
if err := gitRun(ctx, "checkout", "-b", newBranch, baseBranch); err != nil {
    return err
}
```

## Suggested Fix

Check if branch exists first. If it does, either:
- Switch to it with `git checkout <branch>` (safe, preserves work)
- Use `git checkout -B` to force-reset (destructive, needs confirmation)
- Print an informative error explaining the branch exists

## Subtasks

- [x] 005.1 Detect existing branch before `checkout -b`
- [x] 005.2 Choose behavior: switch vs error-with-hint vs reset (chose switch - safe, preserves work)
- [ ] 005.3 Add unit/integration test for re-run scenario
