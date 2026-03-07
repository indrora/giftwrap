package builder

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/indrora/giftwrap/internal/types"
)

// copyFiles copies files matched by each FileSpec into destDir.
// Patterns use doublestar glob syntax, resolved against fsys.
// For each match, the destination path is:
//   - filepath.Join(destDir, spec.Dest, matchPath) if spec.Dest is set
//   - filepath.Join(destDir, matchPath) otherwise (preserves relative structure)
//
// Returns an error if any spec matches no files.
func copyFiles(specs []types.FileSpec, destDir string, fsys fs.FS) error {
	for _, spec := range specs {
		matches, err := doublestar.Glob(fsys, spec.Src)
		if err != nil {
			return fmt.Errorf("glob %q: %w", spec.Src, err)
		}
		if len(matches) == 0 {
			return fmt.Errorf("no files matched pattern %q", spec.Src)
		}
		for _, matchPath := range matches {
			// Skip directories — only copy files
			info, err := fs.Stat(fsys, matchPath)
			if err != nil {
				return fmt.Errorf("stat %q: %w", matchPath, err)
			}
			if info.IsDir() {
				continue
			}

			destPath := filepath.Join(destDir, spec.Dest, matchPath)
			if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
				return fmt.Errorf("mkdir for %q: %w", destPath, err)
			}

			if err := copyFile(fsys, matchPath, destPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(fsys fs.FS, srcPath, destPath string) error {
	src, err := fsys.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open %q: %w", srcPath, err)
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("create %q: %w", destPath, err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("copy %q → %q: %w", srcPath, destPath, err)
	}
	return nil
}
