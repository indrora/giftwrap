# giftwrap

`giftwrap` cross-compiles Go applications for multiple OS/architecture targets and packages releases. It reads a project config file (a *wrapfile*), runs pre/post build hooks, sets environment variables, and invokes `go build` for each target you define.

## Building

```
go build .
```

No working copy of giftwrap required.

## Usage

```
giftwrap init             # generate a starter wrapfile from go.mod
giftwrap build            # build the default target
giftwrap build my-target  # build a named target
giftwrap release          # build and package into release archives
giftwrap clean            # remove build and dist directories
```

giftwrap searches for a config file in this order:
`.wrapfile`, `giftwrap.yml`, `.github/giftwrap.yml`, `.github/.wrapfile`, `.giftwrap.yml`

Pass `--wrapfile <path>` to use a specific file.

## Example wrapfile

```yaml
name: myapp

env:
  VERSION: "1.0.0"

exec:
  pre: go generate ./...
  post: echo "build complete"

include:
  - README.md
  - LICENSE

targets:
  default:
    package: .
    platforms:
      - linux/amd64
      - linux/arm64
      - darwin/amd64
      - darwin/arm64
      - windows/amd64
    env:
      CGO_ENABLED: "0"
    exec:
      post: echo "built $BUILD_TARGET for $GOOS/$GOARCH"

  with-cgo:
    package: .
    platforms: linux/amd64
    cgo: true
    include:
      - src: vendor/libs/**
        dest: lib/
```

## Generating documentation

The docs site lives in `_doc/` and is built with [Hugo](https://gohugo.io). CLI reference pages are generated from the Cobra command tree.

```
go run -tags docgen . docgen   # regenerate _doc/content/docs/cli.md
cd _doc && hugo server         # preview locally
cd _doc && hugo --minify       # production build
```

Run `go run ./tools/docgen` after any change to command flags or descriptions.

## Architecture

**Entry point:** `main.go` → `cmd.Execute()` (Cobra root command)

**Data flow:**

1. `cmd/load.go:LoadProject` (Cobra `PreRunE`) — locates the wrapfile and parses it into a `project.Project`
2. `project.ReifyConfig(targetName)` — merges project-level and target-level env, include, and exec into a resolved `BuildConfig`
3. `builder.Builder` orchestrates the build: `Setup()` (project pre-exec) → per-variant `BuildTarget()` → `Teardown()` (project post-exec)
4. `packager.Packager` (release only) — packages each variant's output directory into a release archive

**Packages:**

| Package | Purpose |
|---------|---------|
| `cmd/` | Cobra commands: `build`, `release`, `init`, `clean` |
| `internal/types/project/` | `Project` and `BuildConfig` structs, YAML loading, validation |
| `internal/types/` | Shared types: `CommandList`, `FileSpec`, `MultiString` |
| `internal/builder/` | Orchestrates `go build` per target variant |
| `internal/runner/` | `Runner` interface; `ExecRunner` (real), `PrintRunner` (dry-run) |
| `internal/packager/` | Packages variant output into release archives |
| `internal/archiver/` | Archive creation (`tar.gz`, `tar.zst`, `zip`) and filename generation |
| `internal/compiler/` | Wraps `go tool dist list` to enumerate valid GOOS/GOARCH pairs |
| `tools/docgen/` | Generates `_doc/content/docs/cli.md` from the Cobra command tree |
