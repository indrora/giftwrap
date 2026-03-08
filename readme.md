# giftwrap

`giftwrap` cross-compiles Go applications for multiple OS/architecture targets and packages releases. It reads a project config file (a *wrapfile*), runs pre/post build hooks, sets environment variables, and invokes `go build` for each target you define.

## Usage

```
giftwrap build            # build the default target
giftwrap build --target my-target
giftwrap init             # generate a starter wrapfile
```

giftwrap searches for a config file in this order:
`.wrapfile`, `giftwrap.yml`, `.github/giftwrap.yml`, `.github/.wrapfile`, `.giftwrap.yml`

You can also pass `--wrapfile <path>` to use a specific file.

## Example wrapfile

```yaml
name: myapp
buildPath: build
distPath: dist

shell:
  unix: bash -c
  windows: cmd /c

env:
  VERSION: "1.0.0"

exec:
  pre:
    - go generate ./...
  pre.windows:
    - windows-codegen.bat
  post: echo "Build complete"

include:
  - README.md
  - LICENSE
  - src: docs/**
    dest: docs/

targets:
  default:
    package: .
    targets:
      - linux/amd64
      - linux/arm64
      - darwin/amd64
      - darwin/arm64
      - windows/amd64
    env:
      CGO_ENABLED: "0"
    exec:
      post: echo "Built $BUILD_TARGET for $GOOS/$GOARCH"

  with-cgo:
    package: .
    targets: linux/amd64
    cgo: true
    include:
      - src: vendor/libs/**
        dest: lib/
```
