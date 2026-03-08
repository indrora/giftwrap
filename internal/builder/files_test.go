package builder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/indrora/giftwrap/internal/types"
)

// writeTempFile creates a file at path (relative to dir) with the given content.
func writeTempFile(t *testing.T, dir, path, content string) {
	t.Helper()
	full := filepath.Join(dir, path)
	if err := os.MkdirAll(filepath.Dir(full), os.ModePerm); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(full, []byte(content), 0644); err != nil {
		t.Fatalf("write %q: %v", path, err)
	}
}

func TestCopyFiles_SingleFile(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	writeTempFile(t, src, "README.md", "hello")

	specs := []types.FileSpec{{Src: "README.md"}}
	if err := copyFiles(specs, dst, os.DirFS(src)); err != nil {
		t.Fatalf("copyFiles: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dst, "README.md"))
	if err != nil {
		t.Fatalf("read dest: %v", err)
	}
	if string(data) != "hello" {
		t.Errorf("expected 'hello', got %q", string(data))
	}
}

func TestCopyFiles_GlobPattern(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	writeTempFile(t, src, "a.md", "a")
	writeTempFile(t, src, "b.md", "b")
	writeTempFile(t, src, "c.txt", "c")

	specs := []types.FileSpec{{Src: "*.md"}}
	if err := copyFiles(specs, dst, os.DirFS(src)); err != nil {
		t.Fatalf("copyFiles: %v", err)
	}

	for _, name := range []string{"a.md", "b.md"} {
		if _, err := os.Stat(filepath.Join(dst, name)); err != nil {
			t.Errorf("expected %q to exist: %v", name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(dst, "c.txt")); !os.IsNotExist(err) {
		t.Error("c.txt should not have been copied")
	}
}

func TestCopyFiles_RecursiveGlob(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	writeTempFile(t, src, "docs/guide.md", "guide")
	writeTempFile(t, src, "docs/sub/ref.md", "ref")

	specs := []types.FileSpec{{Src: "docs/**"}}
	if err := copyFiles(specs, dst, os.DirFS(src)); err != nil {
		t.Fatalf("copyFiles: %v", err)
	}

	for _, rel := range []string{"docs/guide.md", "docs/sub/ref.md"} {
		if _, err := os.Stat(filepath.Join(dst, rel)); err != nil {
			t.Errorf("expected %q to exist: %v", rel, err)
		}
	}
}

func TestCopyFiles_WithDest(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	writeTempFile(t, src, "assets/logo.png", "img")

	specs := []types.FileSpec{{Src: "assets/**", Dest: "resources/"}}
	if err := copyFiles(specs, dst, os.DirFS(src)); err != nil {
		t.Fatalf("copyFiles: %v", err)
	}

	want := filepath.Join(dst, "resources", "logo.png")
	if _, err := os.Stat(want); err != nil {
		t.Errorf("expected %q to exist: %v", want, err)
	}
}

func TestCopyFiles_WithDest_TopLevelGlob(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	writeTempFile(t, src, "README.md", "hello")

	specs := []types.FileSpec{{Src: "*.md", Dest: "docs/"}}
	if err := copyFiles(specs, dst, os.DirFS(src)); err != nil {
		t.Fatalf("copyFiles: %v", err)
	}

	want := filepath.Join(dst, "docs", "README.md")
	if _, err := os.Stat(want); err != nil {
		t.Errorf("expected %q to exist: %v", want, err)
	}
}

func TestCopyFiles_PreservesStructure(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	writeTempFile(t, src, "nested/dir/file.txt", "content")

	specs := []types.FileSpec{{Src: "nested/**"}}
	if err := copyFiles(specs, dst, os.DirFS(src)); err != nil {
		t.Fatalf("copyFiles: %v", err)
	}

	want := filepath.Join(dst, "nested", "dir", "file.txt")
	if _, err := os.Stat(want); err != nil {
		t.Errorf("expected %q to exist: %v", want, err)
	}
}

func TestCopyFiles_NoMatchesError(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	specs := []types.FileSpec{{Src: "nonexistent/**"}}
	err := copyFiles(specs, dst, os.DirFS(src))
	if err == nil {
		t.Fatal("expected error for unmatched pattern, got nil")
	}
}

func TestCopyFiles_EmptySpecs(t *testing.T) {
	dst := t.TempDir()
	if err := copyFiles(nil, dst, os.DirFS(dst)); err != nil {
		t.Fatalf("expected no error for empty specs, got: %v", err)
	}
	if err := copyFiles([]types.FileSpec{}, dst, os.DirFS(dst)); err != nil {
		t.Fatalf("expected no error for empty specs, got: %v", err)
	}
}
