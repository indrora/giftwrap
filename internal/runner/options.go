package runner

import (
	"io"
	"maps"
	"os"
	"strings"
)

type Options struct {
	Env    map[string]string
	Shell  string // shell invocation, e.g. "sh -c" or "cmd /c"; empty = runtime default
	Stdout io.Writer
	Stderr io.Writer
}

func NewOptions() Options {
	return Options{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    make(map[string]string),
	}
}

func (o Options) WithSysEnv() Options {
	sysEnv := os.Environ()
	sysEnvMap := make(map[string]string)

	for _, v := range sysEnv {
		key, val, _ := strings.Cut(v, "=")
		sysEnvMap[key] = val
	}

	maps.Insert(o.Env, maps.All(sysEnvMap))

	return o
}

func (o Options) WithEnv(env map[string]string) Options {
	maps.Insert(o.Env, maps.All(env))
	return o
}

func (o Options) WithStdout(stdout io.Writer) Options {
	o.Stdout = stdout
	return o
}

func (o Options) WithShell(shell string) Options {
	o.Shell = shell
	return o
}

func (o Options) WithStderr(stderr io.Writer) Options {
	o.Stderr = stderr
	return o
}
