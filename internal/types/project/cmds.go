package project

import "github.com/indrora/giftwrap/internal/types"

type BuildCmds struct {
	PreExec  types.CommandList `yaml:"pre"`
	PostExec types.CommandList `yaml:"post"`
}
