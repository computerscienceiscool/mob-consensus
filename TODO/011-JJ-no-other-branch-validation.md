# 011-JJ: No validation of OTHER_BRANCH before merge

**Source:** Code review (new finding)
**Severity:** Medium
**Category:** UX / Error Handling

## Problem

`parseArgs()` at `main.go:113-114` accepts any string as `otherBranch`:

```go
if len(rest) > 0 {
    opts.otherBranch = rest[0]
}
```

No validation is performed before passing it to `git merge`. If the user
mistypes a branch name, they get a raw git error rather than a helpful message.

## Expected Behavior

Before starting the merge flow, verify that `otherBranch` resolves to a valid
ref (local or remote). Print a clear error with suggestions (e.g., "did you
mean `origin/alice/twig`?") if it doesn't exist.

## Impact

- Confusing raw git errors for typos
- No guidance toward the correct branch name
- Related to 008-JJ (remote branch resolution)

## Resolution

Fixed by 008-JJ. The branch resolution code added in `runMerge()` at lines 362-371
validates that the branch exists (locally or as origin/<branch>) before attempting
the merge. Clear error message is provided if branch is not found.

## Subtasks

- [x] 011.1 Add ref validation before entering merge flow
- [ ] 011.2 On failure, list similar branch names as suggestions (deferred - optional enhancement)
- [x] 011.3 Consider: integrate with 008-JJ remote branch auto-resolution
