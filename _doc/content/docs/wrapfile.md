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

## General usage

### `name`

**Required.** The project name, used to form output filenames and directory names.

```yaml
name: myapp
```

---

### `buildPath`

**Optional. Default: `_build`.** Directory where compiled binaries are written. giftwrap creates a subdirectory per variant: `<buildPath>/<os>-<arch>/`.

```yaml
buildPath: _output/build
```

---

### `distPath`

**Optional. Default: `_dist`.** Directory for release archives produced by `giftwrap release`.

```yaml
distPath: _output/dist
```

---

### `defaultTarget`

**Optional. Default: `default`.** Name of the target to build when none is given on the command line. Must match a key in `targets`.

```yaml
defaultTarget: release
```

---

### `env`

**Optional.** Environment variables applied to all pre/post commands and `go build` invocations. Target-level `env` overrides project-level entries for the same key.

```yaml
env:
  VERSION: "1.2.3"
  CGO_ENABLED: "0"
```

---

### `targets`

**Required.** A map of named build configurations. Each key is used on the command line: `giftwrap build <name>`.

```yaml
targets:
  default:
    package: .
    targets:
      - linux/amd64
      - darwin/arm64
      - windows/amd64
```

---

### `targets.<name>.package`

**Required.** The Go package path passed to `go build`. Use `.` for the module root.

```yaml
targets:
  default:
    package: .
    targets: linux/amd64
```

```yaml
targets:
  mycmd:
    package: ./cmd/mycmd
    targets: linux/amd64
```

---

### `targets.<name>.targets`

**Required.** One or more `GOOS/GOARCH` pairs to build. Any pair listed by `go tool dist list` is valid. Accepts a single string or a list of strings.

```yaml
targets:
  default:
    package: .
    targets:
      - linux/amd64
      - linux/arm64
      - darwin/amd64
      - darwin/arm64
      - windows/amd64
```

---

### `targets.<name>.env`

**Optional.** Environment variables for this target. Merged with project-level `env`; target values win on collision.

```yaml
targets:
  release:
    package: .
    targets: linux/amd64
    env:
      GOFLAGS: "-trimpath"
```

---

### `targets.<name>.flags`

**Optional.** Arbitrary flags appended to the `go build` command line.

```yaml
targets:
  release:
    package: .
    targets: linux/amd64
    flags: "-ldflags=-s -w"
```

---

### `targets.<name>.tags`

**Optional.** Build tags passed as `-tags` to `go build`. Accepts a single string or a list of strings.

```yaml
targets:
  with-sqlite:
    package: .
    targets: linux/amd64
    tags:
      - sqlite
      - netgo
```

---

### `targets.<name>.cgo`

**Optional. Default: `false`.** Whether to enable CGo. When `false`, giftwrap sets `CGO_ENABLED=0` automatically.

```yaml
targets:
  native:
    package: .
    targets: linux/amd64
    cgo: true
```

---

## Shell actions

Pre/post hooks run shell commands before and after each build. The OS qualifier on `exec.pre.<os>` and `exec.post.<os>` refers to the **host OS** — the machine running giftwrap — not the `GOOS` of the build target.

The following variables are available in all hook commands:

| Variable | Value |
|----------|-------|
| `GOOS` | Target OS being built (e.g. `linux`) |
| `GOARCH` | Target architecture being built (e.g. `amd64`) |
| `BUILD_PATH` | Output directory for this variant |
| `BUILD_TARGET` | Name of the current target |

---

### `shell`

**Optional.** Configures which shell is used to run hook commands. Keys are `runtime.GOOS` values or the special key `unix` (any non-Windows OS).

Resolution order:
1. Exact OS match (`linux`, `darwin`, `windows`, `plan9`, …)
2. `unix`
3. Built-in default: `sh -c` on non-Windows, `cmd /c` on Windows

```yaml
shell:
  unix: bash -c
  darwin: zsh -c
  windows: powershell -Command
```

---

### `exec`

**Optional.** Hooks that run **once per build**, around the entire set of variants. Use `targets.<name>.exec` for hooks that run once per variant.

---

### `exec.pre` / `exec.post`

Commands to run before or after the build. Accepts a single string or a list of strings.

```yaml
exec:
  pre:
    - go generate ./...
    - echo "starting build"
  post: echo "all variants built"
```

---

### `exec.pre.<os>` / `exec.post.<os>`

Host-OS-specific variant of `exec.pre` or `exec.post`. `<os>` is any `runtime.GOOS` value, or `unix` (any non-Windows host).

Resolution order for each hook: exact OS match → `unix` → unqualified key.

```yaml
exec:
  pre.unix:
    - chmod +x scripts/prepare.sh
    - ./scripts/prepare.sh
  pre.windows: scripts\prepare.bat
  post.darwin: codesign --deep $BUILD_PATH
```

---

### `targets.<name>.exec`

**Optional.** Hooks scoped to a specific target. These run **once per variant** (once per `GOOS/GOARCH` pair), unlike project-level `exec` which runs once per build. Supports the same `pre`, `pre.<os>`, `post`, `post.<os>` keys.

```yaml
targets:
  default:
    package: .
    targets:
      - linux/amd64
      - windows/amd64
    exec:
      pre: echo "building $BUILD_TARGET for $GOOS/$GOARCH"
      post.windows: sign.bat %BUILD_PATH%
```

---

## Release filename templates

By default, `giftwrap release` names archives as:

```
<projectName>-<targetName>-<goos>-<goarch>.<format>
```

When there is exactly one target and it is named `default`, the target name is omitted:

```
<projectName>-<goos>-<goarch>.<format>
```

---

### `nameTemplate`

**Optional.** A Go [`text/template`](https://pkg.go.dev/text/template) string that overrides the archive filename. The format extension is always appended automatically.

Available variables:

| Variable | Value |
|----------|-------|
| `.Name` | Project name |
| `.Version` | Full version string (e.g. `v1.2.3` or `v1.2.3+4.gabc1234`) |
| `.VersionRC` | Release-candidate form; same as `.Version` on an exact tag |
| `.Major` | Major version number |
| `.Minor` | Minor version number |
| `.Patch` | Patch version number |
| `.Commits` | Commits ahead of the last semver tag (0 = exact tag) |
| `.OS` | Target GOOS |
| `.Arch` | Target GOARCH |
| `.Target` | Target name from the wrapfile |
| `.Format` | Archive format extension (e.g. `tar.gz`) |

```yaml
nameTemplate: "{{.Name}}-{{.Version}}-{{.OS}}-{{.Arch}}"
```

Produces: `myapp-v1.2.3-linux-amd64.tar.gz`

---

### `archiveFormat`

**Optional. Default: `zip` on Windows, `tar.gz` on all other platforms.** Format for release archives. Valid values: `tar.gz`, `tar.zst`, `zip`.

Scalar form applies to all platforms:

```yaml
archiveFormat: tar.gz
```

Map form sets per-platform formats. The `default` key applies to platforms not otherwise listed:

```yaml
archiveFormat:
  default: tar.gz
  windows: zip
  linux: tar.zst
```

---

### `targets.<name>.archiveFormat`

**Optional.** Overrides the project-level `archiveFormat` for a specific target. Same scalar or map form.

---

## Additional files in release

Use `include` to copy extra files into each variant's output directory before it is packaged. Project-level `include` entries are applied first; target-level entries are appended after.

---

### `include`

**Optional.** Files or glob patterns to copy into every variant's output directory.

A plain string copies matching files preserving their relative paths:

```yaml
include:
  - README.md
  - LICENSE
  - docs/**
```

---

### `include[*].src` / `include[*].dest`

To remap files to a different destination, use the mapping form. The non-wildcard prefix of `src` is stripped before prepending `dest` — similar to rsync's trailing-slash behavior.

```yaml
include:
  - src: assets/**     # assets/logo.png → res/logo.png
    dest: res/
  - src: "*.md"        # README.md → docs/README.md
    dest: docs/
```

Glob patterns follow [doublestar](https://github.com/bmatcuk/doublestar) syntax: `**`, `*`, `?`, `[range]`.

---

### `targets.<name>.include`

**Optional.** Additional files for a specific target's variants. Appended after any project-level `include` entries.

```yaml
targets:
  default:
    package: .
    targets: linux/amd64
    include:
      - src: scripts/run.sh
        dest: bin/
```
