---
title: "CLI Reference"
draft: false
---

## giftwrap

A tool to build Go applications

### Synopsis

Giftwrap is a tool to build Go applications for
	multiple operating systems and architectures at a time.

	Additionally, it packages releases for you.

```
giftwrap [flags]
```

### Options

```
  -h, --help               help for giftwrap
      --log-level string   Log level to use (debug, info, warn, error) (default "info")
      --wrapfile string    Path to the .wrapfile in use
```


---

## giftwrap build

Build the project

### Synopsis

build compiles each configured target variant for the specified GOOS/GOARCH pairs. If no target name is given, the default target is used.

```
giftwrap build [flags]
```

### Options

```
  -h, --help           help for build
      --shell string   Specify the shell to use for building
```

### Options inherited from parent commands

```
      --log-level string   Log level to use (debug, info, warn, error) (default "info")
      --wrapfile string    Path to the .wrapfile in use
```


---

## giftwrap clean

Clean build and release artifacts

```
giftwrap clean [flags]
```

### Options

```
  -h, --help   help for clean
```

### Options inherited from parent commands

```
      --log-level string   Log level to use (debug, info, warn, error) (default "info")
      --wrapfile string    Path to the .wrapfile in use
```


---

## giftwrap init

Initialize a project

### Synopsis

Initialize a giftwrap project. This will attempt to find
a go.mod in the current file. If this does not exist, it will stop.

```
giftwrap init [flags]
```

### Options

```
  -h, --help             help for init
      --modpath string   Path to go.mod file (default "go.mod")
```

### Options inherited from parent commands

```
      --log-level string   Log level to use (debug, info, warn, error) (default "info")
      --wrapfile string    Path to the .wrapfile in use
```


---

## giftwrap release

Build and package project variants into release archives

### Synopsis

release builds each configured target variant and packages the output into archives in the distribution directory.

```
giftwrap release [flags]
```

### Options

```
  -h, --help           help for release
      --shell string   Specify the shell to use for building
```

### Options inherited from parent commands

```
      --log-level string   Log level to use (debug, info, warn, error) (default "info")
      --wrapfile string    Path to the .wrapfile in use
```


---

