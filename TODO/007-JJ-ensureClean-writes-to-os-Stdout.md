# 007-JJ: ensureClean writes to os.Stdout directly

**Source:** Code review (new finding)
**Severity:** Medium
**Category:** Code Correctness

## Problem

`ensureClean()` at `main.go:436` writes directly to `os.Stdout`:

```go
fmt.Fprintln(os.Stdout, "you have uncommitted changes")
```

Every other output path in the codebase uses the `stdout` or `stderr` writer
parameters passed through the call chain. This function doesn't accept a writer
parameter at all.

## Impact

- Breaks testability — output goes to real stdout even in tests
- Inconsistent with the rest of the codebase's I/O pattern
- Would cause problems if `run()` is ever used as a library function

## Suggested Fix

Add `stdout io.Writer` parameter to `ensureClean()` and use it instead of
`os.Stdout`. Update callers (`runCreateBranch`, `runDiscovery`, `runMerge`) to
pass their `stdout` through.

## Subtasks

- [x] 007.1 Add `stdout io.Writer` parameter to `ensureClean()`
- [x] 007.2 Replace `os.Stdout` with the parameter
- [x] 007.3 Update all callers to pass `stdout`
