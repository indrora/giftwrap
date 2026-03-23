---
title: "giftwrap(1)"
draft: false
---

## Synopsis

**giftwrap** — cross-compile Go applications and package releases.

Giftwrap is a tool to build Go applications for
multiple operating systems and architectures at a time.

Additionally, it packages releases for you.

## Usage

    giftwrap [--wrapfile path] [--log-level level] <command> [options]

Available commands:

    build        Build the project
    clean        Clean build and release artifacts
    init         Initialize a project
    release      Build and package project variants into release archives

## Global Options

    -h, --help               Print help
        --log-level string   Log level to use (debug, info, warn, error) (default: info)
        --wrapfile string    Path to the .wrapfile in use

## Commands

### build — Build the project

build compiles each configured target variant for the specified GOOS/GOARCH pairs. If no target name is given, the default target is used.

**Options:**

        --shell string   Specify the shell to use for building

---

### clean — Clean build and release artifacts

clean removes the build directory (_build) and the distribution directory (_dist).

---

### init — Initialize a project

Initialize a giftwrap project. This will attempt to find
a go.mod in the current file. If this does not exist, it will stop.

**Options:**

        --modpath string   Path to go.mod file (default: go.mod)

---

### release — Build and package project variants into release archives

release builds each configured target variant and packages the output into archives in the distribution directory.

**Options:**

        --shell string   Specify the shell to use for building


## Remarks

### Configuration file search order

giftwrap looks for a configuration file in the following locations, in order
(first match wins):

    .wrapfile
    giftwrap.yml
    .github/giftwrap.yml
    .github/.wrapfile
    .giftwrap.yml

Pass `--wrapfile <path>` to override the search and use a specific file.

### Environment variables

giftwrap sets the following environment variables when running pre/post exec hooks
and `go build`:

    GOOS        Target operating system
    GOARCH      Target architecture
    BUILD_TARGET   Name of the target being built (e.g. "default")
    BUILD_DIR      Root of the build output directory
