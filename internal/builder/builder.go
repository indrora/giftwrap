package builder

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/indrora/giftwrap/internal"
	"github.com/indrora/giftwrap/internal/runner"
	"github.com/indrora/giftwrap/internal/types/project"
)

// Builder is the interface that generates commands to run and their environment.

type Builder struct {
	proj      project.Project
	outWriter io.Writer
	errWriter io.Writer
	runner    runner.Runner
	runOpts   runner.Options

	realTargets map[string]project.BuildConfig
}

func NewBuilder(p project.Project, r runner.Runner) (*Builder, error) {
	b := &Builder{}
	b.proj = p
	b.runner = r
	b.runOpts = runner.NewOptions()

	// reify all configurations

	b.realTargets = make(map[string]project.BuildConfig)

	for k := range p.Targets {
		realConfig, err := p.ReifyConfig(k)
		if err != nil {
			return nil, err
		}
		b.realTargets[k] = *realConfig
	}

	return b, nil
}

func (b *Builder) SetIO(out, err io.Writer) {
	b.errWriter = err
	b.outWriter = out
	b.runOpts.Stdout = out
	b.runOpts.Stderr = err
}

func (b *Builder) Setup() error {

	// Create any paths that we have to make

	if err := os.MkdirAll(b.proj.BuildDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create path %s: %v", b.proj.BuildDir, err)
	}

	if err := os.MkdirAll(b.proj.DistDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create path %s: %v", b.proj.DistDir, err)
	}

	if err := b.proj.Exec.PreExec.Run(b.runner, b.runOpts); err != nil {
		return fmt.Errorf("failed to run pre-build functions: %v", err)
	}

	return nil
}

func (b *Builder) BuildTarget(target string) error {

	config, ok := b.realTargets[target]

	if !ok {
		return errors.New("no such target")
	}

	// Run pre-exec calls for func

	for _, variantName := range config.Targets {
		fmt.Printf("Building target %s:%s\n", target, variantName)

		config.Exec.PreExec.Run(b.runner, b.runOpts)

		// TODO: Implement actual build logic...

		buildpath := path.Join(b.proj.BuildDir, internal.Slugify(variantName))
		os.MkdirAll(buildpath, os.ModePerm)

		varsplit := strings.SplitN(variantName, "/", 2)

		opts := b.runOpts.WithSysEnv().WithEnv(map[string]string{
			"GOOS":   varsplit[0],
			"GOARCH": varsplit[1],
		})
		b.runner.RunArgs("go", []string{"build", "-o", buildpath, config.Package}, opts)

		fmt.Printf("Building to path %s\n", buildpath)

		config.Exec.PostExec.Run(b.runner, b.runOpts)

	}

	return nil
}

func (b *Builder) Teardown() error {
	// Call post-build functions

	if err := b.proj.Exec.PostExec.Run(b.runner, b.runOpts); err != nil {
		return err
	}
	return nil
}
