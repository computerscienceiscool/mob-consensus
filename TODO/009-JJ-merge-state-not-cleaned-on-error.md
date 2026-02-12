# 009-JJ: Merge state not cleaned up on error paths

**Source:** consensus-bugs.md #6, #8
**Severity:** High
**Category:** Workflow Bug

## Problem

In `runMerge()` (main.go:338-425), if any step after `git merge --no-commit`
fails, the repository is left in a "still merging" state with MERGE_HEAD present.

Error paths that leave dirty merge state:
- `mergetool` fails or user aborts (line 387-389)
- `difftool` fails or user aborts (line 411-413)
- `commit` fails (line 415-418) — prints "don't forget to push" but doesn't
  clean up merge state
- Any panic or signal kill

The user must manually run `git merge --abort` or `git reset --merge` to recover.

## Expected Behavior

On error after merge has started, either:
- Automatically run `git merge --abort` to restore clean state
- Print clear instructions: "repo is in merge state; run `git merge --abort` to
  undo or `git commit` to finalize"

## Impact

- Repository left dirty and confusing
- Subsequent git operations behave unexpectedly
- New users won't know how to recover

## Subtasks

- [ ] 009.1 Add deferred merge abort on error after merge starts
- [ ] 009.2 Or: print clear recovery instructions on each error path
- [ ] 009.3 Consider: should commit failure abort the merge or preserve it?
