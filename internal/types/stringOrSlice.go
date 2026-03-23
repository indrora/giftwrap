package types

import (
	"fmt"

	"go.yaml.in/yaml/v4"
)

type MultiString []string

func (f *MultiString) UnmarshalYAML(node *yaml.Node) error {
	// DocumentNode wraps the real value; recurse into it.
	if node.Kind == yaml.DocumentNode {
		if len(node.Content) == 0 {
			return fmt.Errorf("MultiString: empty document node")
		}
		return f.UnmarshalYAML(node.Content[0])
	}
	switch node.Kind {
	case yaml.ScalarNode:
		*f = []string{node.Value}
	case yaml.SequenceNode:
		items := make([]string, len(node.Content))
		for i, n := range node.Content {
			items[i] = n.Value
		}
		*f = items
	case yaml.MappingNode:
		return fmt.Errorf("MultiString: cannot unmarshal a mapping into a string list")
	default:
		return fmt.Errorf("MultiString: cannot unmarshal %v into list", node.Kind)
	}
	return nil
}
