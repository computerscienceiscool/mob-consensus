# 008-JJ: Merge doesn't auto-create local tracking branch

**Source:** consensus-bugs.md #5
**Severity:** Critical
**Category:** Automation Gap

## Problem

`runMerge()` at `main.go:379` runs `git merge --no-commit --no-ff <otherBranch>`
directly. If the user passes a branch name like `alice/twig` but only
`remotes/origin/alice/twig` exists (no local tracking branch), the merge fails
with a confusing git error.

The user must manually run:
```
git switch -c alice/twig origin/alice/twig
git switch <their-branch>
mob-consensus alice/twig
```

## Expected Behavior

If `otherBranch` doesn't resolve locally, check if `origin/<otherBranch>` (or
another remote) has it and either:
- Use the remote ref directly in the merge command
- Auto-create a local tracking branch first

## Impact

- Core merge flow broken for the most common case (merging a collaborator's
  remote branch)
- Requires manual git intervention, defeating the tool's purpose

## Subtasks

- [x] 008.1 Before merge, check if otherBranch resolves as a ref
- [x] 008.2 If not, try `origin/<otherBranch>` or search remotes
- [x] 008.3 Either use remote ref directly or create local tracking branch (chose remote ref directly - no mutation)
- [ ] 008.4 Test with remote-only branches
