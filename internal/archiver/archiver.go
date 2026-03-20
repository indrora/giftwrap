package archiver

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
)

// ArchiveFormat identifies the archive type to produce.
type ArchiveFormat string

const (
	FormatTarGz  ArchiveFormat = "tar.gz"
	FormatTarZst ArchiveFormat = "tar.zst"
	FormatZip    ArchiveFormat = "zip"
)

// ParseFormat validates and returns the ArchiveFormat for a user-supplied string.
func ParseFormat(s string) (ArchiveFormat, error) {
	switch ArchiveFormat(s) {
	case FormatTarGz, FormatTarZst, FormatZip:
		return ArchiveFormat(s), nil
	case "":
		return FormatTarGz, nil
	default:
		return "", fmt.Errorf("unknown archive format %q: must be one of tar.gz, tar.zst, zip", s)
	}
}

// ArchiveDir creates an archive of all files in srcDir at destPath.
// The archive format is determined by format. destPath must include the correct extension.
func ArchiveDir(format ArchiveFormat, srcDir, destPath string) error {
	switch format {
	case FormatTarGz:
		return writeTarGz(srcDir, destPath)
	case FormatTarZst:
		return writeTarZst(srcDir, destPath)
	case FormatZip:
		return writeZip(srcDir, destPath)
	default:
		return fmt.Errorf("unsupported archive format: %s", format)
	}
}

func writeTarGz(srcDir, destPath string) error {
	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()

	gw := gzip.NewWriter(f)
	defer gw.Close()

	return writeTar(srcDir, gw)
}

func writeTarZst(srcDir, destPath string) error {
	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()

	zw, err := zstd.NewWriter(f)
	if err != nil {
		return err
	}
	defer zw.Close()

	return writeTar(srcDir, zw)
}

func writeTar(srcDir string, w io.Writer) error {
	tw := tar.NewWriter(w)
	defer tw.Close()

	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		hdr := &tar.Header{
			Name:    filepath.ToSlash(rel),
			Mode:    int64(info.Mode()),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(tw, file)
		return err
	})
}

func writeZip(srcDir, destPath string) error {
	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		hdr, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		hdr.Name = filepath.ToSlash(rel)
		hdr.Method = zip.Deflate

		w, err := zw.CreateHeader(hdr)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(w, file)
		return err
	})
}
