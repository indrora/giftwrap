package runner

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type ExecRunner struct {
	// :)
}

func (r ExecRunner) Run(cmd string, options Options) error {
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
	process := exec.Command(c, args...)

	// Format the command environment

	env := make([]string, 0, len(options.Env))
	for k, v := range options.Env {

		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	process.Env = env
	process.Stdout = options.Stdout
	process.Stderr = options.Stderr

	if err := process.Start(); err != nil {
		return ProcessFailedError{Cmd: c, Code: -1, Reason: err.Error()}
	}

	if err := process.Wait(); err != nil {
		return ProcessFailedError{Cmd: process.String(), Code: process.ProcessState.ExitCode(), Reason: err.Error()}
	}

	return nil
}
