package runner

import (
	"bytes"
	"os"
	"testing"
)

func TestNewOptions(t *testing.T) {
	opts := NewOptions()

	if opts.Stdout != os.Stdout {
		t.Errorf("Stdout = %v, want os.Stdout", opts.Stdout)
	}
	if opts.Stderr != os.Stderr {
		t.Errorf("Stderr = %v, want os.Stderr", opts.Stderr)
	}
	if opts.Env == nil {
		t.Error("Env is nil, want initialized map")
	}
	if len(opts.Env) != 0 {
		t.Errorf("Env has %d entries, want 0", len(opts.Env))
	}
}

func TestWithStdout(t *testing.T) {
	buf := &bytes.Buffer{}
	opts := NewOptions().WithStdout(buf)
	if opts.Stdout != buf {
		t.Errorf("Stdout = %v, want buf", opts.Stdout)
	}
	// original unchanged
	orig := NewOptions()
	if orig.Stdout != os.Stdout {
		t.Error("WithStdout mutated original Options")
	}
}

func TestWithStderr(t *testing.T) {
	buf := &bytes.Buffer{}
	opts := NewOptions().WithStderr(buf)
	if opts.Stderr != buf {
		t.Errorf("Stderr = %v, want buf", opts.Stderr)
	}
	orig := NewOptions()
	if orig.Stderr != os.Stderr {
		t.Error("WithStderr mutated original Options")
	}
}

func TestWithEnv(t *testing.T) {
	opts := NewOptions().WithEnv(map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	})

	if opts.Env["FOO"] != "bar" {
		t.Errorf("Env[FOO] = %q, want %q", opts.Env["FOO"], "bar")
	}
	if opts.Env["BAZ"] != "qux" {
		t.Errorf("Env[BAZ] = %q, want %q", opts.Env["BAZ"], "qux")
	}
}

func TestWithEnvMerges(t *testing.T) {
	opts := NewOptions().
		WithEnv(map[string]string{"FOO": "bar"}).
		WithEnv(map[string]string{"BAZ": "qux"})

	if opts.Env["FOO"] != "bar" {
		t.Errorf("Env[FOO] = %q, want %q", opts.Env["FOO"], "bar")
	}
	if opts.Env["BAZ"] != "qux" {
		t.Errorf("Env[BAZ] = %q, want %q", opts.Env["BAZ"], "qux")
	}
}

func TestWithEnvOverwrites(t *testing.T) {
	opts := NewOptions().
		WithEnv(map[string]string{"FOO": "original"}).
		WithEnv(map[string]string{"FOO": "overwritten"})

	if opts.Env["FOO"] != "overwritten" {
		t.Errorf("Env[FOO] = %q, want %q", opts.Env["FOO"], "overwritten")
	}
}

func TestWithSysEnv(t *testing.T) {
	// Set a known env var and verify it appears after WithSysEnv.
	t.Setenv("GIFTWRAP_TEST_VAR", "hello")

	opts := NewOptions().WithSysEnv()

	if opts.Env["GIFTWRAP_TEST_VAR"] != "hello" {
		t.Errorf("Env[GIFTWRAP_TEST_VAR] = %q, want %q", opts.Env["GIFTWRAP_TEST_VAR"], "hello")
	}
}

func TestWithSysEnvDoesNotLoseExisting(t *testing.T) {
	t.Setenv("GIFTWRAP_TEST_VAR", "hello")

	opts := NewOptions().
		WithEnv(map[string]string{"CUSTOM": "value"}).
		WithSysEnv()

	if opts.Env["CUSTOM"] != "value" {
		t.Errorf("Env[CUSTOM] = %q, want %q", opts.Env["CUSTOM"], "value")
	}
	if opts.Env["GIFTWRAP_TEST_VAR"] != "hello" {
		t.Errorf("Env[GIFTWRAP_TEST_VAR] = %q, want %q", opts.Env["GIFTWRAP_TEST_VAR"], "hello")
	}
}
