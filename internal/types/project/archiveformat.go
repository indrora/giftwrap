package project

import (
	"fmt"

	"go.yaml.in/yaml/v4"
)

// ArchiveFormatConfig specifies archive formats, optionally per GOOS.
//
// In a wrapfile it can be a scalar string:
//
//	archiveFormat: tar.gz
//
// Or a GOOS-keyed map with an optional "default" fallback:
//
//	archiveFormat:
//	  default: tar.gz
//	  windows: zip
//	  linux: tar.zst
//
// Valid format values: "tar.gz", "tar.zst", "zip".
type ArchiveFormatConfig struct {
	Default string
	ByOS    map[string]string
}

// IsSet reports whether any format has been explicitly configured.
func (a ArchiveFormatConfig) IsSet() bool {
	return a.Default != "" || len(a.ByOS) > 0
}

// ForOS returns the format for the given GOOS value.
// It checks ByOS for an exact match, then falls back to Default.
// Returns "" if nothing is configured for goos and Default is empty.
func (a ArchiveFormatConfig) ForOS(goos string) string {
	if f, ok := a.ByOS[goos]; ok {
		return f
	}
	return a.Default
}

// UnmarshalYAML populates an ArchiveFormatConfig from either a scalar or a mapping node.
//
//	scalar: archiveFormat: tar.gz     → Default="tar.gz", ByOS=nil
//	map:    archiveFormat:
//	          default: tar.gz         → Default="tar.gz"
//	          windows: zip            → ByOS={"windows":"zip"}
func (a *ArchiveFormatConfig) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.DocumentNode {
		if len(node.Content) == 0 {
			return nil
		}
		return a.UnmarshalYAML(node.Content[0])
	}

	switch node.Kind {
	case yaml.ScalarNode:
		a.Default = node.Value
		a.ByOS = nil
		return nil

	case yaml.MappingNode:
		raw := map[string]string{}
		if err := node.Decode(&raw); err != nil {
			return fmt.Errorf("archiveFormat: %w", err)
		}
		a.Default = raw["default"]
		delete(raw, "default")
		if len(raw) > 0 {
			a.ByOS = raw
		}
		return nil

	default:
		return fmt.Errorf("archiveFormat: expected string or map, got node kind %v", node.Kind)
	}
}
