---
title: "giftwrap"
---

![Giftwrap logo](full.png)
{class="bodylogo" srcset="full.png 1x, full@2x.png 2x"}

**giftwrap** cross-compiles Go applications for multiple OS/architecture targets and packages release archives. You define targets in a YAML config file; giftwrap handles the rest.

```yaml
name: myapp

targets:
  default:
    package: .
    platforms:
      - linux/amd64
      - darwin/arm64
      - windows/amd64
```

```
giftwrap init     # generate a starter config from go.mod
giftwrap build    # compile for all configured targets
giftwrap release  # build and package into archives
```

For most projects, `giftwrap init` and `giftwrap build` is all you need.
[Get started →](/docs/quickstart/)

Source: [github.com/indrora/giftwrap](https://github.com/indrora/giftwrap)
