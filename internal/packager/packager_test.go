package packager

import (
	"path/filepath"
	"testing"

	"github.com/indrora/giftwrap/internal/runner"
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
	r := runner.PrintRunner{}

	pkg, err := NewPackager(proj, r)
	if err != nil {
		t.Fatalf("NewPackager: %v", err)
	}

	if err := pkg.Setup(); err != nil {
		t.Fatalf("Setup: %v", err)
	}

	if err := pkg.PackageTarget("default", "linux/amd64"); err != nil {
		t.Errorf("PackageTarget: %v", err)
	}
}

func TestPackageTargetUnknownTarget(t *testing.T) {
	proj := makeTestProject(t)
	r := runner.PrintRunner{}

	pkg, err := NewPackager(proj, r)
	if err != nil {
		t.Fatalf("NewPackager: %v", err)
	}

	if err := pkg.Setup(); err != nil {
		t.Fatalf("Setup: %v", err)
	}

	if err := pkg.PackageTarget("nonexistent", "linux/amd64"); err == nil {
		t.Error("expected error for unknown target, got nil")
	}
}

func TestPackageTargetBadArchiveFormat(t *testing.T) {
	proj := makeTestProject(t)
	proj.ArchiveFormat = project.ArchiveFormatConfig{Default: "invalid"}
	r := runner.PrintRunner{}

	pkg, err := NewPackager(proj, r)
	if err != nil {
		t.Fatalf("NewPackager: %v", err)
	}

	if err := pkg.Setup(); err != nil {
		t.Fatalf("Setup: %v", err)
	}

	if err := pkg.PackageTarget("default", "linux/amd64"); err == nil {
		t.Error("expected error for invalid archive format, got nil")
	}
}
