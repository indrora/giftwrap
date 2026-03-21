package packager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/indrora/giftwrap/internal"
	"github.com/indrora/giftwrap/internal/archiver"
	"github.com/indrora/giftwrap/internal/packager/version"
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
	parts := strings.SplitN(variantName, "/", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid variant %q: expected goos/goarch", variantName)
	}
	goos, goarch := parts[0], parts[1]

	formatStr := p.proj.EffectiveArchiveFormat(target, goos)
	format, err := archiver.ParseFormat(formatStr)
	if err != nil {
		return fmt.Errorf("invalid archive format for %s (%s): %w", target, variantName, err)
	}

	var archiveName string
	if p.proj.NameTemplate != "" {
		ver, verErr := version.GetVersion(p.proj.Path)
		if verErr != nil {
			return fmt.Errorf("nameTemplate requires a semver git tag, but version could not be resolved: %w", verErr)
		}
		data := archiver.ArchiveNameData{
			Name:          p.proj.Name,
			Version:       ver.String(),
			VersionRC:     ver.RCString(),
			Major:         ver.Major(),
			Minor:         ver.Minor(),
			Patch:         ver.Patch(),
			Commits:       ver.Commits,
			OS:            goos,
			Arch:          goarch,
			Target:        target,
			Format:        string(format),
			SingleDefault: p.singleDefault,
		}
		archiveName, err = archiver.ArchiveNameFromTemplate(p.proj.NameTemplate, data)
		if err != nil {
			return fmt.Errorf("could not render nameTemplate for %s (%s): %w", target, variantName, err)
		}
	} else {
		archiveName, err = archiver.ArchiveName(p.proj.Name, target, variantName, format, p.singleDefault)
		if err != nil {
			return fmt.Errorf("could not construct archive name for %s (%s): %w", target, variantName, err)
		}
	}

	srcDir := filepath.Join(p.proj.BuildDir, internal.Slugify(target), internal.Slugify(variantName))
	destPath := filepath.Join(p.proj.DistDir, archiveName)

	if err := archiver.ArchiveDir(format, srcDir, destPath); err != nil {
		return fmt.Errorf("failed to create archive %s: %w", destPath, err)
	}

	return nil
}
