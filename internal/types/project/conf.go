package project

import (
	"errors"
	"maps"

	"github.com/indrora/giftwrap/internal/types"
)

// Config holds the configuration for the giftwrap tool.

type Project struct {
	// name of the project (slugified for release files)
	Name string `yaml:"name"`
	// Additional files to copy into each variant's output directory after build.
	// Supports doublestar glob patterns. See types.FileSpec for YAML forms.
	AdditionalFiles []types.FileSpec `yaml:"include,omitempty"`
	// Environment variables to set during pre- and post-build commands.
	Environment map[string]string `yaml:"env,omitempty"`
	// Commands to run pre- and post-build
	Exec *BuildCmds `yaml:"exec,omitempty"`
	// Directory to place build artifacts
	BuildDir      string `yaml:"buildPath,omitempty"`     // defaults to "build"
	DistDir       string `yaml:"distPath,omitempty"`      // defaults to "dist"
	DefaultTarget string `yaml:"defaultTarget,omitempty"` // defaults to "default"
	// Build configurations. Must have at least one.
	Targets map[string]BuildConfig `yaml:"targets"`
}

var (
	TargetNotFoundErr = errors.New("target not found")
	NoPackageErr      = errors.New("package was not specified")
	NoTargetsErr      = errors.New("no targets were specified")
)

func (p *Project) ReifyConfig(target string) (*BuildConfig, error) {

	buildconfig := &BuildConfig{}
	tgt, ok := p.Targets[target]
	if !ok {
		return nil, TargetNotFoundErr
	}

	// convenience, since we're overwriting some things.
	*buildconfig = tgt

	// The environment variables start as the project's enviornment variables, then the target,
	// with target overriding project

	envs := make(map[string]string)
	maps.Insert(envs, maps.All(p.Environment))
	maps.Insert(envs, maps.All(tgt.Environment))

	buildconfig.Environment = envs

	// Additional files: project-level entries come first, target entries append after.
	buildconfig.AdditionalFiles = append(p.AdditionalFiles, tgt.AdditionalFiles...)

	return buildconfig, nil
}
