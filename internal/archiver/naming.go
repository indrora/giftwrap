package archiver

import (
	"fmt"
	"strings"
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
