# giftwrap Configuration Reference

A wrapfile is a YAML file that defines how giftwrap builds your project. It has two sections: global configuration (top-level keys) and per-target configuration (entries under `targets`).

---

## Global Configuration

These keys apply to the whole project and provide defaults that targets can extend or override.

### `name`

**Type:** string
**Required:** yes

The project name. Used to slug-ify output file and directory names.

```yaml
name: myapp
```

---

### `buildPath`

**Type:** string
**Default:** `build`

Directory where compiled binaries are placed. giftwrap creates subdirectories per variant (`<buildPath>/<os>-<arch>/`).

```yaml
buildPath: _output/build
```

---

### `distPath`

**Type:** string
**Default:** `dist`

Directory for packaged release archives.

```yaml
distPath: _output/dist
```

---

### `defaultTarget`

**Type:** string
**Default:** `default`

Name of the target to build when none is specified on the command line.

```yaml
defaultTarget: release
```

---

### `env`

**Type:** map of string → string

Environment variables set for all pre/post commands and `go build` invocations. Target-level `env` entries override project-level entries of the same name.

```yaml
env:
  VERSION: "1.2.3"
  CGO_ENABLED: "0"
```

---

### `shell`

**Type:** map of OS name → shell invocation string

Configures which shell is used to run pre/post commands on each OS. Keys are `runtime.GOOS` values or the special key `unix`.

**Resolution order:**
1. Exact OS match (`linux`, `darwin`, `windows`, `plan9`, …)
2. `unix` — any non-Windows OS
3. Built-in default: `sh -c` (non-Windows), `cmd /c` (Windows)

```yaml
shell:
  unix: bash -c       # linux, darwin, freebsd, etc.
  darwin: zsh -c      # macOS only — overrides unix
  windows: powershell -Command
  plan9: rc -c
```

**Short form** `shell.<x>` refers to the OS-keyed entries shown above.

---

### `exec`

Pre/post commands to run around the entire build (once, not once per target variant). See [exec format](#exec-format) below.

```yaml
exec:
  pre:
    - go generate ./...
  post: echo "all targets built"
```

---

### `include`

**Type:** list of [FileSpec](#filespec-format)

Files or glob patterns to copy into every variant's output directory after each build. Project-level entries are prepended to any target-level entries.

```yaml
include:
  - README.md
  - LICENSE
  - src: docs/**
    dest: docs/
```

---

### `targets`

**Type:** map of name → [Target configuration](#per-target-configuration)

Defines one or more named build configurations. At least one target is required.

```yaml
targets:
  default:
    package: .
    targets:
      - linux/amd64
      - windows/amd64
```

---

## Per-Target Configuration

Each entry under `targets` is a named build configuration. The name is used on the command line and for output paths.

### `package`

**Type:** string
**Required:** yes

The Go package path passed to `go build`. Use `.` for the module root, or a full import path for a subdirectory.

```yaml
targets:
  mycmd:
    package: ./cmd/mycmd
```

---

### `targets`

**Type:** string or list of strings
**Required:** yes

One or more `GOOS/GOARCH` pairs to build. Accepts any target supported by `go tool dist list`.

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

  plan9-only:
    package: .
    targets: plan9/amd64   # single target can be a scalar
```

---

### `env`

**Type:** map of string → string

Additional environment variables for this target. Merged with the project-level `env`; target values override project values for the same key.

The following variables are always set automatically during each variant build:

| Variable | Value |
|---|---|
| `GOOS` | Target OS (e.g. `linux`) |
| `GOARCH` | Target architecture (e.g. `amd64`) |
| `BUILD_PATH` | Output directory for this variant |
| `BUILD_TARGET` | Name of the current target |

```yaml
targets:
  release:
    package: .
    targets: linux/amd64
    env:
      GOFLAGS: "-trimpath"
```

---

### `flags`

**Type:** string
**Optional**

Arbitrary flags appended to the `go build` command line.

```yaml
targets:
  release:
    package: .
    targets: linux/amd64
    flags: "-ldflags=-s -w"
```

---

### `tags`

**Type:** string or list of strings
**Optional**

Build tags passed to `go build -tags`.

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

### `cgo`

**Type:** boolean
**Default:** `false`

Whether to enable CGo. When `false`, `CGO_ENABLED=0` is implied.

```yaml
targets:
  native:
    package: .
    targets: linux/amd64
    cgo: true
```

---

### `exec`

Pre/post commands scoped to this target. Run once per variant (i.e. per `GOOS/GOARCH` pair). See [exec format](#exec-format) below.

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

### `include`

**Type:** list of [FileSpec](#filespec-format)

Additional files or globs to copy into this target's variant directories. Appended after any project-level `include` entries.

```yaml
targets:
  default:
    package: .
    targets: linux/amd64
    include:
      - src: scripts/run.sh
        dest: bin/
```

---

## exec Format

`exec` maps `pre`/`post` keys (with optional OS suffixes) to [CommandList](#commandlist-format) values.

**Keys:**

| Key | When it runs |
|---|---|
| `pre` | Default pre-build hook (all OSes not matched below) |
| `pre.<x>` | Pre-build hook for a specific OS (e.g. `pre.linux`, `pre.darwin`, `pre.windows`) |
| `post` | Default post-build hook |
| `post.<x>` | Post-build hook for a specific OS |

`<x>` is any `runtime.GOOS` value or the special value `unix` (matches any non-Windows OS).

**Resolution order** (same as `shell`):
1. Exact OS match
2. `unix`
3. Unqualified default

```yaml
exec:
  pre:
    - go generate ./...
    - echo "starting build"
  pre.unix:
    - chmod +x scripts/prepare.sh
    - ./scripts/prepare.sh
  pre.windows:
    - scripts\prepare.bat
  post: echo "done"
  post.darwin: dsymutil $BUILD_PATH
```

---

## FileSpec Format

A `FileSpec` specifies a file or glob pattern to copy into the output directory. It has two forms:

**Plain string** — copies matching files, preserving their relative paths:

```yaml
include:
  - README.md
  - docs/**
  - assets/*.png
```

**Mapping form** — remaps files to a destination subdirectory. The non-wildcard prefix of `src` is stripped before prepending `dest` (rsync-style):

```yaml
include:
  - src: assets/**     # assets/logo.png → res/logo.png
    dest: res/
  - src: "*.md"        # README.md → docs/README.md
    dest: docs/
```

Glob patterns follow [doublestar](https://github.com/bmatcuk/doublestar) syntax: `**`, `*`, `?`, `[range]`.

---

## CommandList Format

Anywhere a command list is accepted, you can use either a single string or a sequence:

```yaml
# Single command
exec:
  pre: go generate ./...

# Multiple commands
exec:
  pre:
    - go generate ./...
    - echo "generated"
```

Each string is passed to the configured shell for the current OS (see [`shell`](#shell)).

---

## Inheritance and Merge Rules

| Field | Merge behavior |
|---|---|
| `env` | Project env merged first; target env overrides on key collision |
| `include` | Project entries prepended; target entries appended |
| `exec` | Not merged — project-level runs once around the whole build; target-level runs per variant |
| `shell` | Global only; targets share the project shell config |
