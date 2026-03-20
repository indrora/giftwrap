package project

import (
	"testing"

	"go.yaml.in/yaml/v4"
)

func TestArchiveFormatConfig_UnmarshalYAML_Scalar(t *testing.T) {
	var cfg ArchiveFormatConfig
	if err := yaml.Unmarshal([]byte(`tar.gz`), &cfg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if cfg.Default != "tar.gz" {
		t.Errorf("Default = %q, want %q", cfg.Default, "tar.gz")
	}
	if len(cfg.ByOS) != 0 {
		t.Errorf("ByOS should be empty, got %v", cfg.ByOS)
	}
}

func TestArchiveFormatConfig_UnmarshalYAML_MapWithDefault(t *testing.T) {
	input := `
default: tar.gz
windows: zip
linux: tar.zst
`
	var cfg ArchiveFormatConfig
	if err := yaml.Unmarshal([]byte(input), &cfg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if cfg.Default != "tar.gz" {
		t.Errorf("Default = %q, want %q", cfg.Default, "tar.gz")
	}
	if cfg.ByOS["windows"] != "zip" {
		t.Errorf("ByOS[windows] = %q, want %q", cfg.ByOS["windows"], "zip")
	}
	if cfg.ByOS["linux"] != "tar.zst" {
		t.Errorf("ByOS[linux] = %q, want %q", cfg.ByOS["linux"], "tar.zst")
	}
	if _, hasDefault := cfg.ByOS["default"]; hasDefault {
		t.Error("ByOS should not contain \"default\" key")
	}
}

func TestArchiveFormatConfig_UnmarshalYAML_MapWithoutDefault(t *testing.T) {
	input := `
windows: zip
`
	var cfg ArchiveFormatConfig
	if err := yaml.Unmarshal([]byte(input), &cfg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if cfg.Default != "" {
		t.Errorf("Default = %q, want empty", cfg.Default)
	}
	if cfg.ByOS["windows"] != "zip" {
		t.Errorf("ByOS[windows] = %q, want %q", cfg.ByOS["windows"], "zip")
	}
}

func TestArchiveFormatConfig_IsSet(t *testing.T) {
	zero := ArchiveFormatConfig{}
	if zero.IsSet() {
		t.Error("zero value IsSet() should be false")
	}

	withDefault := ArchiveFormatConfig{Default: "tar.gz"}
	if !withDefault.IsSet() {
		t.Error("config with Default should IsSet() == true")
	}

	withByOS := ArchiveFormatConfig{ByOS: map[string]string{"windows": "zip"}}
	if !withByOS.IsSet() {
		t.Error("config with ByOS should IsSet() == true")
	}
}

func TestArchiveFormatConfig_ForOS(t *testing.T) {
	cfg := ArchiveFormatConfig{
		Default: "tar.gz",
		ByOS:    map[string]string{"windows": "zip"},
	}

	tests := []struct {
		goos string
		want string
	}{
		{"windows", "zip"},
		{"linux", "tar.gz"},
		{"darwin", "tar.gz"},
	}
	for _, tt := range tests {
		t.Run(tt.goos, func(t *testing.T) {
			got := cfg.ForOS(tt.goos)
			if got != tt.want {
				t.Errorf("ForOS(%q) = %q, want %q", tt.goos, got, tt.want)
			}
		})
	}
}

func TestArchiveFormatConfig_ForOS_Empty(t *testing.T) {
	cfg := ArchiveFormatConfig{}
	if got := cfg.ForOS("linux"); got != "" {
		t.Errorf("zero value ForOS() = %q, want empty", got)
	}
}
