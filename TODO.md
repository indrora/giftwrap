# giftwrap TODO

Items are grouped by severity. Fix Critical and High items before shipping.

---

## Critical — crashes today

### Nil pointer panic when `exec` is absent

`Project.Exec` is a `*BuildCmds` pointer. If `exec` is not present in the wrapfile, it is `nil`.
`Setup()` and `Teardown()` dereference it unconditionally:

```go
b.proj.Exec.PreExec.Run(b.runner, b.runOpts)
```

This will panic. `BuildTarget()` has the same problem with `config.Exec`.

---

## High — documented features that don't work, or silently broken output

### `cmd/default.yml` uses wrong key

The embedded `default.yml` template uses `config:` as the top-level key for targets:

```yaml
config:
  default:
    ...
```

But the `Project` struct uses `targets:` as the YAML tag. The template is malformed and would
produce an empty target map if loaded. Every `giftwrap init` generates a broken wrapfile.

### Target-level `exec` errors silently ignored

In `BuildTarget()`, both `config.Exec.PreExec.Run(...)` (line 106) and
`config.Exec.PostExec.Run(...)` (line 122) discard the returned error. A failing pre/post hook
will not stop the build.

This is also the foundation for the stop-on-error / keep-going behavior described below.

### Fields parsed but never used

#### `flags` (`BuildConfig.BuildFlags`)
The `flags` field is unmarshalled from the wrapfile into `BuildConfig.BuildFlags` but is never
read in `builder/builder.go`. The `go build` invocation (`builder.go:110`) is hardcoded:

```go
b.runner.RunArgs("go", []string{"build", "-o", buildpath, config.Package}, opts)
```

`BuildFlags` needs to be appended to this args slice.

#### `tags` (`BuildConfig.BuildTags`)
Same situation as `flags`. `BuildTags` is parsed but never passed to `go build` as `-tags`.

#### `cgo` (`BuildConfig.UseCgo`)
`UseCgo` is parsed but never acted on. The builder never sets `CGO_ENABLED` based on this field.
The docs say `CGO_ENABLED=0` is implied when `cgo: false`, but that doesn't happen.

---

## Medium — works, but inconsistently or incompletely

### Global `exec` ignores `shell` config and system environment

In `Setup()` and `Teardown()`, the global pre/post hooks run with the base `b.runOpts`:

```go
b.proj.Exec.PreExec.Run(b.runner, b.runOpts)
```

`b.runOpts` has no shell set and does not include the system environment (`WithSysEnv()` is never
called for it). Per-variant builds in `BuildTarget()` correctly call `.WithSysEnv().WithShell(b.Shell)`,
but global hooks silently fall back to `sh -c`/`cmd /c` regardless of the `shell` config, and run
without the user's `PATH` or other environment variables.

### Stop-on-error / keep-going behavior

Currently, errors from variant builds are not surfaced and the build charges on regardless.
Two behaviors should be offered:

- Stop all remaining builds when any variant fails (default / safe mode)
- Keep going when a variant fails and report all failures at the end (`--keep-going` flag)

---

## Low — docs, stubs, polish

### `giftwrap build` target selection (readme wrong)

`readme.md` shows:

```
giftwrap build --target my-target
```

There is no `--target` flag. Targets are positional arguments:

```
giftwrap build my-target
giftwrap build target1 target2
```

### `giftwrap init` behavior undocumented

- Writes to the path from `--wrapfile` (defaults to `.wrapfile`).
- Accepts `--modpath` to specify a non-default `go.mod` path (undocumented in readme/doc).
- Reads the module path from `go.mod` to set the default package; fails if `go.mod` is absent.

### `giftwrap release` is a stub

`cmd/release.go` only prints `"release called"`. The `distPath` directory is created during
`Setup()` but is never used. The docs describe release packaging as a feature of the tool.

### Cleanup command

Add a `giftwrap clean` command to remove the `buildPath` and `distPath` directories.
