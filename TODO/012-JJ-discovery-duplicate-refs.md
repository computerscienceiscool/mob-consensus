# 012-JJ: Discovery shows duplicate local and remote refs

**Source:** Code review (new finding)
**Severity:** Low
**Category:** UX / Status Logic

## Problem

`runDiscovery()` uses `git branch -a` (line 295) which returns both local and
remote-tracking refs. `relatedBranches()` (line 460) includes all matches, so
the output shows both:

```
                        alice/twig  is synced
       remotes/origin/alice/twig  is synced
```

This is redundant when a local tracking branch exists for the remote ref.

## Impact

- Noisy output in discovery mode
- Confusing for users who see the same branch listed twice
- Related to consensus-bugs.md #10 (ahead/behind noise)

## Suggested Fix

Either:
- Prefer local branches and skip `remotes/` entries when a local equivalent exists
- Or: deduplicate by twig owner (group `alice/twig` and `remotes/origin/alice/twig`)

## Subtasks

- [ ] 012.1 Detect when local branch tracks a remote ref
- [ ] 012.2 Deduplicate or group output accordingly
