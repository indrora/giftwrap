package project

import "github.com/indrora/giftwrap/internal/types"

type BuildConfig struct {
	// Arbitrary flags passed to `go build`
	BuildFlags *string `yaml:"flags,omitempty"`
	// Environment variables to set during execution (incl. Exec)
	Environment map[string]string `yaml:"env,omitempty"`
	// Tags to pass to `go build`
	BuildTags types.CommandList `yaml:"tags,omitempty"`
	// Should CGo be enabled?
	UseCgo bool `yaml:"cgo,omitempty"`
	// Pre- and Post-build commands
	Exec *BuildCmds `yaml:"exec,omitempty"`
	// Package to build (REQUIRED)
	Package string `yaml:"package"`
	// Targets to build for (REQUIRED)
	Targets types.CommandList `yaml:"targets"`
	// Additional files that this target should include during release
	AdditionalFiles types.CommandList `yaml:"include,omitempty"`
}
