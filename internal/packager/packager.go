package packager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/indrora/giftwrap/internal"
	"github.com/indrora/giftwrap/internal/archiver"
	"github.com/indrora/giftwrap/internal/types/project"
)

// Packager archives built release targets into the distribution directory.
// It does not build — callers are responsible for building first (e.g. via builder.Builder).
type Packager struct {
	proj          project.Project
	singleDefault bool
}

// NewPackager constructs a Packager for the given project.
func NewPackager(p project.Project) *Packager {
	return &Packager{
		proj:          p,
		singleDefault: p.IsSingleDefault(),
	}
}

// Setup creates the distribution directory.
func (p *Packager) Setup() error {
	if err := os.MkdirAll(p.proj.DistDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create dist dir %s: %w", p.proj.DistDir, err)
	}
	return nil
}

// PackageTarget archives the build output for the given target variant.
// The build output must already exist at the expected path under BuildDir.
func (p *Packager) PackageTarget(target, variantName string) error {
	goos := strings.SplitN(variantName, "/", 2)[0]

	formatStr := p.proj.EffectiveArchiveFormat(target, goos)
	format, err := archiver.ParseFormat(formatStr)
	if err != nil {
		return fmt.Errorf("invalid archive format for %s (%s): %w", target, variantName, err)
	}

	archiveName, err := archiver.ArchiveName(p.proj.Name, target, variantName, format, p.singleDefault)
	if err != nil {
		return fmt.Errorf("could not construct archive name for %s (%s): %w", target, variantName, err)
	}

	srcDir := filepath.Join(p.proj.BuildDir, internal.Slugify(target), internal.Slugify(variantName))
	destPath := filepath.Join(p.proj.DistDir, archiveName)

	if err := archiver.ArchiveDir(format, srcDir, destPath); err != nil {
		return fmt.Errorf("failed to create archive %s: %w", destPath, err)
	}

	return nil
}
