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
