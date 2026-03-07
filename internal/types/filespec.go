package types

import (
	"fmt"

	"go.yaml.in/yaml/v4"
)

// FileSpec describes a file or glob pattern to include in a build output directory.
// It can be specified as a plain string (src only) or a mapping with src and dest.
//
// YAML forms:
//
//	Plain string — copies matching files preserving their relative path:
//	  include:
//	    - README.md
//	    - docs/**
//
//	Struct form — copies matching files into a specific destination subdirectory:
//	  include:
//	    - src: "assets/**"
//	      dest: "resources/"
//
// Glob patterns follow doublestar syntax (**, *, ?, [range]).
// If dest is omitted or empty, the source's relative path is preserved.
// dest is treated as a directory prefix: assets/logo.png with dest "res/" → res/assets/logo.png
type FileSpec struct {
	Src  string `yaml:"src"`
	Dest string `yaml:"dest,omitempty"`
}

// UnmarshalYAML handles both scalar ("README.md") and
// mapping ({src: "...", dest: "..."}) YAML nodes.
func (f *FileSpec) UnmarshalYAML(node *yaml.Node) error {
	// DocumentNode wraps the real value; recurse into it.
	if node.Kind == yaml.DocumentNode {
		if len(node.Content) == 0 {
			return fmt.Errorf("filespec: empty document node")
		}
		return f.UnmarshalYAML(node.Content[0])
	}
	switch node.Kind {
	case yaml.ScalarNode:
		f.Src = node.Value
		f.Dest = ""
	case yaml.MappingNode:
		// Decode as a plain struct using a local alias to avoid recursion
		type fileSpecRaw struct {
			Src  string `yaml:"src"`
			Dest string `yaml:"dest,omitempty"`
		}
		var raw fileSpecRaw
		if err := node.Decode(&raw); err != nil {
			return err
		}
		f.Src = raw.Src
		f.Dest = raw.Dest
	default:
		return fmt.Errorf("filespec: cannot unmarshal %v into FileSpec", node.Kind)
	}
	return nil
}
