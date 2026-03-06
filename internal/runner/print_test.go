package runner

import (
	"bytes"
	"testing"
)

func TestPrintRunnerRun(t *testing.T) {
	buf := &bytes.Buffer{}
	opts := NewOptions().WithStdout(buf)
	r := PrintRunner{}

	if err := r.Run("echo hello", opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got != "echo hello\n" {
		t.Errorf("output = %q, want %q", got, "echo hello\n")
	}
}

func TestPrintRunnerRunArgs(t *testing.T) {
	tests := []struct {
		cmd  string
		args []string
		want string
	}{
		{"echo", []string{"hello", "world"}, "cmd: echo args: [hello world]\n"},
		{"ls", []string{"-la"}, "cmd: ls args: [-la]\n"},
		{"true", []string{}, "cmd: true args: []\n"},
		{"true", nil, "cmd: true args: []\n"},
	}

	for _, tt := range tests {
		t.Run(tt.cmd, func(t *testing.T) {
			buf := &bytes.Buffer{}
			opts := NewOptions().WithStdout(buf)
			r := PrintRunner{}

			if err := r.RunArgs(tt.cmd, tt.args, opts); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := buf.String(); got != tt.want {
				t.Errorf("output = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPrintRunnerImplementsRunner(t *testing.T) {
	var _ Runner = PrintRunner{}
}
