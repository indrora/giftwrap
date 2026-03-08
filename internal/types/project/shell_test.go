package project

import (
	"testing"
)

func TestShellConfigForOS(t *testing.T) {
	tests := []struct {
		name   string
		config ShellConfig
		goos   string
		want   string
	}{
		{
			name:   "empty config windows default",
			config: ShellConfig{},
			goos:   "windows",
			want:   "cmd /c",
		},
		{
			name:   "empty config unix default",
			config: ShellConfig{},
			goos:   "linux",
			want:   "sh -c",
		},
		{
			name:   "empty config darwin default",
			config: ShellConfig{},
			goos:   "darwin",
			want:   "sh -c",
		},
		{
			name:   "empty config plan9 default",
			config: ShellConfig{},
			goos:   "plan9",
			want:   "sh -c",
		},
		{
			name:   "exact windows match",
			config: ShellConfig{"windows": "powershell -Command"},
			goos:   "windows",
			want:   "powershell -Command",
		},
		{
			name:   "exact linux match",
			config: ShellConfig{"linux": "bash -c", "unix": "sh -c"},
			goos:   "linux",
			want:   "bash -c",
		},
		{
			name:   "unix fallback for darwin",
			config: ShellConfig{"unix": "bash -c"},
			goos:   "darwin",
			want:   "bash -c",
		},
		{
			name:   "unix fallback for freebsd",
			config: ShellConfig{"unix": "bash -c"},
			goos:   "freebsd",
			want:   "bash -c",
		},
		{
			name:   "unix key does not apply on windows",
			config: ShellConfig{"unix": "bash -c"},
			goos:   "windows",
			want:   "cmd /c",
		},
		{
			name:   "darwin overrides unix",
			config: ShellConfig{"unix": "bash -c", "darwin": "zsh -c"},
			goos:   "darwin",
			want:   "zsh -c",
		},
		{
			name:   "arbitrary OS (netbsd) exact match",
			config: ShellConfig{"netbsd": "ksh -c"},
			goos:   "netbsd",
			want:   "ksh -c",
		},
		{
			name:   "arbitrary OS (plan9) exact match",
			config: ShellConfig{"plan9": "rc -c"},
			goos:   "plan9",
			want:   "rc -c",
		},
		{
			name:   "arbitrary OS falls back to unix",
			config: ShellConfig{"unix": "bash -c"},
			goos:   "netbsd",
			want:   "bash -c",
		},
		{
			name:   "arbitrary OS falls back to built-in default",
			config: ShellConfig{},
			goos:   "netbsd",
			want:   "sh -c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.ForOS(tt.goos)
			if got != tt.want {
				t.Errorf("ForOS(%q) = %q, want %q", tt.goos, got, tt.want)
			}
		})
	}
}
