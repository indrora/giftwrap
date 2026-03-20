package archiver

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/klauspost/compress/zstd"
)

// setupSrcDir creates a temp directory with two test files and returns its path.
func setupSrcDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "hello.txt"), []byte("hello world"), 0644); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "deep.txt"), []byte("deep file"), 0644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input   string
		want    ArchiveFormat
		wantErr bool
	}{
		{"tar.gz", FormatTarGz, false},
		{"tar.zst", FormatTarZst, false},
		{"zip", FormatZip, false},
		{"", FormatTarGz, false},
		{"bz2", "", true},
		{"invalid", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseFormat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFormat(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ParseFormat(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestArchiveDirTarGz(t *testing.T) {
	srcDir := setupSrcDir(t)
	destPath := filepath.Join(t.TempDir(), "out.tar.gz")

	if err := ArchiveDir(FormatTarGz, srcDir, destPath); err != nil {
		t.Fatalf("ArchiveDir: %v", err)
	}

	files := extractTarGz(t, destPath)
	assertFiles(t, files, map[string]string{
		"hello.txt":    "hello world",
		"sub/deep.txt": "deep file",
	})
}

func TestArchiveDirTarZst(t *testing.T) {
	srcDir := setupSrcDir(t)
	destPath := filepath.Join(t.TempDir(), "out.tar.zst")

	if err := ArchiveDir(FormatTarZst, srcDir, destPath); err != nil {
		t.Fatalf("ArchiveDir: %v", err)
	}

	files := extractTarZst(t, destPath)
	assertFiles(t, files, map[string]string{
		"hello.txt":    "hello world",
		"sub/deep.txt": "deep file",
	})
}

func TestArchiveDirZip(t *testing.T) {
	srcDir := setupSrcDir(t)
	destPath := filepath.Join(t.TempDir(), "out.zip")

	if err := ArchiveDir(FormatZip, srcDir, destPath); err != nil {
		t.Fatalf("ArchiveDir: %v", err)
	}

	files := extractZip(t, destPath)
	assertFiles(t, files, map[string]string{
		"hello.txt":    "hello world",
		"sub/deep.txt": "deep file",
	})
}

func assertFiles(t *testing.T, got map[string]string, want map[string]string) {
	t.Helper()
	for name, wantContent := range want {
		gotContent, ok := got[name]
		if !ok {
			t.Errorf("missing file %q in archive", name)
			continue
		}
		if gotContent != wantContent {
			t.Errorf("file %q: content = %q, want %q", name, gotContent, wantContent)
		}
	}
	for name := range got {
		if _, ok := want[name]; !ok {
			t.Errorf("unexpected file %q in archive", name)
		}
	}
}

func extractTarGz(t *testing.T, path string) map[string]string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	defer gr.Close()
	return readTar(t, gr)
}

func extractTarZst(t *testing.T, path string) map[string]string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	zr, err := zstd.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	defer zr.Close()
	return readTar(t, zr)
}

func readTar(t *testing.T, r io.Reader) map[string]string {
	t.Helper()
	files := make(map[string]string)
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		body, err := io.ReadAll(tr)
		if err != nil {
			t.Fatal(err)
		}
		files[hdr.Name] = string(body)
	}
	return files
}

func extractZip(t *testing.T, path string) map[string]string {
	t.Helper()
	zr, err := zip.OpenReader(path)
	if err != nil {
		t.Fatal(err)
	}
	defer zr.Close()
	files := make(map[string]string)
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			t.Fatal(err)
		}
		body, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			t.Fatal(err)
		}
		files[filepath.ToSlash(f.Name)] = string(body)
	}
	return files
}

// Ensure ArchiveDir returns an error for invalid format.
func TestArchiveDirUnknownFormat(t *testing.T) {
	err := ArchiveDir("bz2", t.TempDir(), filepath.Join(t.TempDir(), "out.bz2"))
	if err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}
