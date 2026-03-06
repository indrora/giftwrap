package types

import (
	"fmt"

	"github.com/indrora/giftwrap/internal"
	"go.yaml.in/yaml/v4"
)

type StringOrSlice []string

func (s *StringOrSlice) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.ScalarNode:
		*s = []string{node.Value}
	case yaml.SequenceNode:
		*s = internal.SliceDice(node.Content, func(n *yaml.Node) string {
			return n.Value
		})
	default:
		return fmt.Errorf("Couldnt make stringSlice out of %T", node.Kind)
	}
	return nil
}
