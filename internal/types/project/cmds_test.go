package project

import (
	"runtime"
	"testing"

	"github.com/indrora/giftwrap/internal/runner"
	"github.com/indrora/giftwrap/internal/types"
	"go.yaml.in/yaml/v4"
)

func TestBuildCmdsUnmarshalYAML(t *testing.T) {
	tests := []struct {
		name        string
		yaml        string
		wantPreDef  []string
		wantPostDef []string
		wantPreOS   map[string][]string
		wantPostOS  map[string][]string
		wantErr     bool
	}{
		{
			name:       "simple pre and post",
			yaml:       "pre:\n  - echo hello\npost:\n  - echo bye\n",
			wantPreDef: []string{"echo hello"},
			wantPostDef: []string{"echo bye"},
		},
		{
			name:      "pre only scalar",
			yaml:      "pre: echo hello\n",
			wantPreDef: []string{"echo hello"},
		},
		{
			name: "os-specific pre",
			yaml: "pre:\n  - fallback\npre.windows:\n  - win-cmd\npre.unix:\n  - unix-cmd\n",
			wantPreDef: []string{"fallback"},
			wantPreOS:  map[string][]string{"windows": {"win-cmd"}, "unix": {"unix-cmd"}},
		},
		{
			name: "os-specific post",
			yaml: "post:\n  - cleanup\npost.darwin:\n  - mac-cleanup\n",
			wantPostDef: []string{"cleanup"},
			wantPostOS:  map[string][]string{"darwin": {"mac-cleanup"}},
		},
		{
			name: "arbitrary OS keys (plan9, netbsd)",
			yaml: "pre.plan9:\n  - rc script\npre.netbsd:\n  - nbsh cmd\n",
			wantPreOS: map[string][]string{"plan9": {"rc script"}, "netbsd": {"nbsh cmd"}},
		},
		{
			name:    "unknown key",
			yaml:    "build:\n  - something\n",
			wantErr: true,
		},
		{
			name:    "not a mapping",
			yaml:    "- item\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b BuildCmds
			err := yaml.Unmarshal([]byte(tt.yaml), &b)
			if (err != nil) != tt.wantErr {
				t.Fatalf("UnmarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			checkList(t, "PreExec.Default", []string(b.PreExec.Default), tt.wantPreDef)
			checkList(t, "PostExec.Default", []string(b.PostExec.Default), tt.wantPostDef)

			for os, want := range tt.wantPreOS {
				checkList(t, "PreExec.ByOS["+os+"]", []string(b.PreExec.ByOS[os]), want)
			}
			for os, want := range tt.wantPostOS {
				checkList(t, "PostExec.ByOS["+os+"]", []string(b.PostExec.ByOS[os]), want)
			}
		})
	}
}

func checkList(t *testing.T, label string, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("%s: got %v, want %v", label, got, want)
		return
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("%s[%d]: got %q, want %q", label, i, got[i], want[i])
		}
	}
}

// recordRunner captures which commands were run.
type recordRunner struct {
	ran []string
}

func (r *recordRunner) Run(cmd string, _ runner.Options) error {
	r.ran = append(r.ran, cmd)
	return nil
}

func (r *recordRunner) RunArgs(cmd string, args []string, _ runner.Options) error {
	r.ran = append(r.ran, cmd)
	return nil
}

func TestPlatformCommandListRun(t *testing.T) {
	goos := runtime.GOOS

	// Determine what "unix" means for this test run.
	isUnix := goos != "windows"

	tests := []struct {
		name    string
		list    PlatformCommandList
		wantRan string // expected command that gets run
	}{
		{
			name: "exact OS match wins",
			list: PlatformCommandList{
				Default: types.CommandList{"default-cmd"},
				ByOS:    map[string]types.CommandList{goos: {"exact-cmd"}, "unix": {"unix-cmd"}},
			},
			wantRan: "exact-cmd",
		},
		{
			name: "unix fallback on non-windows",
			list: PlatformCommandList{
				Default: types.CommandList{"default-cmd"},
				ByOS:    map[string]types.CommandList{"unix": {"unix-cmd"}},
			},
			wantRan: func() string {
				if isUnix {
					return "unix-cmd"
				}
				return "default-cmd"
			}(),
		},
		{
			name: "default when no match",
			list: PlatformCommandList{
				Default: types.CommandList{"default-cmd"},
				ByOS:    map[string]types.CommandList{},
			},
			wantRan: "default-cmd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := &recordRunner{}
			opts := runner.NewOptions()
			if err := tt.list.Run(rec, opts); err != nil {
				t.Fatalf("Run() error = %v", err)
			}
			if len(rec.ran) == 0 {
				if tt.wantRan != "" {
					t.Errorf("no commands ran, want %q", tt.wantRan)
				}
				return
			}
			if rec.ran[0] != tt.wantRan {
				t.Errorf("ran %q, want %q", rec.ran[0], tt.wantRan)
			}
		})
	}
}
