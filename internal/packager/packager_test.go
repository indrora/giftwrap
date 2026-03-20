package packager

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/indrora/giftwrap/internal/types/project"
)

func makeTestProject(t *testing.T) project.Project {
	t.Helper()
	tmp := t.TempDir()
	return project.Project{
		Name:           "testapp",
		BuildDir:       filepath.Join(tmp, "_build"),
		DistDir:        filepath.Join(tmp, "_dist"),
		DefaultTargets: []string{"default"},
		Exec:           &project.BuildCmds{},
		Targets: map[string]project.BuildConfig{
			"default": {
				Package: ".",
				Targets: []string{"linux/amd64"},
				Exec:    &project.BuildCmds{},
			},
		},
	}
}

func TestPackageTargetSuccess(t *testing.T) {
	proj := makeTestProject(t)
	pkg := NewPackager(proj)

	if err := pkg.Setup(); err != nil {
		t.Fatalf("Setup: %v", err)
	}

	// Simulate a prior build by creating the expected source directory.
	srcDir := filepath.Join(proj.BuildDir, "default", "linux-amd64")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatalf("creating srcDir: %v", err)
	}

	if err := pkg.PackageTarget("default", "linux/amd64"); err != nil {
		t.Errorf("PackageTarget: %v", err)
	}
}

func TestPackageTargetBadArchiveFormat(t *testing.T) {
	proj := makeTestProject(t)
	proj.ArchiveFormat = project.ArchiveFormatConfig{Default: "invalid"}
	pkg := NewPackager(proj)

	if err := pkg.Setup(); err != nil {
		t.Fatalf("Setup: %v", err)
	}

	if err := pkg.PackageTarget("default", "linux/amd64"); err == nil {
		t.Error("expected error for invalid archive format, got nil")
	}
}
