# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What is giftwrap

`giftwrap` is a CLI tool (built with Cobra) that cross-compiles Go applications for multiple OS/architecture targets and packages releases. It reads a project config file (`.wrapfile` or `giftwrap.yml`) and runs pre/post build hooks, sets environment variables, and invokes `go build` for each target variant.

## Commands

```bash
go build ./...          # Build the project
go test ./...           # Run all tests
go test ./internal/runner/...   # Run tests in a specific package
go run . build          # Run the CLI build command
```

## Architecture

**Entry point**: `main.go` → `cmd.Execute()` (Cobra root command)

**Config file search order** (from `cmd/root.go`):
`.wrapfile`, `giftwrap.yml`, `.github/giftwrap.yml`, `.github/.wrapfile`, `.giftwrap.yml`

**Key data flow**:
1. `cmd/load.go:LoadProject` (Cobra `PreRunE`) reads the wrapfile into `project.Project`
2. `project.ReifyConfig(target)` merges project-level and target-level env/files into a `BuildConfig`
3. `builder.Builder` orchestrates the build: `Setup()` → `BuildTarget(name)` → `Teardown()`
4. The `runner.Runner` interface abstracts command execution — `ExecRunner` actually runs processes, `PrintRunner` only prints (dry-run/testing)

**Package layout**:
- `cmd/` — Cobra commands (`build`, `release`, `init`)
- `internal/types/project/` — `Project`, `BuildConfig`, `BuildCmds` structs and YAML loading
- `internal/types/yml_stringslice.go` — `CommandList` type: unmarshals a YAML scalar or sequence into `[]string` and can `.Run()` commands via a `Runner`
- `internal/runner/` — `Runner` interface, `ExecRunner`, `PrintRunner`, `Options` (env + stdio)
- `internal/builder/` — `Builder` that calls `go build` per target variant (`GOOS/GOARCH`)
- `internal/compiler/` — wraps `go tool dist list -json` to enumerate valid targets
- `internal/` — `Slugify` (for build output paths), `SliceDice` helper

**BuildConfig** fields of note: `Package` (Go package path), `Targets` (`[]string` of `GOOS/GOARCH`), `Exec.PreExec`/`PostExec` (`CommandList`), `Environment` (merged with project env, target overrides project).

Defaults loaded in `LoadProject`: `BuildDir="build"`, `DistDir="dist"`, `DefaultTarget="default"`.

# Standards

All agents must follow these standards:

* Comments should be short -- No more than a single line or two -- unless there is a specific
  reason to explain a non-trivial choice
* Prefer simple over complicated: Don't add Another Layer of Abstraction if you can help it. 
* Follow good practices: Don't write a huge function that does 100 things, write 25 functions that do a mix of things in context. 
* Follow the standard Go style guide, with the exception that clear variable names are better than one-letter variables
* Don't reinvent the wheel; If there is a package to do something, check with the user to see if it's appropriate. 
* Don't assume Linux/Windows/MacOS at any time. Write OS-Agnostic code wherever possible. Go gives you that ability, Don't squander it.
* New structs and new functionality should come with tests wherever reasonable.

# Expected behavior

* **NEVER** assume you are right. 
* **ALWAYS** check to see if the approach is correct
