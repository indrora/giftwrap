package builder

import (
	"fmt"
	"io"
	"os"
	"path"
	"slices"
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
	Shell     string

	realTargets map[string]project.BuildConfig
}

func NewBuilder(p project.Project, r runner.Runner, o runner.Options) (*Builder, error) {
	b := &Builder{}
	b.proj = p
	b.runner = r
	b.Shell = p.Shell.ForHost()
	b.runOpts = o

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

	opts := b.runOpts.WithSysEnv().
		WithShell(b.Shell).
		WithEnv(b.proj.Environment)

	if err := b.proj.Exec.PreExec.Run(b.runner, opts); err != nil {
		return fmt.Errorf("failed to run pre-build functions: %v", err)
	}

	return nil
}

func (b *Builder) BuildTarget(target, variantName string) error {

	config, ok := b.realTargets[target]

	if !ok {
		return fmt.Errorf("no such target %s", target)
	}

	ok = slices.Contains(config.Targets, variantName)
	if !ok {
		return fmt.Errorf("No such variant %s", variantName)
	}

	buildpath := path.Join(b.proj.BuildDir, internal.Slugify(target), internal.Slugify(variantName))

	varsplit := strings.SplitN(variantName, "/", 2)

	opts := b.runOpts.WithSysEnv().
		WithShell(b.Shell).
		WithEnv(config.Environment).
		WithEnv(map[string]string{
			"GOOS":         varsplit[0],
			"GOARCH":       varsplit[1],
			"BUILD_PATH":   buildpath,
			"BUILD_TARGET": target,
		})

	err := config.Exec.PreExec.Run(b.runner, opts)
	if err != nil {
		return fmt.Errorf("Failed pre-exec stage for %s (%s): %v", target, variantName, err)
	}

	err = os.MkdirAll(buildpath, os.ModePerm)

	if err != nil {
		return fmt.Errorf("failed creating output directory for %s (%s): %v", target, variantName, err)
	}

	flags, err := internal.SplitArgs(*config.BuildFlags)
	if err != nil {
		return fmt.Errorf("failed to parse build flags for %s (%s): %v", target, variantName, err)
	}

	args := []string{
		"build", "-o", buildpath,
	}

	if len(config.BuildTags) > 0 {
		args = append(args, "-tags", strings.Join(config.BuildTags, ","))
	}

	if len(flags) > 0 {
		args = append(args, flags...)
	}

	args = append(args, config.Package)

	err = b.runner.RunArgs("go", args, opts)
	if err != nil {
		return fmt.Errorf("build failed for %s (%s): %v", target, variantName, err)
	}

	if len(config.AdditionalFiles) > 0 {
		if err := copyFiles(config.AdditionalFiles, buildpath, os.DirFS(".")); err != nil {
			return fmt.Errorf("copying files for %s: %w", variantName, err)
		}
	}

	err = config.Exec.PostExec.Run(b.runner, opts)
	if err != nil {
		return fmt.Errorf("Failed post-exec stage for %s (%s): %v", target, variantName, err)
	}

	return nil
}

func (b *Builder) Teardown() error {
	// Call post-build functions

	opts := b.runOpts.WithSysEnv().
		WithShell(b.Shell).
		WithEnv(b.proj.Environment)

	if err := b.proj.Exec.PostExec.Run(b.runner, opts); err != nil {
		return fmt.Errorf("Failed pre-exec stage: %v ", err)
	}
	return nil
}
