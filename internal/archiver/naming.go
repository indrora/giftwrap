package archiver

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// ArchiveName returns the archive filename for a built variant.
//
// If singleDefault is true (exactly one target and it is the default),
// the target name is omitted:
//
//	<projectName>-<goos>-<goarch>.<format>
//
// Otherwise:
//
//	<projectName>-<targetName>-<goos>-<goarch>.<format>
//
// variant must be in "goos/goarch" form (e.g. "linux/amd64").
func ArchiveName(projectName, targetName, variant string, format ArchiveFormat, singleDefault bool) (string, error) {
	parts := strings.SplitN(variant, "/", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid variant %q: expected goos/goarch", variant)
	}
	goos, goarch := parts[0], parts[1]
	if singleDefault {
		return fmt.Sprintf("%s-%s-%s.%s", projectName, goos, goarch, format), nil
	}
	return fmt.Sprintf("%s-%s-%s-%s.%s", projectName, targetName, goos, goarch, format), nil
}

// ArchiveNameData holds the variables available inside a nameTemplate.
type ArchiveNameData struct {
	Name          string // project name
	Version       string // full version string, e.g. "v1.2.3" or "v1.2.3+4.gabc1234"
	VersionRC     string // vM.m.(p+1)+rcN form; same as Version on an exact tag
	Major         uint64
	Minor         uint64
	Patch         uint64
	Commits       int    // commits ahead of the last semver tag (0 = exact tag)
	OS            string // GOOS
	Arch          string // GOARCH
	Target        string // target name from the wrapfile
	Format        string // archive format extension, e.g. "tar.gz"
	SingleDefault bool   // true when exactly one target is the sole default
}

// ArchiveNameFromTemplate renders tmplStr with data to produce an archive filename.
// The template controls the stem; the format extension is always appended.
// variant must be in "goos/goarch" form (e.g. "linux/amd64").
func ArchiveNameFromTemplate(tmplStr string, data ArchiveNameData) (string, error) {
	tmpl, err := template.New("archiveName").Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("invalid nameTemplate: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("nameTemplate execution failed: %w", err)
	}
	return buf.String() + "." + data.Format, nil
}
