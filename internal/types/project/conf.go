package project

import (
	"errors"
	"maps"

	sh "github.com/indrora/giftwrap/internal/shell"
	"github.com/indrora/giftwrap/internal/types"
)

// a Project is the core structure of the Wrapfile
type Project struct {
	// (non-yaml) Directory that holds the wrapfile
	Path string `yaml:"-"`

	// name of the project (slugified for release files)
	Name string `yaml:"name"`
	// Additional files to copy into each variant's output directory after build.
	// Supports doublestar glob patterns. See types.FileSpec for YAML forms.
	AdditionalFiles []types.FileSpec `yaml:"include,omitempty"`
	// Environment variables to set during pre- and post-build commands.
	Environment map[string]string `yaml:"env,omitempty"`
	// Commands to run pre- and post-build
	Exec *BuildCmds `yaml:"exec,omitempty"`
	// Shell configures which shell to use per OS when running pre/post commands.
	// See ShellConfig for resolution order and YAML format.
	Shell sh.ShellConfig `yaml:"shell,omitempty"`
	// Directory to place build artifacts
	BuildDir       string            `yaml:"buildPath,omitempty"`      // defaults to "_build"
	DistDir        string            `yaml:"distPath,omitempty"`       // defaults to "_dist"
	DefaultTargets types.MultiString `yaml:"defaultTargets,omitempty"` // defaults to "default"
	// Archive format for release output.
	// Accepts a string ("tar.gz") or a GOOS map ({default: tar.gz, windows: zip}).
	// Built-in defaults: "zip" for windows, "tar.gz" for everything else.
	ArchiveFormat ArchiveFormatConfig `yaml:"archiveFormat,omitempty"`
	// NameTemplate is an optional Go text/template for the archive filename stem
	// (without the format extension, which is always appended automatically).
	// Available variables: .Name .Version .Major .Minor .Patch .OS .Arch
	// .Target .Format .SingleDefault
	// If unset, the default naming logic is used (backward-compatible).
	NameTemplate string `yaml:"nameTemplate,omitempty"`
	// Build configurations. Must have at least one.
	Targets map[string]BuildConfig `yaml:"targets"`
}

var (
	TargetNotFoundErr = errors.New("target not found")
	NoPackageErr      = errors.New("package was not specified")
	NoTargetsErr      = errors.New("no targets were specified")
)

// IsSingleDefault reports whether the project has exactly one target and it is the only default.
// Used by the release command to decide whether to include the target name in archive filenames.
func (p *Project) IsSingleDefault() bool {
	if len(p.Targets) != 1 || len(p.DefaultTargets) != 1 {
		return false
	}
	_, exists := p.Targets[p.DefaultTargets[0]]
	return exists
}

// EffectiveArchiveFormat returns the archive format for a given target and GOOS.
// Resolution order: target-level (per-GOOS) → project-level (per-GOOS) → built-in default.
// Built-in defaults: "zip" for windows, "tar.gz" for everything else.
func (p *Project) EffectiveArchiveFormat(targetName, goos string) string {
	if tgt, ok := p.Targets[targetName]; ok && tgt.ArchiveFormat.IsSet() {
		if f := tgt.ArchiveFormat.ForOS(goos); f != "" {
			return f
		}
	}
	if p.ArchiveFormat.IsSet() {
		if f := p.ArchiveFormat.ForOS(goos); f != "" {
			return f
		}
	}
	if goos == "windows" {
		return "zip"
	}
	return "tar.gz"
}

func (p *Project) ReifyConfig(target string) (*BuildConfig, error) {

	buildconfig := &BuildConfig{}
	tgt, ok := p.Targets[target]
	if !ok {
		return nil, TargetNotFoundErr
	}

	// convenience, since we're overwriting some things.
	*buildconfig = tgt

	// Clean up a few things

	if buildconfig.Package == "" {
		return nil, NoPackageErr
	}
	if len(buildconfig.Platforms) < 1 {
		return nil, NoTargetsErr
	}

	if buildconfig.Exec == nil {
		buildconfig.Exec = &BuildCmds{}
	}

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
