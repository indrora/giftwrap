package project

import (
	"errors"
	"io"
	"maps"
	"os"
	"regexp"
	"strings"

	"github.com/indrora/giftwrap/internal/types"
	"go.yaml.in/yaml/v4"
)

// Config holds the configuration for the giftwrap tool.

type Project struct {
	// name of the project (slugified for release files)
	Name string `yaml:"name"`
	// Build configurations. Must have at least one.
	Configurations map[string]BuildConfig `yaml:"config"`
	// Additional files to include in the release
	AdditionalFiles types.StringOrSlice `yaml:"include,omitempty"`
	// Environment variables to set during pre- and post-build commands.
	Environment map[string]string `yaml:"env,omitempty"`
	// Commands to run pre- and post-build
	Exec *BuildCmds `yaml:"exec,omitempty"`
	// Directory to place build artifacts
	BuildDir string `yaml:"buildPath,omitempty"` // defaults to "build"
}

type BuildConfig struct {
	// Arbitrary flags passed to `go build`
	BuildFlags *string `yaml:"flags,omitempty"`
	// Environment variables to set during execution (incl. Exec)
	Environment map[string]string `yaml:"env,omitempty"`
	// Tags to pass to `go build`
	BuildTags types.StringOrSlice `yaml:"tags,omitempty"`
	// Should CGo be enabled?
	UseCgo bool `yaml:"cgo,omitempty"`
	// Pre- and Post-build commands
	Exec *BuildCmds `yaml:"exec,omitempty"`
	// Package to build (REQUIRED)
	Package string `yaml:"package"`
	// Targets to build for (REQUIRED)
	Targets types.StringOrSlice `yaml:"targets"`
	// Additional files that this target should include during release
	AdditionalFiles types.StringOrSlice `yaml:"include,omitempty"`
}

type BuildCmds struct {
	PreExec  types.StringOrSlice `yaml:"pre"`
	PostExec types.StringOrSlice `yaml:"post"`
}

var (
	TargetNotFoundErr = errors.New("target not found")
	NoPackageErr      = errors.New("package was not specified")
	NoTargetsErr      = errors.New("no targets were specified")
)

func (p *Project) ReifyConfig(target string) (*BuildConfig, error) {

	buildconfig := &BuildConfig{}
	tgt, ok := p.Configurations[target]
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

	// Additional files is similar.
	files := make(types.StringOrSlice, len(p.AdditionalFiles)+len(buildconfig.AdditionalFiles))
	files = append(p.AdditionalFiles, tgt.AdditionalFiles...)
	buildconfig.AdditionalFiles = files

	return buildconfig, nil
}

func (p *Project) GetSlug() string {

	slugRE, _ := regexp.Compile(`[^a-z0-9]+`)

	lower := strings.ToLower(p.Name)

	return slugRE.ReplaceAllLiteralString(lower, "-")
}

func LoadProject(path string) (*Project, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	project := &Project{}
	body, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	err = yaml.Load(body, project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (p *Project) Validate() error {
	return nil
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
