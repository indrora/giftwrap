package runner

import (
	"fmt"
	"os/exec"
	"strings"
)

type ExecRunner struct {
	// :)
}

func (r ExecRunner) Run(cmd string, options Options) error {
	c, a, e := splitCommand(cmd)
	if e != nil {
		return e
	}
	return r.RunArgs(c, a, options)
}

func (r ExecRunner) RunArgs(c string, args []string, options Options) error {
	process := exec.Command(c, args...)

	// Format the command environment

	env := make([]string, 0, len(options.Env))
	for k, v := range options.Env {
		env = append(env, fmt.Sprintf("%s=%s", strings.ToUpper(k), v))
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
