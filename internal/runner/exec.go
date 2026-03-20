package runner

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/log"
)

type ExecRunner struct {
	// :)
	logger *log.Logger
}

func NewExecRunner(l *log.Logger) *ExecRunner {
	return &ExecRunner{logger: l}
}

func (r ExecRunner) Run(cmd string, options Options) error {
	r.logger.Debug("Running command", "cmd", cmd, "options", options)
	shell := options.Shell
	if shell == "" {
		if runtime.GOOS == "windows" {
			shell = "cmd /c"
		} else {
			shell = "sh -c"
		}
	}
	parts := strings.Fields(shell)
	return r.RunArgs(parts[0], append(parts[1:], cmd), options)
}

func (r ExecRunner) RunArgs(c string, args []string, options Options) error {
	r.logger.Debug("Running command with args", "cmd", c, "args", args, "options", options)
	process := exec.Command(c, args...)

	// Format the command environment

	env := make([]string, 0, len(options.Env))
	for k, v := range options.Env {

		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	process.Env = env

	var outBuf, errBuf bytes.Buffer
	process.Stdout = &outBuf
	process.Stderr = &errBuf

	if err := process.Start(); err != nil {
		return ProcessFailedError{Cmd: c, Code: -1, Reason: err.Error()}
	}

	if err := process.Wait(); err != nil {
		io.Copy(options.Stdout, &outBuf)
		io.Copy(options.Stderr, &errBuf)
		return ProcessFailedError{Cmd: process.String(), Code: process.ProcessState.ExitCode(), Reason: err.Error()}
	}

	return nil
}
