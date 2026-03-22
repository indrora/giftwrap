---
title: "Quickstart"
draft: false
---

This guide covers the common case: you have a Go module, you want to build binaries for Linux, macOS, and Windows, and you want to package them as release archives.

## Prerequisites

- Go 1.18 or later
- `giftwrap` installed (`go install github.com/indrora/giftwrap@latest`)

## 1. Initialize

Run `giftwrap init` in your project's root directory (where `go.mod` lives):

```
giftwrap init
```

This writes a `.wrapfile` with a `default` target pointing at your module. The generated file looks like this:

```yaml
name: MyProject

targets:
  default:
    package: github.com/example/myapp
    targets:
      - linux/arm64
      - linux/amd64
      - darwin/arm64
      - darwin/amd64
      - windows/arm64
      - windows/amd64
```

Edit `name` to match your project. Remove any `targets` entries for platforms you don't need.

## 2. Build

```
giftwrap build
```

Binaries are placed in `_build/<os>-<arch>/`. For the config above:

```
_build/
  linux-arm64/myapp
  linux-amd64/myapp
  darwin-arm64/myapp
  darwin-amd64/myapp
  windows-arm64/myapp.exe
  windows-amd64/myapp.exe
```

To build a specific named target:

```
giftwrap build my-target
```

## 3. Package a release

```
giftwrap release
```

Packages each variant into an archive under `_dist/`. Default formats: `.zip` for Windows, `.tar.gz` for everything else.

```
_dist/
  myapp-linux-arm64.tar.gz
  myapp-linux-amd64.tar.gz
  myapp-darwin-arm64.tar.gz
  myapp-darwin-amd64.tar.gz
  myapp-windows-arm64.zip
  myapp-windows-amd64.zip
```

## Including extra files

Add `include` to copy files into every variant's output directory before packaging:

```yaml
name: myapp

include:
  - README.md
  - LICENSE

targets:
  default:
    package: .
    targets:
      - linux/amd64
      - darwin/arm64
      - windows/amd64
```

## Pre/post build hooks

Use `exec` to run commands before or after each variant build. Environment variables `GOOS`, `GOARCH`, `BUILD_PATH`, and `BUILD_TARGET` are available:

```yaml
name: myapp

targets:
  default:
    package: .
    targets:
      - linux/amd64
      - darwin/arm64
      - windows/amd64
    exec:
      pre: go generate ./...
      post: echo "built $BUILD_TARGET for $GOOS/$GOARCH"
```

For OS-specific hooks, use `pre.<os>` and `post.<os>`:

```yaml
    exec:
      pre.unix: ./scripts/prepare.sh
      pre.windows: scripts\prepare.bat
```

## Cleaning up

```
giftwrap clean
```

Removes `_build/` and `_dist/`.

## Next steps

- [Wrapfile reference](/docs/wrapfile/) — all available configuration fields
- [CLI reference](/docs/cli/) — all commands and flags
