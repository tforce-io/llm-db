# Coding Conventions

## General

The principles below are for all languages. They use camelCase for demonstration purposes; adapt casing to the language's style.

### Formats

- Let the project's formatter own the format: run it when the project define an formatter or linter, its check must report no files.
- Keep file content reads top-down: high-level intent first, details descending (newspaper metaphor).
- Keep related things close: declarations near first use, callers near callees.
- Always brace control-flow bodies (if/else, for, while, do, switch, when) even for a single statement; never rely on brace-less bodies.
- Start every code file, including test files, with the license header defined in [license-header.md](license-header.md).
<!-- Project-specific / Formats -->

### Names

- Use one word per concept across the codebase; never one word for two concepts.
- Use intention-revealing names: the name says why it exists and what it does, with no mental mapping from a terse abbreviation to a concept.
- Include units in names where values carry them (e.g. `timeoutMs`, `maxBytes`).
- Function names: start with a verb, and disclose side effects (`save`, `delete`, `publish`, or `mutate`).
- Acronyms: stay consistently cased (e.g. `ID`, `URL`, `HTTP`, not `Id`, `Url`, `Http`).
- Booleans: name as predicates (`isReady`, `hasItems`, `canRetry`).
- Avoid generic names: (`data`, `manager`, `helper`, `util`, `info`, `result`) when a domain name exists.
- Avoid noise: names that differ only by noise words or loose abbreviations.
<!-- Project-specific / Names -->

### Packages

<!-- Project-specific / Packages -->

### Types

<!-- Project-specific / Types -->

### Functions

<!-- Project-specific / Functions -->

### Comments

- Document every exported identifier (constants, vars, types, methods, functions) with a doc comment.
- Write comments explaining *why*, not *what*; the code already shows what it does.
- Cut articles, filler, and hedging; keep IDs, numbers, links, and warnings whole.
- Keep comments short and concise: they should be one to three lines; longer explanations belong in a doc file the comment points to, or an extracted function.
- Keep comments up to date with the code: a stale comment is worse than no comment.
- Keep comments about: non-obvious constraints, legal requirements, algorithmic tradeoffs, external system quirks, warnings of consequences, and TODOs with context and ownership.
- Delete comment noise: redundant restatements, change journals, attributions, banner or position markers, closing-brace labels, and commented-out code.
- Never delete or reword comments a tool reads (shebangs, encoding lines, license and SPDX headers, build and compiler directives, suppressions, doctests, generated-file headers). An unexplained suppression is a finding, not a comment to delete.
<!-- Project-specific / Comments -->

### Errors

- Handle errors explicitly where a meaningful decision can be made.
- Wrap errors with original cause and context when propagating.
- Error strings should be lowercase and not end with punctuation (they're often composed into larger messages).
- Make retry, fallback, timeout, and cancellation behavior explicit.
- Never log secrets, tokens, and sensitive data in error messages and logs.
- Never swallow errors or catch broadly without handling, wrapping, or propagating with context.
<!-- Project-specific / Errors -->

## Go

### Formats

- Keep imports sorted and grouped: standard library, then third-party, then local packages, separated by a blank line.
<!-- Project-specific / Go-Formats -->

### Names

- Case marks visibility: MixedCaps (exported) or mixedCaps (unexported); never underscores.
- Package names: short, lowercase, single word, no underscores or mixedCaps; never `util`, `common`, or `base`.
- Interface names: prefer the `-er` suffix for single-method interfaces (e.g. `Reader`, `Closer`); never an `I` prefix.
- Avoid stuttering: don't repeat the package name in exported identifiers (e.g. `log.Logger`, not `log.LogLogger`).
- Receiver names: short (one or two letters) and identical across every method of the type.
- Getters: name after the value returned, without a `Get` prefix (e.g. `Owner`, not `GetOwner`).
- Never write a vague name (`data`, `info`, `obj`, `ret`), a bare-verb function with no object (`handle`, `process`), or a numbered or versioned name (`user2`, `dataV2`).
- Never suffix a type with a noise word (`Manager`, `Helper`, `Util`) or name it after a verb (`ProcessOrder`).
<!-- Project-specific / Go-Names -->

### Packages

- Default to unexported; export only the names another package must call.
- Keep one package per directory; never split a single package's files across unrelated directories.
- Never add to a `utils` or `common` package; name the concept the code implements.
- Keep `init()` to registration only; never reach out to a database, file, or network from one.
<!-- Project-specific / Go-Packages -->

### Types

- Accept an interface as a parameter, return a concrete struct; declare the interface at the consumer that needs it, never beside the implementation.
- Guard a type whose zero value is unsafe to use behind a constructor (`NewX`); document the zero value when it is safe.
<!-- Project-specific / Go-Types -->

### Functions

- Keep exported functions to a handful of parameters; past three related values, group them into a struct.
- Never switch behavior on a `bool` parameter; write two functions or a small typed enum instead.
- Use named returns only to document a short function's intent, never to smuggle a late, hidden mutation.
<!-- Project-specific / Go-Functions -->

### Comments

- Refer to [Go doc conventions](https://go.dev/doc/comment) for specific convention.
- State a package's purpose once, in its doc comment; do not restate it per file.
<!-- Project-specific / Go-Comments -->

### Errors

- Wrap with context and keep the cause: `fmt.Errorf("loading order %s: %w", id, err)`.
- Compare with `errors.Is` and extract with `errors.As`; never compare an error's `.Error()` string.
- Return `(T, error)` rather than a sentinel zero value standing in for failure.
- Never panic for an expected failure in library code; reserve `panic` for a programmer error the caller cannot act on.
- Check every returned error where it is returned; assigning one to `_` outside a test is a defect, not a shortcut.
<!-- Project-specific / Go-Errors -->

### Tests

- Write table-driven tests with `t.Run(tt.name, func(t *testing.T) {...})`; name each case for behavior, not the input shape.
- Call `t.Helper()` in every test helper so a failure reports the caller's line.
- Run with `-race` whenever a test touches a goroutine or shared state; a flake there is a bug to fix, not noise to rerun.
- Reach for `t.Parallel()` only once a test's fixtures are verified independent of its siblings'.
<!-- Project-specific / Go-Tests -->

### Concurrency

- Give every goroutine an owner for its shutdown; never start one with `go func() { ... }()` and no way to wait for or cancel it.
- Take `context.Context` as the first parameter and check `ctx.Done()`/`ctx.Err()` on any path that can block; never store it on a struct field.
- Use `errgroup.Group` for goroutines that share a lifetime and must fail together, rather than a `sync.WaitGroup` plus a hand-rolled error slice.
- The sender closes a channel, never the receiver; closing one with a live sender still writing to it panics.
- Keep package-level state immutable; mutable package state read or written from more than one goroutine is a defect.
- Since Go 1.22 each `for` iteration owns its own loop variables; add a `v := v` capture only when the module still targets an older Go version.
<!-- Project-specific / Go-Concurrency -->

### Smells

- A `util` or `common` package that keeps absorbing unrelated helpers.
- Package-level mutable state read or written from more than one goroutine.
- An `init()` that reaches out to a database, file, or network instead of only registering something.
- A `context.Context` stored on a struct field instead of threaded through as a parameter.
- An error compared by its formatted string instead of `errors.Is`/`errors.As`.
- A wide interface declared next to its one implementation, instead of the small one its consumer needs.
<!-- Project-specific / Go-Smells -->

## Project-specific

<!-- Project-specific -->
