# 006-JJ: .lock file not in .gitignore

**Source:** consensus-bugs.md #7
**Severity:** High
**Category:** Repository Hygiene

## Problem

`.lock` files appear in diffs, `git status`, and ahead/behind calculations.
The `.gitignore` does not exclude them.

## Impact

- False divergence signals in discovery mode
- Polluted diff output
- Incorrect status reporting
- Noise in merge operations

## Current .gitignore

Does not contain any `.lock` entry.

## Suggested Fix

Add `*.lock` to `.gitignore`.

## Subtasks

- [x] 006.1 Add `*.lock` to `.gitignore`
- [x] 006.2 Remove any tracked `.lock` files from the index if present (none found)
