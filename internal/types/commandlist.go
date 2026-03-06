package types

import (
	"fmt"

	"github.com/indrora/giftwrap/internal"
	"github.com/indrora/giftwrap/internal/runner"
	"go.yaml.in/yaml/v4"
)

type CommandList []string

func (s *CommandList) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.ScalarNode:
		*s = []string{node.Value}
	case yaml.SequenceNode:
		*s = internal.SliceDice(node.Content, func(n *yaml.Node) string {
			return n.Value
		})
	default:
		return fmt.Errorf("Couldnt make CommandList out of %T", node.Kind)
	}
	return nil
}

func (c *CommandList) Run(r runner.Runner, options runner.Options) error {
	for _, cmd := range *c {
		err := r.Run(cmd, options)
		if err != nil {
			return err
		}
	}
	return nil
}
