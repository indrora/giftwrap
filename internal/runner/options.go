package runner

import (
	"io"
	"maps"
	"os"
	"strings"
)

type Options struct {
	Env    map[string]string
	Shell  string // shell invocation, e.g. "sh -c" or "cmd /c"
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
	newEnv := make(map[string]string, len(o.Env))
	for _, v := range os.Environ() {
		key, val, _ := strings.Cut(v, "=")
		newEnv[key] = val
	}
	maps.Insert(newEnv, maps.All(o.Env)) // existing entries override system env
	o.Env = newEnv
	return o
}

func (o Options) WithEnv(env map[string]string) Options {
	newEnv := make(map[string]string, len(o.Env)+len(env))
	maps.Insert(newEnv, maps.All(o.Env))
	maps.Insert(newEnv, maps.All(env))
	o.Env = newEnv
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
