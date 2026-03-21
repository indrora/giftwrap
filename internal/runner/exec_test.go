package runner

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/charmbracelet/log"
)

func newTestOptions() Options {
	outBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	return NewOptions().WithStdout(outBuf).WithStderr(errBuf)
}

// TestExecRunner_Success verifies that output from a successful command is
// silently discarded and nothing is written to options.Stdout/Stderr.
func TestExecRunner_Success(t *testing.T) {
	runner := newExecRunnerForTest(t)
	opts := newTestOptions()

	// "go version" always succeeds and prints output.
	if err := runner.RunArgs("go", []string{"version"}, opts); err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	out := opts.Stdout.(*bytes.Buffer).String()
	errOut := opts.Stderr.(*bytes.Buffer).String()

	if out == "" {
		t.Errorf("expected text on stdout, got nothing.")
	}
	if errOut != "" {
		t.Errorf("expected no stderr on success, got %q", errOut)
	}
}

// TestExecRunner_Failure verifies that output is flushed to options.Stdout/Stderr
// when a command exits non-zero.
func TestExecRunner_Failure(t *testing.T) {
	runner := newExecRunnerForTest(t)
	opts := newTestOptions()

	// "go build" on a nonexistent package exits non-zero and writes to stderr.
	err := runner.RunArgs("go", []string{"build", "nonexistent/package/giftwrap_test"}, opts)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var pfe ProcessFailedError
	if !errors.As(err, &pfe) {
		t.Fatalf("expected ProcessFailedError, got %T: %v", err, err)
	}
	if pfe.Code == 0 {
		t.Errorf("expected non-zero exit code, got 0")
	}

	// At least stderr should have been flushed.
	errOut := opts.Stderr.(*bytes.Buffer).String()
	if errOut == "" {
		t.Error("expected stderr output to be flushed on failure, got empty string")
	}
}

// TestExecRunner_StartFailure verifies that a non-existent binary returns a
// ProcessFailedError with exit code -1.
func TestExecRunner_StartFailure(t *testing.T) {
	runner := newExecRunnerForTest(t)
	opts := newTestOptions()

	err := runner.RunArgs("/definitely/does/not/exist/giftwrap_binary", []string{}, opts)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var pfe ProcessFailedError
	if !errors.As(err, &pfe) {
		t.Fatalf("expected ProcessFailedError, got %T: %v", err, err)
	}
	if pfe.Code != -1 {
		t.Errorf("expected exit code -1, got %d", pfe.Code)
	}
}

func newExecRunnerForTest(t *testing.T) *ExecRunner {
	t.Helper()
	logger := log.New(io.Discard)
	return NewExecRunner(logger)
}
