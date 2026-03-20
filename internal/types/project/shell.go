package project

import "runtime"

// ShellConfig maps OS names to shell invocation strings used when running
// pre/post commands. The map keys are runtime.GOOS values or the special
// key "unix" as a catch-all for non-Windows systems.
//
// Example wrapfile usage:
//
//	shell:
//	  unix: bash -c      # linux, darwin, freebsd, etc.
//	  darwin: zsh -c     # macOS specifically (overrides unix)
//	  windows: cmd /c
//	  plan9: rc -c
//
// Resolution order for a given OS:
//  1. Exact OS name (e.g. "linux", "darwin", "netbsd", "plan9")
//  2. "unix" — matches any non-Windows OS
//  3. Built-in default: "cmd /c" on Windows, "sh -c" everywhere else
type ShellConfig map[string]string

func (s ShellConfig) ForHost() string {
	return s.ForOS(runtime.GOOS)
}

// ForOS returns the shell invocation string to use for the given GOOS value.
func (s ShellConfig) ForOS(goos string) string {
	if shell, ok := s[goos]; ok {
		return shell
	} else if _UNIX_FALLBACK {
		// Check if we have a shell set for Unix
		if shell, ok := s["unix"]; ok {
			return shell
		}
	}

	return _DEFAULT_SHELL
}
