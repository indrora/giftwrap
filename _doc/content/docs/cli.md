---
title: "giftwrap(1)"
draft: false
---

# Synopsis

**giftwrap** — cross-compile Go applications and package releases.

Giftwrap is a tool to build Go applications for
multiple operating systems and architectures at a time.

Additionally, it packages releases for you.

# Usage

    giftwrap [--wrapfile path] [--log-level level] <command> [options]

Available commands:

* **`build`**: Build the project
* **`clean`**: Clean build and release artifacts
* **`init`**: Initialize a project
* **`release`**: Build and package project variants into release archives

# Global Options

    -h, --help               Print help
        --log-level string   Log level to use (debug, info, warn, error) (default: info)
        --wrapfile string    Path to the .wrapfile in use

# Commands

## build — Build the project {#cmd_build}

build compiles each configured target variant for the specified GOOS/GOARCH pairs. If no target name is given, the default target is used.

**Options:**

        --shell string   Specify the shell to use for building

---

## clean — Clean build and release artifacts {#cmd_clean}

clean removes the build directory (_build) and the distribution directory (_dist).

---

## init — Initialize a project {#cmd_init}

Initialize a giftwrap project. This will attempt to find
a go.mod in the current file. If this does not exist, it will stop.

**Options:**

        --modpath string   Path to go.mod file (default: go.mod)

---

## release — Build and package project variants into release archives {#cmd_release}

release builds each configured target variant and packages the output into archives in the distribution directory.

**Options:**

        --shell string   Specify the shell to use for building


# Remarks

## Working Directory & caveats therein

At startup, giftwrap automatically changes directory to the location of the wrapfile.
All paths must be below the wrapfile directory, with the exception of any paths called by the `exec` block.
This may be confusing for some folks. If you need to move files into the build directory from below/outside
the directory the wrapfile is stored in, the enviornment variable `BUILD_DIR` is provided so that you can
add whatever you need postfacto. As an example:

```yaml
targets:
  default:
    exec:
      post: cp -r ../some_path $BUILD_DIR/
```

is valid to perform. All shell actions are performed from within the working directory of the path the
wrapfile is located.

## Compiling for the Raspberry Pi Zero

If you're compiling for the Raspberry Pi Zero -- not the Zero 2 -- you may find that your executables
for arm don't work right. In order to properly target the Pi Zero (W), add the following to your project:

```yaml
targets:
  pizero:
    platforms:
      - linux/arm
    env:
      GOARM: "5"
```

you should now be able to build your package for ARM5/Pi Zero. This is not required on the Pi Zero 2 or any
other ARM64-based device. For more information on required build configurations on ARM, refer to the [Go Wiki](https://go.dev/wiki/GoArm)
for more informaion.

## wrapfile file search order

giftwrap looks for a configuration file in the following locations, in order
(first match wins):

    .wrapfile
    giftwrap.yml
    .github/giftwrap.yml
    .github/.wrapfile
    .giftwrap.yml

Pass `--wrapfile <path>` to override the search and use a specific file.

## Environment variables

giftwrap sets the following environment variables when running pre/post exec hooks
and `go build`:

    GOOS        Target operating system
    GOARCH      Target architecture
    BUILD_TARGET   Name of the target being built (e.g. "default")
    BUILD_DIR      Root of the build output directory

# License

Giftwrap is licensed under the terms of the MIT license.
