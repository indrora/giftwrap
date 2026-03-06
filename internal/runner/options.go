package runner

import (
	"io"
	"maps"
	"os"
	"strings"
)

type Options struct {
	Env    map[string]string
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
		// Split the line on the first '='
		parts := strings.SplitN(v, "=", 2)
		sysEnvMap[parts[0]] = parts[1]
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

func (o Options) WithStderr(stderr io.Writer) Options {
	o.Stderr = stderr
	return o
}
