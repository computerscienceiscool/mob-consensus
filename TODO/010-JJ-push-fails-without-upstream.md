# 010-JJ: Push may fail without upstream tracking

**Source:** consensus-bugs.md #9
**Severity:** Medium
**Category:** Workflow Logic

## Problem

`runMerge()` at `main.go:424` runs bare `git push` after a successful commit.
If the branch doesn't have an upstream tracking ref set, this fails.

The branch creation flow (`runCreateBranch`) only *prints advice* to push
(line 240) but doesn't actually push or set upstream. So a user who creates a
branch with `-b`, commits, and then merges will hit a push failure.

The old Bash version (line 97) did `git push --set-upstream origin $new_branch`
during branch creation, which set up tracking automatically.

## Impact

- Inconsistent sync state
- Users confused about why push fails after a successful merge
- Must manually set upstream before mob-consensus push works

## Related

`ensureClean()` at `main.go:453` also does bare `git push` with the same issue.

## Subtasks

- [x] 010.1 Detect if upstream is set before pushing
- [x] 010.2 If not set, use `git push -u <remote> <branch>` with best remote
- [x] 010.3 Or: prompt user / print guidance on first push failure (chose auto-fix approach)
