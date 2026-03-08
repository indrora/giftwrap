package project

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/indrora/giftwrap/internal/runner"
	"github.com/indrora/giftwrap/internal/types"
	"go.yaml.in/yaml/v4"
)

// PlatformCommandList holds a default command list plus optional per-OS overrides.
//
// In a wrapfile, OS variants are expressed with a dot-separated suffix on the
// "pre" or "post" key:
//
//	exec:
//	  pre:                   # default — runs on any OS not matched below
//	    - echo hello
//	  pre.unix:              # any non-Windows OS (linux, darwin, freebsd, …)
//	    - echo $USER
//	  pre.darwin:            # macOS only — overrides pre.unix on darwin
//	    - sw_vers
//	  pre.windows:           # Windows only
//	    - echo %USERNAME%
//	  pre.plan9:             # Plan 9 — any runtime.GOOS value is valid
//	    - echo $user
//
// Resolution order (mirrors ShellConfig):
//  1. Exact runtime.GOOS match (e.g. "darwin", "linux", "netbsd", "plan9")
//  2. "unix" — for any non-Windows OS
//  3. Default (unqualified key)
type PlatformCommandList struct {
	Default types.CommandList
	ByOS    map[string]types.CommandList
}

// Run executes the command list appropriate for the current OS.
func (p *PlatformCommandList) Run(r runner.Runner, opts runner.Options) error {
	goos := runtime.GOOS
	if cmds, ok := p.ByOS[goos]; ok {
		return cmds.Run(r, opts)
	}
	if goos != "windows" {
		if cmds, ok := p.ByOS["unix"]; ok {
			return cmds.Run(r, opts)
		}
	}
	return p.Default.Run(r, opts)
}

// BuildCmds holds pre- and post-build command lists, each with optional
// per-OS variants. See PlatformCommandList for the supported key formats.
//
// Example:
//
//	exec:
//	  pre:
//	    - go generate ./...
//	  pre.windows:
//	    - windows-codegen.bat
//	  post:
//	    - echo done
type BuildCmds struct {
	PreExec  PlatformCommandList
	PostExec PlatformCommandList
}

// UnmarshalYAML populates BuildCmds from a mapping whose keys are "pre",
// "post", "pre.<os>", or "post.<os>" (e.g. "pre.windows", "post.darwin").
func (b *BuildCmds) UnmarshalYAML(node *yaml.Node) error {
	// Unwrap document node when called directly via yaml.Unmarshal.
	if node.Kind == yaml.DocumentNode {
		if len(node.Content) == 0 {
			return nil
		}
		node = node.Content[0]
	}
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("exec: expected mapping, got node kind %v", node.Kind)
	}

	b.PreExec.ByOS = make(map[string]types.CommandList)
	b.PostExec.ByOS = make(map[string]types.CommandList)

	for i := 0; i+1 < len(node.Content); i += 2 {
		key := node.Content[i].Value
		valNode := node.Content[i+1]

		var cmdList types.CommandList
		if err := valNode.Decode(&cmdList); err != nil {
			return fmt.Errorf("exec.%s: %w", key, err)
		}

		parts := strings.SplitN(key, ".", 2)
		prefix := parts[0]
		osSuffix := ""
		if len(parts) == 2 {
			osSuffix = parts[1]
		}

		switch prefix {
		case "pre":
			if osSuffix == "" {
				b.PreExec.Default = cmdList
			} else {
				b.PreExec.ByOS[osSuffix] = cmdList
			}
		case "post":
			if osSuffix == "" {
				b.PostExec.Default = cmdList
			} else {
				b.PostExec.ByOS[osSuffix] = cmdList
			}
		default:
			return fmt.Errorf("exec: unknown key %q (expected \"pre\" or \"post\" with optional OS suffix)", key)
		}
	}

	return nil
}
