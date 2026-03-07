package types

import (
	"testing"

	"go.yaml.in/yaml/v4"
)

func TestFileSpec_UnmarshalScalar(t *testing.T) {
	input := `README.md`
	var spec FileSpec
	if err := yaml.Unmarshal([]byte(input), &spec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if spec.Src != "README.md" {
		t.Errorf("expected Src=README.md, got %q", spec.Src)
	}
	if spec.Dest != "" {
		t.Errorf("expected Dest empty, got %q", spec.Dest)
	}
}

func TestFileSpec_UnmarshalStruct(t *testing.T) {
	input := `src: "assets/**"
dest: "res/"`
	var spec FileSpec
	if err := yaml.Unmarshal([]byte(input), &spec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if spec.Src != "assets/**" {
		t.Errorf("expected Src=assets/**, got %q", spec.Src)
	}
	if spec.Dest != "res/" {
		t.Errorf("expected Dest=res/, got %q", spec.Dest)
	}
}

func TestFileSpec_UnmarshalList(t *testing.T) {
	input := `
- README.md
- src: "assets/**"
  dest: "res/"
- docs/**
`
	var specs []FileSpec
	if err := yaml.Unmarshal([]byte(input), &specs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(specs) != 3 {
		t.Fatalf("expected 3 specs, got %d", len(specs))
	}

	if specs[0].Src != "README.md" || specs[0].Dest != "" {
		t.Errorf("specs[0]: got src=%q dest=%q", specs[0].Src, specs[0].Dest)
	}
	if specs[1].Src != "assets/**" || specs[1].Dest != "res/" {
		t.Errorf("specs[1]: got src=%q dest=%q", specs[1].Src, specs[1].Dest)
	}
	if specs[2].Src != "docs/**" || specs[2].Dest != "" {
		t.Errorf("specs[2]: got src=%q dest=%q", specs[2].Src, specs[2].Dest)
	}
}
