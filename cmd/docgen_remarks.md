
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
