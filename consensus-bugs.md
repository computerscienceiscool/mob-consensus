# consensus-bugs.md

## Mob Consensus -- Bugs Identified During Live Demo

Date: 2026-02-05\
Context: Live walkthrough of Go implementation of `mob-consensus` tool.

------------------------------------------------------------------------

## 1. Username Not Pulled from Local Git Config

**Expected Behavior:**\
Branch prefix should be derived from the local git config: -
`git config user.email` - Username = portion left of `@`

**Observed Behavior:**\
Incorrect username was used when creating personal branch.

**Impact:**\
- Incorrect branch naming - Multi-user flow compromised

**Category:** Core Logic Bug\
**Severity:** Critical

------------------------------------------------------------------------

## 2. Wrong Binary Executed (PATH Issue)

**Observed Behavior:**\
Old shell script version of `mob-consensus` was executed instead of the
newly built Go binary due to PATH precedence.

**Impact:**\
- Debugging confusion - False test results - Demo instability

**Category:** Developer Experience\
**Severity:** Low

------------------------------------------------------------------------

## 3. Branch Creation Logic Inconsistent

**Observed Behavior:**\
- `git switch -c` failed when branch already existed - Help instructions
assumed clean state - Flow not idempotent

**Impact:**\
- Setup process fragile - Onboarding brittle

**Category:** Flow Logic\
**Severity:** High

------------------------------------------------------------------------

## 4. Merge Command Construction Failure

**Observed Behavior:**\
Running:

    mob consensus <branch>

Caused merge failure due to incorrectly constructed git command in Go
implementation.

**Impact:**\
- Core diff + merge flow broken - Tool unusable for intended sync
operation

**Category:** Core Logic Bug\
**Severity:** Critical

------------------------------------------------------------------------

## 5. Local Branch Not Auto-Created Before Merge

**Observed Behavior:**\
Merge failed because corresponding remote branch had not been created
locally.

Manual steps required: - `git switch -c origin/<branch>` - Switch back -
Retry merge

**Expected Behavior:**\
Tool should auto-create local tracking branch if missing.

**Impact:**\
- Automation promise broken - Requires manual Git intervention

**Category:** Automation Gap\
**Severity:** Critical

------------------------------------------------------------------------

## 6. Merge State Not Cleaned Up Automatically

**Observed Behavior:**\
Repository left in "still merging" state after tool operations.

Manual `git commit` required to finalize merge.

**Impact:**\
- Repository left dirty - Confusing state - Breaks smooth UX

**Category:** Workflow Bug\
**Severity:** High

------------------------------------------------------------------------

## 7. `.lock` File Not Ignored

**Observed Behavior:**\
`.lock` file appeared in diffs and ahead/behind calculations.

**Impact:**\
- False divergence signals - Polluted diff output - Incorrect status
reporting

**Recommended Fix:**\
Add `.lock` to `.gitignore`.

**Category:** Repository Hygiene\
**Severity:** High

------------------------------------------------------------------------

## 8. Diff Tool Invocation Unreliable

**Expected Behavior:**\
Tool should: 1. Start merge with `--no-commit --no-ff` 2. Open vimdiff
3. Allow edits 4. Auto-commit

**Observed Behavior:**\
Manual Git commands were required to complete this flow.

**Impact:**\
- Core UX not functioning - Demo required fallback to manual Git

**Category:** Core UX Bug\
**Severity:** High

------------------------------------------------------------------------

## 9. Automatic Push Not Reliable

**Observed Behavior:**\
Tool sometimes required manual `git push`.

Unclear whether push logic is consistently triggered.

**Impact:**\
- Inconsistent sync state - Confusion about canonical repo state

**Category:** Workflow Logic\
**Severity:** Medium

------------------------------------------------------------------------

## 10. Ahead/Behind Status Polluted by Non-Content Changes

**Observed Behavior:**\
Merge commits and lock files influenced ahead/behind reporting.

**Impact:**\
- Misleading collaboration signals - Noisy output

**Category:** Status Logic\
**Severity:** Medium

------------------------------------------------------------------------

## 11. Help Text Not Context-Aware

**Observed Behavior:**\
Help examples hard-coded branch names instead of dynamically using: -
Current branch - Current twig - Current user

**Impact:**\
- Documentation drift - Confusing onboarding

**Category:** UX / Documentation\
**Severity:** Low

------------------------------------------------------------------------

# Severity Summary

### Critical

-   Merge command construction failure
-   Local branch auto-creation missing
-   Username extraction incorrect

### High

-   Merge cleanup incomplete
-   Diff flow unreliable
-   Lock file not ignored
-   Branch creation flow fragile

### Medium

-   Push logic ambiguity
-   Ahead/behind noise

### Low

-   Help text not dynamic
-   PATH confusion

------------------------------------------------------------------------

# Estimated Effort (Rough Order)

  Category                        Estimated Hours
  ------------------------------- -----------------
  Core merge logic fixes          6--10
  Branch automation fixes         4--6
  Username config fix             1--2
  Diff + merge UX stabilization   4--6
  Status reporting cleanup        3--5
  Git hygiene fixes               1--2
  Help text improvements          2--3

**Total Potential Work:** 20--35 hours

------------------------------------------------------------------------

End of document.
