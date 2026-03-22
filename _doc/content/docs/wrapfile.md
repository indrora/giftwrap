---
title: "Wrapfile Reference"
draft: false
---

A wrapfile is a YAML file that configures how giftwrap builds your project. giftwrap searches for one of these filenames in the current directory:

```
.wrapfile
giftwrap.yml
.github/giftwrap.yml
.github/.wrapfile
.giftwrap.yml
```

You can also specify a path explicitly: `giftwrap --wrapfile path/to/file.yml`.

---

## Top-level configuration

`name`
: **Required.** Project name, used in output filenames and directory names.

`buildPath`
: **Optional. Default: `_build`.** Directory where compiled binaries are written. Subdirectories are created per variant: `<buildPath>/<os>-<arch>/`.

`distPath`
: **Optional. Default: `_dist`.** Directory for release archives produced by `giftwrap release`.

`defaultTarget`
: **Optional. Default: `default`.** Target name to build when none is given on the command line. Must match a key in `targets`.

`env`
: **Optional.** Project-wide environment variables. Applied to all `go build` invocations and hook commands. Target-level `env` overrides project-level on key collision.

`targets`
: **Required.** Map of named build configurations. At least one entry is required.

```yaml
name: myapp
buildPath: _build
distPath: _dist
defaultTarget: default

env:
  VERSION: "1.0.0"

targets:
  default:
    package: .
    platforms:
      - linux/amd64
      - darwin/arm64
      - windows/amd64
```

---

## Target configuration

Each entry under `targets` is a named build configuration. The name is used on the command line (`giftwrap build <name>`) and in output paths.

`targets.<name>.package`
: **Required.** Go package path passed to `go build`. Use `.` for the module root, or a full import path for a subdirectory.

`targets.<name>.platforms`
: **Required.** One or more `GOOS/GOARCH` pairs. Accepts a single string or a list of strings. Any pair listed by `go tool dist list` is valid.

`targets.<name>.env`
: **Optional.** Target-specific environment variables. Merged with project `env`; target values win on collision.

`targets.<name>.flags`
: **Optional.** Arbitrary flags appended to the `go build` command line.

`targets.<name>.tags`
: **Optional.** Build tags, passed as `-tags`. Accepts a single string or a list of strings.

`targets.<name>.cgo`
: **Optional. Default: `false`.** Enables CGo. When `false`, giftwrap sets `CGO_ENABLED=0` automatically.

```yaml
targets:
  default:
    package: .
    platforms:
      - linux/amd64
      - linux/arm64
      - darwin/arm64
      - windows/amd64
    env:
      GOFLAGS: "-trimpath"
    flags: "-ldflags=-s -w"
    tags:
      - netgo
    cgo: false

  with-cgo:
    package: ./cmd/native
    platforms: linux/amd64
    cgo: true
```

---

## Shell actions

Hook commands run in a shell on the **host machine** — the OS running giftwrap. The `.<os>` qualifier on any hook key refers to the host OS, not the `GOOS` build target.

### Execution order

For each `giftwrap build` or `giftwrap release` run:

1. Project `exec.pre[.<host-os>]` runs
2. For each `GOOS/GOARCH` variant:
   1. Target `exec.pre[.<host-os>]` runs
   2. `go build` runs
   3. Target `exec.post[.<host-os>]` runs
3. Project `exec.post[.<host-os>]` runs

The following variables are set automatically and available in all hook commands:

| Variable | Value |
|----------|-------|
| `GOOS` | The target OS being built (e.g. `linux`) |
| `GOARCH` | The target architecture (e.g. `amd64`) |
| `BUILD_PATH` | Output directory for this variant |
| `BUILD_TARGET` | Name of the current target |

### Shell configuration

`shell`
: **Optional.** Map of host OS name → shell invocation string. Resolution order: exact OS match → `unix` → built-in default (`sh -c` on non-Windows, `cmd /c` on Windows).

`shell.<os>`
: A single OS entry. `<os>` is any `runtime.GOOS` value. The special key `unix` matches all non-Windows platforms.

```yaml
shell:
  unix: bash -c
  darwin: zsh -c        # overrides unix on macOS
  windows: powershell -Command
```

### Hooks

`exec.pre` / `exec.post`
: Project-level hooks. Run **once per build**, before or after all variants. Accepts a single string or a list of strings.

`exec.pre.<os>` / `exec.post.<os>`
: OS-qualified project hooks. `<os>` is any `runtime.GOOS` value or `unix`. Resolution: exact OS → `unix` → unqualified key.

`targets.<name>.exec.pre` / `targets.<name>.exec.post`
: Target-level hooks. Run **once per variant** (once per `GOOS/GOARCH` pair).

`targets.<name>.exec.pre.<os>` / `targets.<name>.exec.post.<os>`
: OS-qualified target hooks. Same resolution order as project-level.

```yaml
# Project-level: runs once, wraps the whole build
exec:
  pre:
    - go generate ./...
  pre.unix:
    - chmod +x scripts/sign.sh
  pre.windows: scripts\codegen.bat
  post: echo "build complete"
```

```yaml
# Target-level: runs once per GOOS/GOARCH variant
targets:
  default:
    package: .
    platforms:
      - linux/amd64
      - windows/amd64
    exec:
      pre: echo "building $BUILD_TARGET for $GOOS/$GOARCH"
      post.windows: sign.bat %BUILD_PATH%
      post.darwin: codesign --deep $BUILD_PATH
```

---

## Release archives

`archiveFormat`
: **Optional. Default: `zip` on Windows, `tar.gz` on all other platforms.** Format for archives produced by `giftwrap release`. Valid values: `tar.gz`, `tar.zst`, `zip`. Accepts a scalar string or a map of host OS → format. The `default` key in a map applies to platforms not otherwise listed.

`targets.<name>.archiveFormat`
: **Optional.** Overrides the project-level `archiveFormat` for a specific target.

`nameTemplate`
: **Optional.** Go [`text/template`](https://pkg.go.dev/text/template) string controlling archive filenames. The format extension is always appended automatically. When absent, archives are named `<name>-<target>-<goos>-<goarch>.<format>` (target name omitted when there is exactly one target named `default`).

```yaml
archiveFormat:
  default: tar.gz
  windows: zip
  linux: tar.zst

nameTemplate: "{{.Name}}-{{.Version}}-{{.OS}}-{{.Arch}}"
# produces: myapp-v1.2.3-linux-amd64.tar.gz
```

Variables available in `nameTemplate`:

| Variable | Value |
|----------|-------|
| `.Name` | Project name |
| `.Version` | Full version string (e.g. `v1.2.3` or `v1.2.3+4.gabc1234`) |
| `.VersionRC` | Release-candidate form; same as `.Version` on an exact tag |
| `.Major` / `.Minor` / `.Patch` | Numeric version components |
| `.Commits` | Commits ahead of the last semver tag (0 on an exact tag) |
| `.OS` | Target GOOS |
| `.Arch` | Target GOARCH |
| `.Target` | Target name from the wrapfile |
| `.Format` | Archive format extension (e.g. `tar.gz`) |

---

## Additional files

Use `include` to copy files into each variant's output directory before it is packaged. Project-level entries are applied first; target-level entries are appended after.

`include`
: **Optional.** Project-wide list of file specs. Copied into every variant's output directory.

`targets.<name>.include`
: **Optional.** Target-specific file specs. Appended after project `include` entries.

A plain string copies matching files preserving their relative paths. The mapping form (`src`/`dest`) remaps files to a destination subdirectory — the non-wildcard prefix of `src` is stripped before prepending `dest`. Glob patterns follow [doublestar](https://github.com/bmatcuk/doublestar) syntax.

```yaml
include:
  - README.md           # copied as README.md
  - LICENSE
  - src: docs/**        # docs/guide.md → documentation/guide.md
    dest: documentation/
  - src: "*.md"         # README.md → extra/README.md
    dest: extra/
```
