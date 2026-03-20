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
	case yaml.MappingNode:
		// Decode as a plain struct using a local alias to avoid recursion
		items := []string{}
		if err := node.Decode(items); err != nil {
			return fmt.Errorf("MultiString: Failed decode: %w", err)
		}
		if len(items) > 0 {
			*f = items
		} else {
			*f = []string{} // Empty list if there is no value
		}
	default:
		return fmt.Errorf("MultiString: cannot unmarshal %v into list", node.Kind)
	}
	return nil
}
