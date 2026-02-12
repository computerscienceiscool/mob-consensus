# 013-JJ: Clean merge skips review and commit

**Source:** consensus-bugs.md #8 - "Diff tool invocation unreliable"
**Severity:** High
**Category:** Core UX Bug

## Problem

`runMerge()` at `main.go:476-481` has a critical early return that skips
difftool and commit when a merge succeeds without conflicts.

```go
if _, err := os.Stat(mergeHeadPath); err != nil {
    if errors.Is(err, os.ErrNotExist) {
        return nil  // ← EXITS EARLY!
    }
    return err
}
```

When `git merge --no-commit --no-ff` succeeds cleanly (no conflicts), Git doesn't
create MERGE_HEAD. The code interprets this as "nothing to do" and returns,
skipping:
- Line 495: `git difftool -t vimdiff HEAD` (review changes)
- Line 499: `git commit -e -F <msgPath>` (finalize merge)

## Expected Behavior

The tool should always:
1. Run merge with `--no-commit --no-ff` ✅
2. If conflicts: open mergetool ✅
3. Open difftool for review (ALWAYS, even if no conflicts) ❌
4. Open editor to finalize commit message ❌
5. Commit the merge ❌
6. Push ✅

## Impact

- Core merge UX broken for clean merges
- User must manually `git difftool HEAD` and `git commit`
- Contradicts tool's automation promise
- Demo showed this exact failure

## Root Cause

The code assumes MERGE_HEAD only exists when there are conflicts. Actually:
- MERGE_HEAD exists when `--no-commit` leaves merge uncommitted
- MERGE_HEAD does NOT exist when merge is already committed OR when git
  determines it can't do a merge at all

The check at line 476 is trying to detect "did the merge start?" but it's
checking the wrong thing.

## Suggested Fix

The logic at lines 476-481 should be removed or completely rethought. After
mergetool (line 471), we should ALWAYS proceed to difftool and commit, because
we explicitly used `--no-commit`.

Possible approaches:
1. Remove the MERGE_HEAD check entirely - always run difftool and commit
2. Check `git status --porcelain` to see if there are staged changes
3. Check exit code of the merge command more carefully

## Subtasks

- [ ] 013.1 Understand when MERGE_HEAD exists vs doesn't exist
- [ ] 013.2 Remove or fix the early return at lines 476-481
- [ ] 013.3 Ensure difftool and commit always run after merge (clean or conflicted)
- [ ] 013.4 Test with both clean merges and conflicted merges
