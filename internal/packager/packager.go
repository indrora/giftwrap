package packager

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/indrora/giftwrap/internal"
	"github.com/indrora/giftwrap/internal/archiver"
	"github.com/indrora/giftwrap/internal/builder"
	"github.com/indrora/giftwrap/internal/runner"
	"github.com/indrora/giftwrap/internal/types/project"
)

// Packager builds and archives release targets.
type Packager struct {
	b             *builder.Builder
	proj          project.Project
	Shell         string
	singleDefault bool
}

// NewPackager constructs a Packager backed by the given project and runner.
func NewPackager(p project.Project, r runner.Runner) (*Packager, error) {
	b, err := builder.NewBuilder(p, r)
	if err != nil {
		return nil, err
	}
	return &Packager{
		b:             b,
		proj:          p,
		Shell:         b.Shell,
		singleDefault: p.IsSingleDefault(),
	}, nil
}

// SetIO forwards output writers to the underlying builder.
func (p *Packager) SetIO(out, err io.Writer) {
	p.b.SetIO(out, err)
}

// Setup creates output directories and runs project-level pre-exec hooks.
func (p *Packager) Setup() error {
	p.b.Shell = p.Shell
	return p.b.Setup()
}

// PackageTarget builds the given variant of target and archives the output.
func (p *Packager) PackageTarget(target, variantName string) error {
	p.b.Shell = p.Shell
	goos := strings.SplitN(variantName, "/", 2)[0]

	formatStr := p.proj.EffectiveArchiveFormat(target, goos)
	format, err := archiver.ParseFormat(formatStr)
	if err != nil {
		return fmt.Errorf("invalid archive format for %s (%s): %w", target, variantName, err)
	}

	if err := p.b.BuildTarget(target, variantName); err != nil {
		return fmt.Errorf("build failed for %s (%s): %w", target, variantName, err)
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

// Teardown runs project-level post-exec hooks.
func (p *Packager) Teardown() error {
	p.b.Shell = p.Shell
	return p.b.Teardown()
}
