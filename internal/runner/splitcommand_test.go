package runner

import (
	"testing"
)

func TestSplitCommand(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantExe string
		wantArgs []string
		wantErr bool
	}{
		{
			name:    "empty string",
			input:   "",
			wantExe: "",
			wantArgs: nil,
			wantErr: false,
		},
		{
			name:    "single word",
			input:   "echo",
			wantExe: "echo",
			wantArgs: []string{},
			wantErr: false,
		},
		{
			name:    "simple command with args",
			input:   "echo hello world",
			wantExe: "echo",
			wantArgs: []string{"hello", "world"},
			wantErr: false,
		},
		{
			name:    "double-quoted argument with space",
			input:   `echo "hello world"`,
			wantExe: "echo",
			wantArgs: []string{"hello world"},
			wantErr: false,
		},
		{
			name:    "single-quoted argument with space",
			input:   "echo 'hello world'",
			wantExe: "echo",
			wantArgs: []string{"hello world"},
			wantErr: false,
		},
		{
			name:    "backslash-escaped space",
			input:   `echo hello\ world`,
			wantExe: "echo",
			wantArgs: []string{"hello world"},
			wantErr: false,
		},
		{
			name:    "backslash escape in double quotes",
			input:   `echo "hello\"world"`,
			wantExe: "echo",
			wantArgs: []string{`hello"world`},
			wantErr: false,
		},
		{
			name:    "no escape inside single quotes",
			input:   `echo 'hello\nworld'`,
			wantExe: "echo",
			wantArgs: []string{`hello\nworld`},
			wantErr: false,
		},
		{
			name:    "multiple spaces between args",
			input:   "echo  hello   world",
			wantExe: "echo",
			wantArgs: []string{"hello", "world"},
			wantErr: false,
		},
		{
			name:    "tabs between args",
			input:   "echo\thello\tworld",
			wantExe: "echo",
			wantArgs: []string{"hello", "world"},
			wantErr: false,
		},
		{
			name:    "unclosed single quote",
			input:   "echo 'hello",
			wantErr: true,
		},
		{
			name:    "unclosed double quote",
			input:   `echo "hello`,
			wantErr: true,
		},
		{
			name:    "trailing backslash",
			input:   `echo hello\`,
			wantErr: true,
		},
		{
			name:    "command with path",
			input:   "/usr/bin/env python3 script.py",
			wantExe: "/usr/bin/env",
			wantArgs: []string{"python3", "script.py"},
			wantErr: false,
		},
		{
			name:    "quoted executable",
			input:   `"/path/with spaces/bin" arg`,
			wantExe: "/path/with spaces/bin",
			wantArgs: []string{"arg"},
			wantErr: false,
		},
		{
			name:    "empty quoted string as arg",
			input:   `echo "" foo`,
			wantExe: "echo",
			wantArgs: []string{"", "foo"},
			wantErr: false,
		},
		{
			name:     "backslash-escaped space in executable",
			input:    `foo\ bar baz`,
			wantExe:  "foo bar",
			wantArgs: []string{"baz"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exe, args, err := splitCommand(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("splitCommand(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if exe != tt.wantExe {
				t.Errorf("exe = %q, want %q", exe, tt.wantExe)
			}
			if len(args) != len(tt.wantArgs) {
				t.Errorf("args = %v, want %v", args, tt.wantArgs)
				return
			}
			for i, arg := range args {
				if arg != tt.wantArgs[i] {
					t.Errorf("args[%d] = %q, want %q", i, arg, tt.wantArgs[i])
				}
			}
		})
	}
}
